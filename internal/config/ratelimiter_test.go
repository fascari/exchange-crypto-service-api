package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDefaultRateLimiterConfig(t *testing.T) {
	config := defaultRateLimiterConfig()

	require.Equal(t, float64(10), config.RequestsPerSecond)
	require.Equal(t, 20, config.BurstSize)
	require.Equal(t, time.Minute*5, config.CleanupInterval)
}
