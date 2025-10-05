package config

import (
	"database/sql" //database/sql package for database connection
	"log"

	_ "github.com/lib/pq" //postgres driver for database connection
)

var DB *sql.DB

func InitDB() {
	var err error
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

// RunMigrations will create the required tables in the database
func RunMigrations() {
	createWalletsTable := `
	CREATE TABLE IF NOT EXISTS wallets (
		id SERIAL PRIMARY KEY,
		user_id INTEGER UNIQUE NOT NULL,
		balance DECIMAL(15, 2) DEFAULT 0.00,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	createTransactionsTable := `
	CREATE TABLE IF NOT EXISTS transactions (
		id SERIAL PRIMARY KEY,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		amount DECIMAL(15, 2) NOT NULL, #15 digits before decimal and 2 digits after decimal total integers 13
		status VARCHAR(50) DEFAULT 'pending',
		description TEXT,
		transaction_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(createWalletsTable)
	if err != nil {
		log.Fatal("Failed to create wallets table:", err)
	}

	_, err = DB.Exec(createTransactionsTable)
	if err != nil {
		log.Fatal("Failed to create transactions table:", err)
	}

	log.Println("Database migrations completed successfully")
}
