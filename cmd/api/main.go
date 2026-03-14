package main

import (
	"fmt"
	"pollstream/internal/config"
)

func main() {

	cfg := config.LoadConfig()
	fmt.Println("Server Address: ", cfg.ServerAddress)
	fmt.Println("Database URL: ", cfg.DatabaseURL)

}
