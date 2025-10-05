package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// database configuration to connect database with the application
var DB *sql.DB

func InitDB() {
	var err error
	//load the database URL from env.go
	databaseURL := GetEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/money_transfer_db?sslmode=disable")

	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func RunMigrations() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		is_verified BOOLEAN DEFAULT FALSE,
		verification_token VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	log.Println("Database migrations completed successfully")
}
