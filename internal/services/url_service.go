package services

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/aditip149209/SwalpaUrl/internal/repository"
	"github.com/redis/go-redis/v9"
)

// URLService handles URL shortening and redirection logic
type URLService struct {
	keyService *KeyGenerationService
	repo       repository.Repository
	rdb        *redis.Client
}

// NewURLService creates a new URLService instance
func NewURLService(keyService *KeyGenerationService, repo repository.Repository, rdb *redis.Client) *URLService {
	return &URLService{
		keyService: keyService,
		repo:       repo,
		rdb:        rdb,
	}
}

// GetURL shortens an original URL by:
// 1. Fetching the next available key from the key pool
// 2. Creating a mapping between the short code and original URL
// 3. Returning the short code
func (us *URLService) GetURL(ctx context.Context, original_url string) (string, error) {
	if err := validateURL(original_url); err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// This single call handles matching, marking used, AND inserting into keylink!
	shortCode, err := us.keyService.GetNextKey(ctx, original_url)
	if err != nil {
		return "", fmt.Errorf("failed to allocate short code: %w", err)
	}

	log.Printf("Successfully shortened URL: %s -> %s", shortCode, original_url)
	return shortCode, nil
}

func validateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	_, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	return nil
}

func (us *URLService) Getoriginal_url(ctx context.Context, shortCode string) (string, error) {
	cachedUrl, err := us.rdb.Get(ctx, shortCode).Result()

	if err == nil {
		log.Printf("Redis cache hit for short url: %s", shortCode)
		return cachedUrl, nil
	}

	log.Printf("Checking postgres for original url for %s", shortCode)

	originalUrl, err := us.repo.Getoriginal_url(ctx, shortCode)
	if err != nil {
		return "", err
	}

	err = us.rdb.Set(ctx, shortCode, originalUrl, 24*time.Hour).Err()
	if err != nil {
		log.Printf("Failed to save original url to redis")
	} else {
		log.Printf("Was able to save url to redis")
	}

	return originalUrl, err

}
