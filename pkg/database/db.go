package database

import (
	"database/sql"
	"log"
	"pollstream/internal/config"

	_ "github.com/lib/pq"
)

func ConnectToDB(cfg *config.Config) (*sql.DB, error) {

	DB, err := sql.Open("postgres", cfg.DatabaseURL)

	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}

	err = DB.Ping()

	if err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to database")

	return DB, nil

}
