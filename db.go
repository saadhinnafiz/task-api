package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func connectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("unable to connect to database:", err)
	}

	db = pool
	log.Println("connected to database")
}