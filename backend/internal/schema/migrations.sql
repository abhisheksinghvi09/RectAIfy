-- RectAIfy Database Schema
-- PostgreSQL 15+

-- Create the main analyses table
CREATE TABLE IF NOT EXISTS analyses (
    id TEXT PRIMARY KEY,
    idea JSONB NOT NULL,
    result JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create the evidence table for research citations
CREATE TABLE IF NOT EXISTS evidence (
    id TEXT PRIMARY KEY,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    snippet TEXT,
    published_at TIMESTAMPTZ,
    retrieved_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    source_type TEXT
);

-- Create the many-to-many relationship table for analysis-evidence
CREATE TABLE IF NOT EXISTS analysis_evidence (
    analysis_id TEXT REFERENCES analyses(id) ON DELETE CASCADE,
    evidence_id TEXT REFERENCES evidence(id) ON DELETE CASCADE,
    PRIMARY KEY(analysis_id, evidence_id)
);

-- Create the web cache table for search results
CREATE TABLE IF NOT EXISTS web_cache (
    hash TEXT PRIMARY KEY,
    query TEXT NOT NULL,
    result JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ttl_seconds INTEGER NOT NULL DEFAULT 86400
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_evidence_url_hash ON evidence(MD5(url));
CREATE INDEX IF NOT EXISTS idx_analyses_result_gin ON analyses USING GIN (result jsonb_path_ops);
CREATE INDEX IF NOT EXISTS idx_web_cache_created_at ON web_cache (created_at);
CREATE INDEX IF NOT EXISTS idx_evidence_retrieved_at ON evidence (retrieved_at);
CREATE INDEX IF NOT EXISTS idx_analyses_created_at ON analyses (created_at);

-- Create index for cache expiration cleanup
-- Note: Using a simpler index on created_at since the complex expression isn't IMMUTABLE
CREATE INDEX IF NOT EXISTS idx_web_cache_expiry ON web_cache (created_at, ttl_seconds);

-- Create extension for better JSON operations if available
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
