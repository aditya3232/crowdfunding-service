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

type CampaignTransactionResponse struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionResponse struct {
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
	Status     string `json:"status" validate:"required"`
}

// request to create transaction notification (notifikasi pembayaran) yg dikirim dari midtrans ke service kita
type CreateTransactionNotificationRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}

type SearchCampaignTransactionRequest struct {
	UserName string `json:"user_name" validate:"max=255"`
}

type SearchUserTransactionRequest struct {
	Status string `json:"status" validate:"max=255"`
}

type GetCampaignTransactionRequest struct {
	CampaignID string `json:"-" validate:"required,max=100,uuid"`
}
