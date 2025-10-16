package transaction

import (
	"fmt"
)

type ErrDuplicateIdempotencyKey struct {
	TransactionID string
}

func (e ErrDuplicateIdempotencyKey) Error() string {
	return fmt.Sprintf("transaction already processed with idempotency key, transactionId: %s", e.TransactionID)
}
