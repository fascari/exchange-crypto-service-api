package testdata

import "time"

const (
	ShortTimeout  = 30 * time.Second
	MediumTimeout = 60 * time.Second

	ConcurrentDeposits    = 20
	ConcurrentWithdrawals = 20
)
