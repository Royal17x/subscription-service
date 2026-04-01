package main

import (
	"context"
	"github.com/Royal17x/subscription-service/internal/config"
	"github.com/Royal17x/subscription-service/internal/db"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.RunMigrations(cfg.DB.DSN()); err != nil {
		log.Fatalf("db.RunMigrations: %v", err)
	}

	pool, err := db.NewPool(context.Background(), cfg.DB.DSN())
	if err != nil {
		log.Fatalf("db.NewPool: %v", err)
	}
	defer pool.Close()

	log.Printf("starting on app port - %s", cfg.App.Port)
}
