package handlers

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register a new user
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	err := config.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Generate verification token
	verificationToken := uuid.New().String()

	// Insert user
	var userID int
	err = config.DB.QueryRow(
		"INSERT INTO users (name, email, password, verification_token, is_verified) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		req.Name, req.Email, hashedPassword, verificationToken, false,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send verification email
	go utils.SendVerificationEmail(req.Email, verificationToken)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email to verify your account.",
		"user_id": userID,
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email filtered from the users table
	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, name, email, password, is_verified FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if email is verified
	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please verify your email before logging in"})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	// Update user verification status
	result, err := config.DB.Exec(
		"UPDATE users SET is_verified = true, verification_token = null WHERE verification_token = $1",
		token,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	//If the user is not found
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func SendVerificationEmail(c *gin.Context) {
	var req models.SendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user and check if already verified
	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, is_verified, verification_token FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.IsVerified, &user.VerificationToken)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if user.IsVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already verified"})
		return
	}

	// Generate new verification token if needed
	if user.VerificationToken == "" {
		user.VerificationToken = uuid.New().String()
		_, err = config.DB.Exec(
			"UPDATE users SET verification_token = $1 WHERE id = $2",
			user.VerificationToken, user.ID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
			return
		}
	}

	// Send verification email
	go utils.SendVerificationEmail(req.Email, user.VerificationToken)

	c.JSON(http.StatusOK, gin.H{"message": "Verification email sent successfully"})
}

// Fetching current verified user details
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	err := config.DB.QueryRow(
		"SELECT id, name, email, is_verified, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.IsVerified, &user.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Fetching all verified users details
func GetAllUsers(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rows, err := config.DB.Query(
		"SELECT id, name, email FROM users WHERE id != $1 AND is_verified = true",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
