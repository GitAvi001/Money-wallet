package utils

import (
	"auth-service/config"
	"errors"
	"time" //time library for Go using for handle token expirations

	"github.com/golang-jwt/jwt/v5" //JWT library for Go
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Create JWT token for given user ID and email
func GenerateJWT(userID int, email string) (string, error) {
	jwtSecret := config.GetEnv("JWT_SECRET", "your-secret-key")

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			//jwt expires within 24 hours
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	jwtSecret := config.GetEnv("JWT_SECRET", "your-secret-key")

	//Parse secret key for sign the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
