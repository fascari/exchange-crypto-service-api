package createtransaction

type (
	InputPayload struct {
		Amount         float64 `json:"amount" validate:"required,gt=0"`
		IdempotencyKey string  `json:"idempotencyKey" validate:"max=255"`
	}

	OutputPayload struct {
		TransactionID string `json:"transactionId"`
	}
)
