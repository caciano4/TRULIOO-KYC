package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	host := GetEnv("DB_HOST", "db_kyc")
	port := GetEnv("DB_PORT", "5432")
	user := GetEnv("DB_USER", "caciano4")
	password := GetEnv("DB_PASSWORD", "123caciano")
	dbname := GetEnv("DB_NAME", "trullio")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Successfully connected to the database!")
	return db
}
