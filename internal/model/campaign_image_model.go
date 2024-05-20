package model

import (
	"mime/multipart"
	"time"
)

// bisa create campaign image hanya satu-satu
// tapi hapus campaign image hanya satu-satu
// semua berdasarkan campaign_id
// kalau mau create atau hapus banyak. bisa dari client, dengan mengirimkan request sebanyak jumlah gambar yang diinginkan

type CampaignImageResponse struct {
	ID         string    `json:"id"`
	CampaignID string    `json:"campaign_id"`
	FileName   string    `json:"file_name"`
	IsPrimary  int       `json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateCampaignImageRequest struct {
	CampaignID string                `form:"campaign_id" validate:"required,max=100,uuid"`
	UserID     string                `json:"-" validate:"required,max=100,uuid"` // current user id
	FileName   string                `json:"-"`
	FileImage  *multipart.FileHeader `form:"file_image" validate:"omitempty"`
	IsPrimary  bool                  `form:"is_primary" validate:"required,boolean"`
}
