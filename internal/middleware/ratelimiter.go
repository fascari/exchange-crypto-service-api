package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"exchange-crypto-service-api/internal/config"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

const limitExceededResponse = `{"error": "rate limit exceeded", "message": "too many requests"}`

type (
	clientLimiter struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	RateLimiter struct {
		limiters map[string]*clientLimiter
		mu       sync.RWMutex
		config   *config.RateLimiter
	}
)

func newRateLimiterByIP(rateLimiterConfig *config.RateLimiter) *RateLimiter {
	if rateLimiterConfig == nil {
		defaultConfig := config.LoadRateLimiter()
		rateLimiterConfig = &defaultConfig
	}

	rl := &RateLimiter{
		limiters: make(map[string]*clientLimiter),
		config:   rateLimiterConfig,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Limiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[ip]
	if !exists {
		limiter = &clientLimiter{
			limiter:  rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize),
			lastSeen: time.Now(),
		}
		rl.limiters[ip] = limiter
	}

	limiter.lastSeen = time.Now()
	return limiter.limiter
}

// cleanup removes old entries from the limiters map
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, limiter := range rl.limiters {
			if time.Since(limiter.lastSeen) > rl.config.CleanupInterval {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func rateLimitMiddleware(rateLimiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			limiter := rateLimiter.Limiter(ip)

			if !limiter.Allow() {
				log.Warn().Str("ip", ip).Str("path", r.URL.Path).Msg("request blocked by rate limiter")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(limitExceededResponse))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientIP extracts the real client IP from the request
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func SetupRateLimiter(router *mux.Router) {
	rateLimiterConfig := config.LoadRateLimiter()

	rateLimiter := newRateLimiterByIP(&rateLimiterConfig)
	router.Use(rateLimitMiddleware(rateLimiter))

	log.Info().
		Float64("requests_per_second", rateLimiterConfig.RequestsPerSecond).
		Int("burst_size", rateLimiterConfig.BurstSize).
		Dur("cleanup_interval", rateLimiterConfig.CleanupInterval).
		Msg("rate limiter middleware enabled from configuration")
}
