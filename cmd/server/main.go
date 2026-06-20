package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/aditip149209/SwalpaUrl/internal/handlers"
	"github.com/aditip149209/SwalpaUrl/internal/repository"
	"github.com/aditip149209/SwalpaUrl/internal/services"
)

const (
	defaultPort     = "8080"
	defaultDSN      = "user=root password=root dbname=postgres host=localhost port=5433 sslmode=disable"
	defaultPoolSize = 50000
	shutdownTimeout = 30 * time.Second
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 60 * time.Second
)

func main() {
	// Setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Get configuration from environment or use defaults
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbDSN := os.Getenv("DATABASE_URL")
	if dbDSN == "" {
		dbDSN = defaultDSN
	}

	poolSize := defaultPoolSize

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize database connection
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Verify database connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to verify database connection: %v", err)
	}
	log.Println("✓ Database connection established")

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize repository
	repo := repository.NewPostgresRepository(db)

	// Initialize KeyGenerationService
	keyService := services.NewKeyGenerationService(repo, poolSize)

	// Get base path for loading word lists
	// basePath := "/home/aditi/Documents/SwalpaURL"

	// Initialize key pool (pre-generate and populate database)
	ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := keyService.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize KeyGenerationService: %v", err)
	}
	log.Println("✓ KeyGenerationService initialized with pre-generated keys")

	log.Println("Init redis")
	redisAddr := os.Getenv("REDIS_URL")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to establish contact with redis instance context %v", err)
	}
	log.Printf("Connected to redis")
	// Initialize URLService
	urlService := services.NewURLService(keyService, repo, rdb)
	log.Println("✓ URLService initialized")

	// Initialize HTTP handlers
	urlHandler := handlers.NewURLHandler(urlService)
	log.Println("✓ HTTP handlers initialized")

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("POST /shorten", urlHandler.Shorten)
	mux.HandleFunc("GET /{shortCode}", urlHandler.GetURL)
	mux.HandleFunc("GET /health", healthHandler)
	// Add this right next to your /shorten registration
	mux.HandleFunc("GET /", urlHandler.HomeHandler)

	log.Println("✓ Routes configured")

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting HTTP server on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Setup graceful shutdown signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-sigChan:
		log.Printf("Received signal: %v, initiating graceful shutdown", sig)
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("✓ Server shutdown complete")
}

// healthHandler responds to health checks (Kubernetes-compatible)
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}
