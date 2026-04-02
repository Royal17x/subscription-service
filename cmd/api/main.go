package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Royal17x/subscription-service/internal/config"
	"github.com/Royal17x/subscription-service/internal/db"
	"github.com/Royal17x/subscription-service/internal/handler"
	"github.com/Royal17x/subscription-service/internal/service"
	"github.com/Royal17x/subscription-service/internal/storage/postgres"
)

//	@title			Subscription Service API
//	@version		1.0
//	@description	REST сервис для агрегации данных об онлайн подписках пользователей

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("config.Load", "error", err)
		os.Exit(1)
	}

	if err := db.RunMigrations(cfg.DB.DSN()); err != nil {
		log.Error("db.RunMigrations", "error", err)
		os.Exit(1)
	}

	pool, err := db.NewPool(context.Background(), cfg.DB.DSN())
	if err != nil {
		log.Error("db.NewPool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	repo := postgres.NewSubscriptionRepository(pool)
	svc := service.NewSubscriptionService(repo)
	h := handler.New(svc, log)

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: h.Router(),
	}

	go func() {
		log.Info("starting server", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...", "port", cfg.App.Port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server shutdown", "error", err)
	}
	log.Info("server stopped")
}
