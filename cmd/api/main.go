package main

import (
	"fmt"
	"log"
	"net/http"
	"pollstream/internal/config"
	handler "pollstream/internal/config/http"
	"pollstream/internal/config/poll"
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

	pollRepo := poll.NewPollRepository(db)
	pollService := poll.NewPollService(pollRepo)
	pollHandler := handler.NewPollHandler(pollService)

	mux := http.NewServeMux()

	mux.HandleFunc("/poll/create", pollHandler.Createpoll)
	mux.HandleFunc("/poll/getById", pollHandler.GetPollByID)

	srv := http.Server{
		Addr:    ":" + cfg.ServerAddress,
		Handler: mux,
	}

	log.Printf("Server started on %s", cfg.ServerAddress)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Error starting server: %v", err)
	}
}
