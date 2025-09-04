package httpx

import (
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware provides bearer token authentication
func AuthMiddleware(bearerToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if bearerToken == "" {
				// If no token configured, skip auth
				next.ServeHTTP(w, r)
				return
			}

			// Allow OPTIONS requests to pass through without authentication
			// (CORS preflight requests should not require auth)
			if r.Method == "OPTIONS" {
				next.ServeHTTP(w, r)
				return
			}

			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "Bearer token required", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			if token != bearerToken {
				http.Error(w, "Invalid bearer token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow specific origins for development
		allowedOrigins := []string{
			"http://localhost:5173", // Vite default
			"http://localhost:5174", // Vite alternative
			"http://localhost:5175", // Vite alternative
			"http://localhost:5176", // User's frontend
			"http://localhost:3000", // React default
			"http://localhost:3001", // React alternative
			"http://127.0.0.1:5173",
			"http://127.0.0.1:5174",
			"http://127.0.0.1:5175",
			"http://127.0.0.1:5176",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:3001",
		}

		// Check if origin is allowed
		isAllowed := false
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				isAllowed = true
				break
			}
		}

		// Set CORS headers
		if isAllowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			// Fallback to wildcard for other origins (production, etc.)
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer that captures status code
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		// In a real implementation, use a proper logger
		println(
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration.String(),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
