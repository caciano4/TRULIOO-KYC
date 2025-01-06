package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	AppLogger.Print("Starting connect DB")

	host := GetEnv("DB_HOST", "db_kyc")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "caciano4")
	password := GetEnv("DB_PASSWORD", "123caciano")
	dbname := GetEnv("DB_NAME", "trullio")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		AppLogger.Fatalf("Unable to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		AppLogger.Fatalf("Unable to ping database: %v", err)
	}

	AppLogger.Println("Successfully connected to the database!")
	return db
}

func CloseConnectionDB(db *sql.DB) {
	defer db.Close()

	AppLogger.Print("Database Closed!")
}
