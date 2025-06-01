package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements a rate limiter using token bucket algorithm
type RateLimiter struct {
	limiter map[string]*rate.Limiter
	mu      sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiter: make(map[string]*rate.Limiter),
	}
}

// getLimiter returns the rate limiter for a given key (usually IP address)
func (rl *RateLimiter) getLimiter(key string, limit rate.Limit, burst int) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiter[key]
	if !exists {
		limiter = rate.NewLimiter(limit, burst)
		rl.limiter[key] = limiter
	}

	return limiter
}

// RateLimit middleware limits requests per IP address
func (rl *RateLimiter) RateLimit(rps int, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the key
		key := c.ClientIP()

		// Get the limiter for this IP
		limiter := rl.getLimiter(key, rate.Limit(rps), burst)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Cleanup removes old limiters to prevent memory leaks
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Clear all limiters (in production, you might want more sophisticated cleanup)
	rl.limiter = make(map[string]*rate.Limiter)
}

// AuthMiddleware provides basic authentication validation
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, this is a placeholder for authentication
		// In a real application, you would validate JWT tokens or API keys here

		// Check for API key in header
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Allow requests without API key for public endpoints
			c.Next()
			return
		}

		// Basic API key validation (in production, use proper validation)
		validAPIKeys := map[string]bool{
			"dev-key-123":   true,
			"admin-key-456": true,
		}

		if !validAPIKeys[apiKey] {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds common security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'")

		c.Next()
	}
}

// LoggingMiddleware provides structured logging
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// generateRequestID creates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
