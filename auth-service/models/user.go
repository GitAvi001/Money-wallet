package models

import "time"

// user struct defined for user details table
type User struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Password          string    `json:"-"`
	IsVerified        bool      `json:"is_verified"`
	VerificationToken string    `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Registration takes Name, Email and Passowrd
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login takes Email and Password
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse returns Token and User details
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type SendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}
