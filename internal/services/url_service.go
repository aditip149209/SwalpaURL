package services

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/aditip149209/SwalpaUrl/internal/repository"
)

// URLService handles URL shortening and redirection logic
type URLService struct {
	keyService *KeyGenerationService
	repo       repository.Repository
}

// NewURLService creates a new URLService instance
func NewURLService(keyService *KeyGenerationService, repo repository.Repository) *URLService {
	return &URLService{
		keyService: keyService,
		repo:       repo,
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
	return us.repo.Getoriginal_url(ctx, shortCode)
}
