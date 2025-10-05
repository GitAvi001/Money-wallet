package utils

import "golang.org/x/crypto/bcrypt" //bcrypt package for password hashing

// Hashes the password for imrove security
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compares the hashed password with plain text password for matching the user
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
