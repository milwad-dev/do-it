package models

import (
	"time"
)

type User struct {
	ID              int       `json:"id" validate:"required"`
	Name            string    `json:"name" validate:"required,min=3,max=250"`
	Email           string    `json:"email" validate:"required,email,min=3,max=250"`
	Phone           string    `json:"phone" validate:"required,len=11,unique"`
	Password        string    `json:"-" validate:"required,min=8,max=250"`
	EmailVerifiedAt time.Time `json:"emailVerified_at,omitempty"`
	PhoneVerifiedAt time.Time `json:"phoneVerified_at,omitempty"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

type Label struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Color     string `json:"color"`
	UserId    int    `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	User `json:"user"`
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	LabelId     int    `json:"label_id"`
	UserId      int    `json:"user_id"`
	CompletedAt string `json:"completed_at,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`

	User  `json:"user"`
	Label `json:"label"`
}
