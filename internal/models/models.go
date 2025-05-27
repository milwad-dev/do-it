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
	Title     string `json:"title" validate:"required,min=3,max=250"`
	Color     string `json:"color" validate:"required,min=3,max=250"`
	UserId    int    `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	User `json:"user"  validate:"-"`
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required,min=3,max=250"`
	Description string `json:"description" validate:"required,min=3,max=250"`
	Status      string `json:"status" validate:"required,oneof=pending active inactive"`
	LabelId     int    `json:"label_id" validate:"required"`
	UserId      int    `json:"user_id"`
	CompletedAt string `json:"completed_at,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`

	User  `json:"user" validate:"-"`
	Label `json:"label" validate:"-"`
}
