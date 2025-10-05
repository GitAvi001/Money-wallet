package models

import "time"

// Wallet struct for database table wallets
type Wallet struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddFundsRequest struct for adding funds to wallet
type AddFundsRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// WithdrawFundsRequest struct for withdrawing funds from wallet
type WithdrawFundsRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}
