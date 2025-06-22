// Package middlewares provides rate limiting functionality.
package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/yeferson59/gin-template/pkg/logger"
	"github.com/yeferson59/gin-template/pkg/response"
)

// IPRateLimiter contains the rate limiters for each IP address.
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewIPRateLimiter creates a new IP-based rate limiter.
func NewIPRateLimiter(rps rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rps,
		burst:    burst,
	}
}

// GetLimiter returns the rate limiter for the given IP address.
func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
	}

	return limiter
}

// CleanupOldLimiters removes limiters for IPs that haven't been used recently.
func (rl *IPRateLimiter) CleanupOldLimiters() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// In a production environment, you might want to implement a more sophisticated
	// cleanup mechanism, possibly using a separate goroutine and tracking last access times
	if len(rl.limiters) > 1000 {
		// Simple cleanup: remove half of the limiters
		count := 0
		for ip := range rl.limiters {
			if count > 500 {
				break
			}
			delete(rl.limiters, ip)
			count++
		}
	}
}

var globalRateLimiter = NewIPRateLimiter(rate.Every(time.Second), 10) // 10 requests per second per IP

// RateLimit returns a middleware that limits requests per IP address.
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := globalRateLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			logger.WithField("ip", ip).Warn("Rate limit exceeded")
			response.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded", "Too many requests from your IP address")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitWithConfig returns a middleware with custom rate limiting configuration.
func RateLimitWithConfig(rps rate.Limit, burst int) gin.HandlerFunc {
	rateLimiter := NewIPRateLimiter(rps, burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rateLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			logger.WithField("ip", ip).Warn("Rate limit exceeded")
			response.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded", "Too many requests from your IP address")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimit provides stricter rate limiting for authentication endpoints.
func AuthRateLimit() gin.HandlerFunc {
	authLimiter := NewIPRateLimiter(rate.Every(time.Minute), 5) // 5 attempts per minute per IP

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := authLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			logger.WithField("ip", ip).Warn("Auth rate limit exceeded")
			response.ErrorResponse(c, http.StatusTooManyRequests, "AUTH_RATE_LIMIT_EXCEEDED", "Authentication rate limit exceeded", "Too many authentication attempts from your IP address")
			c.Abort()
			return
		}

		c.Next()
	}
}
