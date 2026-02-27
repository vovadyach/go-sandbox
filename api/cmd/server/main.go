package main

import (
	"log"

	"go-sandbox/api/internal/database"

	"go-sandbox/api/internal/config"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")
}
