package model

type PaymentRequest struct {
	// customer detail
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`

	// transaction detail
	TransactionID string `json:"transaction_id" validate:"required"`
	Amount        int    `json:"amount" validate:"required"`
}
