package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/channinghe/waken/internal/config"
	"github.com/channinghe/waken/internal/database"
	"github.com/channinghe/waken/internal/repository"
	"github.com/channinghe/waken/internal/server"
)

func main() {
	cfg := config.Load()

	if cfg.AuthToken == "" {
		log.Println("WARNING: WOL_AUTH_TOKEN not set, API authentication is disabled")
	}

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	repo := repository.NewDeviceRepository(db)

	var frontendFS fs.FS
	frontendFS = getFrontendFS()

	router := server.NewRouter(cfg, repo, frontendFS)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
