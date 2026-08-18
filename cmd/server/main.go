package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stazoloto/calendar/internal/config"
	"github.com/stazoloto/calendar/internal/repository/psql"
	"github.com/stazoloto/calendar/internal/service"
	"github.com/stazoloto/calendar/internal/transport/rest"
	"github.com/stazoloto/calendar/pkg/database"
)

func main() {
	// подключение конфига
	cfg := config.Load()

	// подключение БД
	db, err := database.NewPostgresConnection(database.DBConfig{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	repo := psql.NewEvents(db)
	svc := service.NewEvents(repo)
	handler := rest.NewHandler(svc)

	srv := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	// запуск сервера
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
