package main

import (
	"go_crud_postgres/internal/config"
	"go_crud_postgres/internal/database"
	"log"
)

func main() {
	config := config.LoadConfig()
	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("Failed to connect to datbase:", err)
	}
	defer db.Close()

}
