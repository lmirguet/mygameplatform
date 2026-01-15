package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"mygameplatform/internal/auth"
	"mygameplatform/internal/lobby_service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	secret := []byte(os.Getenv("JWT_SIGNING_SECRET"))
	verifier, err := auth.NewHS256Signer(secret, nil)
	if err != nil {
		log.Fatalf("jwt: %v", err)
	}

	h := lobby_service.NewHandler(verifier, logger)

	mux := http.NewServeMux()
	mux.Handle("/api/", h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("lobby-service listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http: %v", err)
	}
}
