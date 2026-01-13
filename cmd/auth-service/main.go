package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"mygameplatform/internal/auth"
	"mygameplatform/internal/auth_service"
	"mygameplatform/internal/db"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	databaseURL := os.Getenv("DATABASE_URL")
	pool, err := db.NewPGPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	secret := []byte(os.Getenv("JWT_SIGNING_SECRET"))

	var ttl *time.Duration
	if v := os.Getenv("JWT_TTL_SECONDS"); v != "" {
		secs, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("JWT_TTL_SECONDS: %v", err)
		}
		d := time.Duration(secs) * time.Second
		ttl = &d
	}

	signer, err := auth.NewHS256Signer(secret, ttl)
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	users := db.NewPostgresUserStore(pool)
	h := auth_service.NewHandler(users, signer, logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("auth-service listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http: %v", err)
	}
}
