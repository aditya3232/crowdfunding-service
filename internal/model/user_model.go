package model

import (
	"mime/multipart"
	"time"
)

type UserResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Occupation string    `json:"occupation"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RegisterUserRequest struct {
	Name       string `json:"name" validate:"required,max=255"`
	Occupation string `json:"occupation" validate:"required,max=255"`
	Email      string `json:"email" validate:"required,email,max=255"`
	Password   string `json:"password" validate:"required,max=255"`
	Role       string `json:"role" validate:"required,max=255"`
}

type UpdateUserRequest struct {
	ID         string `json:"-" validate:"required,max=100,uuid"`
	Name       string `json:"name" validate:"required,max=255"`
	Occupation string `json:"occupation" validate:"required,max=255"`
	Email      string `json:"email" validate:"required,email,max=255"`
	Password   string `json:"password" validate:"required,max=255"`
	Role       string `json:"role" validate:"required,max=255"`
}

type UpdateAvatarRequest struct {
	ID     string                `json:"-" validate:"required,max=100,uuid"`
	Avatar *multipart.FileHeader `json:"-" validate:"omitempty"`
}

type SearchUserRequest struct {
	Name  string `json:"name" validate:"max=255"`
	Email string `json:"email" validate:"max=255"`
	Role  string `json:"role" validate:"max=255"`
	Page  int    `json:"page" validate:"min=1"`
	Size  int    `json:"size" validate:"min=1,max=100"`
}

type GetUserRequest struct {
	ID string `json:"-" validate:"required,max=100,uuid"`
}

type GetUserByEmailRequest struct {
	Email string `json:"-" validate:"required,email,max=255"`
}

type DeleteUserRequest struct {
	ID string `json:"-" validate:"required,max=100,uuid"`
}
