package createtransaction

type InputPayload struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}
