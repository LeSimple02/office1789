package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Connectdb() {
	// En local, charge .env si présent (ignore l'erreur si absent)
	_ = godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	fmt.Printf("[DB][DEBUG] DB_HOST=%s\n", host)
	fmt.Printf("[DB][DEBUG] DB_PORT=%s\n", port)
	fmt.Printf("[DB][DEBUG] DB_USER=%s\n", user)
	fmt.Printf("[DB][DEBUG] DB_NAME=%s\n", dbname)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Printf("[DB][DEBUG] Connexion string: %s\n", psqlInfo)

	var err error
	db, err = sql.Open("postgres", psqlInfo)

	if err != nil {
		fmt.Printf("[DB][ERROR] sql.Open failed: %v\n", err)
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		fmt.Printf("[DB][ERROR] db.Ping failed: %v\n", err)
		panic(err)
	}
	fmt.Println("[DB][INFO] Successfully connected to Postgres!")
}
