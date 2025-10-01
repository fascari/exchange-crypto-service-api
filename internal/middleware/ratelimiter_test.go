package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"exchange-crypto-service-api/internal/config"

	"github.com/stretchr/testify/require"
)

func TestNewRateLimiterByIP(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.RateLimiterConfig
		expected float64
	}{
		{
			name:     "nil config uses default",
			config:   nil,
			expected: 10,
		},
		{
			name: "custom config",
			config: &config.RateLimiterConfig{
				RequestsPerSecond: 5,
				BurstSize:         10,
				CleanupInterval:   time.Minute,
			},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl := newRateLimiterByIP(tt.config)
			require.NotNil(t, rl)
			require.Equal(t, tt.expected, rl.config.RequestsPerSecond)
		})
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	cfg := &config.RateLimiterConfig{
		RequestsPerSecond: 1,
		BurstSize:         2,
		CleanupInterval:   time.Minute,
	}
	middleware := rateLimitMiddleware(newRateLimiterByIP(cfg))

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	tests := []struct {
		name           string
		requestCount   int
		expectedStatus int
		expectContains string
	}{
		{
			name:           "allow requests within burst limit",
			requestCount:   2,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "block requests exceeding rate limit",
			requestCount:   3,
			expectedStatus: http.StatusTooManyRequests,
			expectContains: "rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedHandler := middleware(handler)

			var lastResponse *httptest.ResponseRecorder
			for i := 0; i < tt.requestCount; i++ {
				req := httptest.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "192.168.1.1:12345"
				lastResponse = httptest.NewRecorder()
				wrappedHandler.ServeHTTP(lastResponse, req)
			}

			require.Equal(t, tt.expectedStatus, lastResponse.Code)
			if tt.expectContains != "" {
				require.Contains(t, lastResponse.Body.String(), tt.expectContains)
			}
		})
	}
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		name          string
		xForwardedFor string
		xRealIP       string
		remoteAddr    string
		expected      string
	}{
		{
			name:          "X-Forwarded-For priority",
			xForwardedFor: "192.168.1.100",
			remoteAddr:    "10.0.0.1:12345",
			expected:      "192.168.1.100",
		},
		{
			name:       "X-Real-IP fallback",
			xRealIP:    "192.168.1.200",
			remoteAddr: "10.0.0.1:12345",
			expected:   "192.168.1.200",
		},
		{
			name:       "RemoteAddr default",
			remoteAddr: "10.0.0.1:12345",
			expected:   "10.0.0.1:12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tt.remoteAddr

			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			ip := clientIP(req)
			require.Equal(t, tt.expected, ip)
		})
	}
}
