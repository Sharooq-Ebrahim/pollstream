package main

import (
	"fmt"
	"log"
	"net/http"
	"pollstream/internal/api"
	"pollstream/internal/config"
	"pollstream/internal/poll"
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

	hub := poll.NewHub()
	go hub.Run()

	pollRepo := poll.NewPollRepository(db)
	pollService := poll.NewPollService(pollRepo, hub)
	pollHandler := api.NewPollHandler(pollService)
	wsHandler := api.NewWSHandler(hub)

	mux := http.NewServeMux()

	mux.HandleFunc("/poll/create", pollHandler.CreatePoll)
	mux.HandleFunc("/poll/getById", pollHandler.GetPollByID)
	mux.HandleFunc("/ws/poll", wsHandler.HandleWS)

	srv := http.Server{
		Addr:    ":" + cfg.ServerAddress,
		Handler: mux,
	}

	log.Printf("Server started on %s", cfg.ServerAddress)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Error starting server: %v", err)
	}
}
