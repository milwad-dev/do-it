package models

import "time"

type User struct {
	ID              int
	Name            string
	Email           string
	Phone           string
	Password        string
	EmailVerifiedAt time.Time
	PhoneVerifiedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
