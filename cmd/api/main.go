package main

import (
	"github.com/Royal17x/subscription-service/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("app port - %s", cfg.App.Port)
}
