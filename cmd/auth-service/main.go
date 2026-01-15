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
	"mygameplatform/internal/httpx"
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

	rateLimitPerMinute := 10
	if v := os.Getenv("LOGIN_RATE_LIMIT_PER_MINUTE"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed <= 0 {
			log.Fatalf("LOGIN_RATE_LIMIT_PER_MINUTE: %v", err)
		}
		rateLimitPerMinute = parsed
	}

	rateLimitWindowSeconds := 60
	if v := os.Getenv("LOGIN_RATE_LIMIT_WINDOW_SECONDS"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil || parsed <= 0 {
			log.Fatalf("LOGIN_RATE_LIMIT_WINDOW_SECONDS: %v", err)
		}
		rateLimitWindowSeconds = parsed
	}

	trustProxyHeaders := false
	if v := os.Getenv("TRUST_PROXY_HEADERS"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("TRUST_PROXY_HEADERS: %v", err)
		}
		trustProxyHeaders = parsed
	}

	h := auth_service.NewHandler(users, signer, logger)
	h.SetLimiter(auth_service.NewRateLimiter(rateLimitPerMinute, time.Duration(rateLimitWindowSeconds)*time.Second))
	h.SetTrustProxyHeaders(trustProxyHeaders)

	mux := http.NewServeMux()
	mux.Handle("/api/", h)

	webDistDir := os.Getenv("WEB_DIST_DIR")
	if webDistDir == "" {
		webDistDir = "web/dist"
	}
	if _, err := os.Stat(webDistDir); err == nil {
		mux.Handle("/", httpx.NewSPA(http.Dir(webDistDir), "index.html"))
	} else {
		logger.Warn("web dist dir not found; static hosting disabled", "WEB_DIST_DIR", webDistDir, "err", err)
		mux.Handle("/", http.NotFoundHandler())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("auth-service listening", "addr", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http: %v", err)
	}
}
