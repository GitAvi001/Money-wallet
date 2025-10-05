package models

import "time"

// database table transactions
type Transaction struct {
	ID              int       `json:"id"`
	SenderID        int       `json:"sender_id"`
	ReceiverID      int       `json:"receiver_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	Description     string    `json:"description"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Defines the request body for transfer
type TransferRequest struct {
	ReceiverID  int     `json:"receiver_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

// Defines the response body for transfer
type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
	Message     string      `json:"message"`
}
