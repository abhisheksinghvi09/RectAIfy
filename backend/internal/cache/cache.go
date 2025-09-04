package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/singleflight"

	"realitycheck/pkg/types"
)

// Cache provides multi-level caching with LRU + Postgres + singleflight
type Cache struct {
	lru *lru.Cache[string, *CacheEntry]
	db  *pgxpool.Pool
	sf  singleflight.Group
	ttl time.Duration
}

// CacheEntry represents a cached item
type CacheEntry struct {
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
	TTL       time.Duration   `json:"ttl"`
}

// NewCache creates a new cache instance
func NewCache(db *pgxpool.Pool, lruSize int, ttl time.Duration) (*Cache, error) {
	lruCache, err := lru.New[string, *CacheEntry](lruSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create LRU cache: %w", err)
	}

	c := &Cache{
		lru: lruCache,
		db:  db,
		ttl: ttl,
	}

	// Warm up cache on startup
	go c.warmupFromDB(context.Background())

	return c, nil
}

// Get retrieves data from cache with read-through to database
func (c *Cache) Get(ctx context.Context, key string) (json.RawMessage, bool, error) {
	hash := c.hashKey(key)

	// Use singleflight to deduplicate concurrent requests
	result, err, _ := c.sf.Do(hash, func() (interface{}, error) {
		return c.get(ctx, key, hash)
	})

	if err != nil {
		return nil, false, err
	}

	if result == nil {
		return nil, false, nil
	}

	entry, ok := result.(*CacheEntry)
	if !ok || entry == nil {
		return nil, false, nil
	}
	return entry.Data, true, nil
}

// Set stores data in both LRU and database
func (c *Cache) Set(ctx context.Context, key string, data json.RawMessage) error {
	hash := c.hashKey(key)

	entry := &CacheEntry{
		Data:      data,
		CreatedAt: time.Now(),
		TTL:       c.ttl,
	}

	// Store in LRU
	c.lru.Add(hash, entry)

	// Store in database (only if database is available)
	if c.db != nil {
		return c.setDB(ctx, hash, key, data)
	}
	return nil
}

// get implements the actual cache retrieval logic
func (c *Cache) get(ctx context.Context, key, hash string) (*CacheEntry, error) {
	// Check LRU first
	if entry, exists := c.lru.Get(hash); exists {
		if !c.isExpired(entry) {
			return entry, nil
		}
		// Remove expired entry
		c.lru.Remove(hash)
	}

	// Check database (only if database is available)
	var entry *CacheEntry
	var found bool
	var err error

	if c.db != nil {
		entry, found, err = c.getDB(ctx, hash)
		if err != nil {
			return nil, fmt.Errorf("database lookup failed: %w", err)
		}
	}

	if found && !c.isExpired(entry) {
		// Populate LRU with fresh data from DB
		c.lru.Add(hash, entry)
		return entry, nil
	}

	// Clean up expired entry from DB
	if found && c.isExpired(entry) {
		go c.deleteDB(context.Background(), hash)
	}

	return nil, nil
}

// getDB retrieves entry from database
func (c *Cache) getDB(ctx context.Context, hash string) (*CacheEntry, bool, error) {
	var result json.RawMessage
	var createdAt time.Time
	var ttlSeconds int

	err := c.db.QueryRow(ctx,
		"SELECT result, created_at, ttl_seconds FROM web_cache WHERE hash = $1",
		hash,
	).Scan(&result, &createdAt, &ttlSeconds)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, false, nil
		}
		return nil, false, err
	}

	entry := &CacheEntry{
		Data:      result,
		CreatedAt: createdAt,
		TTL:       time.Duration(ttlSeconds) * time.Second,
	}

	return entry, true, nil
}

// setDB stores entry in database
func (c *Cache) setDB(ctx context.Context, hash, key string, data json.RawMessage) error {
	_, err := c.db.Exec(ctx,
		`INSERT INTO web_cache (hash, query, result, created_at, ttl_seconds) 
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (hash) DO UPDATE SET 
		 query = EXCLUDED.query,
		 result = EXCLUDED.result,
		 created_at = EXCLUDED.created_at,
		 ttl_seconds = EXCLUDED.ttl_seconds`,
		hash, key, data, time.Now(), int(c.ttl.Seconds()),
	)
	return err
}

// deleteDB removes entry from database
func (c *Cache) deleteDB(ctx context.Context, hash string) error {
	_, err := c.db.Exec(ctx, "DELETE FROM web_cache WHERE hash = $1", hash)
	return err
}

// isExpired checks if a cache entry has expired
func (c *Cache) isExpired(entry *CacheEntry) bool {
	return time.Since(entry.CreatedAt) > entry.TTL
}

// hashKey creates a stable hash for cache keys
func (c *Cache) hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", hash)
}

// warmupFromDB loads recent cache entries into LRU on startup
func (c *Cache) warmupFromDB(ctx context.Context) {
	if c.db == nil {
		return
	}

	rows, err := c.db.Query(ctx,
		`SELECT hash, result, created_at, ttl_seconds 
		 FROM web_cache 
		 WHERE created_at + (ttl_seconds || ' seconds')::INTERVAL > NOW()
		 ORDER BY created_at DESC 
		 LIMIT 1000`,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var hash string
		var result json.RawMessage
		var createdAt time.Time
		var ttlSeconds int

		if err := rows.Scan(&hash, &result, &createdAt, &ttlSeconds); err != nil {
			continue
		}

		entry := &CacheEntry{
			Data:      result,
			CreatedAt: createdAt,
			TTL:       time.Duration(ttlSeconds) * time.Second,
		}

		if !c.isExpired(entry) {
			c.lru.Add(hash, entry)
		}
	}
}

// CleanupExpired removes expired entries from database
func (c *Cache) CleanupExpired(ctx context.Context) error {
	_, err := c.db.Exec(ctx,
		"DELETE FROM web_cache WHERE created_at + (ttl_seconds || ' seconds')::INTERVAL < NOW()",
	)
	return err
}

// EvidenceCache provides specialized caching for search evidence
type EvidenceCache struct {
	cache *Cache
}

// StartCleanupWorker starts a background worker to clean expired entries
func (ec *EvidenceCache) StartCleanupWorker(ctx context.Context, interval time.Duration) {
	ec.cache.StartCleanupWorker(ctx, interval)
}

// NewEvidenceCache creates a cache specifically for evidence
func NewEvidenceCache(db *pgxpool.Pool, lruSize int, ttl time.Duration) (*EvidenceCache, error) {
	cache, err := NewCache(db, lruSize, ttl)
	if err != nil {
		return nil, err
	}

	return &EvidenceCache{cache: cache}, nil
}

// GetEvidence retrieves cached evidence for a query
func (ec *EvidenceCache) GetEvidence(ctx context.Context, query string) ([]types.Evidence, bool, error) {
	data, found, err := ec.cache.Get(ctx, query)
	if err != nil || !found {
		return nil, found, err
	}

	var evidence []types.Evidence
	if err := json.Unmarshal(data, &evidence); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal evidence: %w", err)
	}

	return evidence, true, nil
}

// SetEvidence stores evidence in cache
func (ec *EvidenceCache) SetEvidence(ctx context.Context, query string, evidence []types.Evidence) error {
	data, err := json.Marshal(evidence)
	if err != nil {
		return fmt.Errorf("failed to marshal evidence: %w", err)
	}

	return ec.cache.Set(ctx, query, data)
}

// StartCleanupWorker starts a background worker to clean expired entries
func (c *Cache) StartCleanupWorker(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.CleanupExpired(ctx); err != nil {
				// Log error but continue
				continue
			}
		}
	}
}
