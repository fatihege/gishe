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
	"github.com/fatihege/gishe/internal/catalog"
	cataloghttp "github.com/fatihege/gishe/internal/catalog/http"
	catalogpostgres "github.com/fatihege/gishe/internal/catalog/postgres"
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

	tokenManager := auth.NewTokenManager([]byte(cfg.JWTSecret), cfg.JWTIssuer, cfg.JWTAudience)

	authService := auth.NewService(authRepository, passwordHasher, tokenManager)
	authHandler := authhttp.NewHandler(authService)

	catalogRepository := catalogpostgres.New(pool)
	catalogService := catalog.NewService(catalogRepository)
	catalogHandler := cataloghttp.NewHandler(catalogService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", health)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/catalog", func(r chi.Router) {
		r.Post("/venues", catalogHandler.CreateVenue)
		r.Get("/venues/{id}/sessions", catalogHandler.GetSessionsByVenueID)
		r.Get("/venues", catalogHandler.GetVenues)

		r.Get("/sessions/{id}", catalogHandler.GetSessionByID)
		r.Post("/sessions", catalogHandler.CreateSession)
		r.Get("/sessions", catalogHandler.GetSessions)
	})

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
