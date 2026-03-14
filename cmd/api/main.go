package main

import (
	"fmt"
	"log"
	"pollstream/internal/config"
	"pollstream/pkg/database"
)

func main() {

	cfg := config.LoadConfig()
	fmt.Println("Server Address: ", cfg.ServerAddress)
	fmt.Println("Database URL: ", cfg.DatabaseURL)

	db, err := database.ConnectToDB(cfg)

	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}

	defer db.Close()


      




}
