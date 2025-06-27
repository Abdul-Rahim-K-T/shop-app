package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shop-backend/config"
	"shop-backend/internal/server"
	"time"
)

func main() {
	cfg := config.LoadConfig()

	srv := server.NewServer(cfg)

	go func() {
		log.Println("Starting server on port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Sutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server foced to shutdown:", err)
	}
}
