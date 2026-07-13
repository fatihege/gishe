package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s environment variable not set", key)
	}

	return value
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if _, err := fmt.Fprintf(w, "OK %s", r.URL.Path); err != nil {
		log.Printf("response write error: %v", err)
	}
}

func main() {
	dbUrl := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			mustEnv("POSTGRES_USER"),
			mustEnv("POSTGRES_PASSWORD"),
		),
		Host: net.JoinHostPort(
			mustEnv("POSTGRES_HOST"),
			mustEnv("POSTGRES_PORT"),
		),
		Path: mustEnv("POSTGRES_DB"),
	}

	query := dbUrl.Query()
	query.Add("sslmode", "disable")
	dbUrl.RawQuery = query.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbUrl.String())
	if err != nil {
		log.Fatalf("create database pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping database: %v", err)
	}

	log.Println("connected to postgres")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("server listening on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("http server error:", err)
	}
}
