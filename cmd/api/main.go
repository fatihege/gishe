package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fatihege/gishe/internal/auth"
	authhttp "github.com/fatihege/gishe/internal/auth/http"
	authpostgres "github.com/fatihege/gishe/internal/auth/postgres"
	"github.com/fatihege/gishe/internal/config"
	"github.com/fatihege/gishe/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := fmt.Fprintf(w, "OK %s", r.URL.Path); err != nil {
		log.Printf("response write error: %v", err)
	}
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer pool.Close()

	log.Println("connected to postgres")

	authRepository := authpostgres.New(pool)
	passwordHasher := auth.NewPasswordHasher()
	authService := auth.NewService(authRepository, passwordHasher)
	authHandler := authhttp.NewHandler(authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", health)
	r.Post("/auth/register", authHandler.Register)

	server := &http.Server{
		Addr:              cfg.HTTPAddress,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server listen and serve: %v", err)
	}
}
