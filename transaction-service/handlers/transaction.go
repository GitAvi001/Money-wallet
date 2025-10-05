package handlers

import (
	"database/sql"
	"net/http"
	"transaction-service/config"
	"transaction-service/models"

	"github.com/gin-gonic/gin"
)

func Transfer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	senderID := userID.(int)

	// Prevent self-transfer
	if senderID == req.ReceiverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to yourself"})
		return
	}

	// Start transaction
	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Check sender balance
	var senderBalance float64
	err = tx.QueryRow("SELECT balance FROM wallets WHERE user_id = $1 FOR UPDATE", senderID).Scan(&senderBalance)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sender wallet not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sender wallet"})
		return
	}

	//validate sender balance below the current balance
	if senderBalance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	// Check if receiver wallet exists, create if not
	var receiverWalletExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM wallets WHERE user_id = $1)", req.ReceiverID).Scan(&receiverWalletExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check receiver wallet"})
		return
	}

	if !receiverWalletExists {
		_, err = tx.Exec("INSERT INTO wallets (user_id, balance) VALUES ($1, 0.00)", req.ReceiverID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create receiver wallet"})
			return
		}
	}

	// Deduct from sender
	_, err = tx.Exec("UPDATE wallets SET balance = balance - $1, updated_at = NOW() WHERE user_id = $2", req.Amount, senderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deduct from sender"})
		return
	}

	// Add the new fund to receiver
	_, err = tx.Exec("UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE user_id = $2", req.Amount, req.ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to receiver"})
		return
	}

	// Transaction details records to transactions table
	var transaction models.Transaction
	err = tx.QueryRow(
		`INSERT INTO transactions (sender_id, receiver_id, amount, status, transaction_type, description) 
		 VALUES ($1, $2, $3, $4, $5, $6) 
		 RETURNING id, sender_id, receiver_id, amount, status, description, transaction_type, created_at, updated_at`,
		senderID, req.ReceiverID, req.Amount, "completed", "transfer", req.Description,
	).Scan(
		&transaction.ID, &transaction.SenderID, &transaction.ReceiverID,
		&transaction.Amount, &transaction.Status, &transaction.Description,
		&transaction.TransactionType, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record transaction"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, models.TransactionResponse{
		Transaction: transaction,
		Message:     "Transfer completed successfully",
	})
}

func GetTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	rows, err := config.DB.Query(
		`SELECT id, sender_id, receiver_id, amount, status, description, transaction_type, created_at, updated_at 
		 FROM transactions 
		 WHERE sender_id = $1 OR receiver_id = $1 
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID, &t.SenderID, &t.ReceiverID, &t.Amount,
			&t.Status, &t.Description, &t.TransactionType,
			&t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan transaction"})
			return
		}
		transactions = append(transactions, t)
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	c.JSON(http.StatusOK, transactions)
}

func GetTransactionByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	transactionID := c.Param("id")

	var transaction models.Transaction
	err := config.DB.QueryRow(
		`SELECT id, sender_id, receiver_id, amount, status, description, transaction_type, created_at, updated_at 
		 FROM transactions 
		 WHERE id = $1 AND (sender_id = $2 OR receiver_id = $2)`,
		transactionID, userID,
	).Scan(
		&transaction.ID, &transaction.SenderID, &transaction.ReceiverID,
		&transaction.Amount, &transaction.Status, &transaction.Description,
		&transaction.TransactionType, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}
