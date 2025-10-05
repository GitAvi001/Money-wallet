package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"transaction-service/config"
	"transaction-service/models"

	"github.com/gin-gonic/gin"
)

func GetWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("GetWallet called for user_id: %v (type: %T)", userID, userID)

	var wallet models.Wallet
	err := config.DB.QueryRow(
		"SELECT id, user_id, balance, created_at, updated_at FROM wallets WHERE user_id = $1",
		userID,
	).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)

	log.Printf("Wallet fetch result - error: %v, wallet: %+v", err, wallet)

	if err == sql.ErrNoRows {
		log.Printf("No wallet found for user_id: %v, creating new wallet", userID)
		// Create wallet if doesn't exist
		err = config.DB.QueryRow(
			"INSERT INTO wallets (user_id, balance) VALUES ($1, 0.00) RETURNING id, user_id, balance, created_at, updated_at",
			userID,
		).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wallet"})
		return
	}

	log.Printf("Returning wallet: %+v", wallet)
	c.JSON(http.StatusOK, wallet)
}

func AddFunds(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AddFundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("AddFunds called - user_id: %v, amount: %f", userID, req.Amount)

	// Start transaction
	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback() // Rollback transaction on error

	// Update wallet balance
	var newBalance float64
	err = tx.QueryRow(
		"UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2 RETURNING balance",
		req.Amount, userID,
	).Scan(&newBalance)

	if err == sql.ErrNoRows {
		// Create wallet if doesn't exist
		err = tx.QueryRow(
			"INSERT INTO wallets (user_id, balance) VALUES ($1, $2) RETURNING balance",
			userID, req.Amount,
		).Scan(&newBalance) // Scan new balance into variable
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
		return
	}

	// Record transaction
	_, err = tx.Exec(
		"INSERT INTO transactions (sender_id, receiver_id, amount, status, transaction_type, description) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, userID, req.Amount, "completed", "deposit", "Added funds to wallet",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	log.Printf("Funds added successfully - user_id: %v, new_balance: %f", userID, newBalance)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Funds added successfully",
		"new_balance": newBalance,
	})
}

func WithdrawFunds(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.WithdrawFundsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check current balance
	var currentBalance float64
	err := config.DB.QueryRow("SELECT balance FROM wallets WHERE user_id = $1", userID).Scan(&currentBalance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wallet"})
		return
	}

	if currentBalance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Start transaction
	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Update wallet balance
	var newBalance float64
	err = tx.QueryRow(
		"UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE user_id = $2 RETURNING balance",
		req.Amount, userID,
	).Scan(&newBalance)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet"})
		return
	}

	// Record transaction
	_, err = tx.Exec(
		"INSERT INTO transactions (sender_id, receiver_id, amount, status, transaction_type, description) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, userID, req.Amount, "completed", "withdrawal", "Withdrew funds from wallet",
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Funds withdrawn successfully",
		"new_balance": newBalance,
	})
}
