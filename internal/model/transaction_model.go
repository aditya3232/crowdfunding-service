package model

import "time"

// response after create transaction
type TransactionResponse struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	UserID     string    `json:"user_id"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       string    `json:"code"`
	PaymentURL string    `json:"payment_url"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetTransactionByCampaignIDResponse struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type GetTransactionByUserIDResponse struct {
	ID        string    `json:"id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Campaign  Campaign  `json:"campaign"` //tidak ambil dari campaign model karena tidak semua field yang dibutuhkan
}

type Campaign struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CreateTransactionRequest struct {
	CampaignID string `json:"campaign_id" validate:"required"`
	UserID     string `json:"user_id" validate:"required"`
	Amount     int    `json:"amount" validate:"required"`
}

// request to create transaction notification (notifikasi pembayaran) yg dikirim dari midtrans ke service kita
type CreateTransactionNotificationRequest struct {
	ID                string `json:"id"` // transaction_id
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

// get all need paginate
// user_id yang melakukan request (yg login)
// tidak perlu json tag
type GetTransactionByCampaignIDRequest struct {
	CampaignID string `json:"-" validate:"required,max=100,uuid"`
	UserID     string `json:"-" validate:"required,max=100,uuid"`
	Page       int    `json:"page" validate:"min=1"`
	Size       int    `json:"size" validate:"min=1,max=100"`
}

// get all need paginate
type GetTransactionByUserIDRequest struct {
	UserID string `json:"-" validate:"required,max=100,uuid"`
	Page   int    `json:"page" validate:"min=1"`
	Size   int    `json:"size" validate:"min=1,max=100"`
}

type GetTransactionRequest struct {
	ID string `json:"-" validate:"required,max=100,uuid"`
}
