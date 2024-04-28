package models

import "time"

type User struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Password        string    `json:"password"`
	EmailVerifiedAt time.Time `json:"emailVerified_at"`
	PhoneVerifiedAt time.Time `json:"phoneVerified_at"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}
