package services

import (
	"context"
	"fmt"
	"log"

	"github.com/aditip149209/SwalpaUrl/internal/repository"
	"github.com/aditip149209/SwalpaUrl/pkg/wacky"
)

// KeyGenerationService manages pre-generation and distribution of unique wacky keys
type KeyGenerationService struct {
	repo      repository.Repository
	generator *wacky.WordLists
	poolSize  int
}

// NewKeyGenerationService creates a new KeyGenerationService instance
func NewKeyGenerationService(repo repository.Repository, poolSize int) *KeyGenerationService {
	return &KeyGenerationService{
		repo:     repo,
		poolSize: poolSize,
	}
}

// Initialize loads word lists and pre-generates keys into the database
// This should be called once at server startup
func (ks *KeyGenerationService) Initialize(ctx context.Context) error {
	log.Printf("Initializing KeyGenerationService with pool size: %d", ks.poolSize)

	// Step 1: Load word lists
	generator, err := wacky.Load()
	if err != nil {
		return fmt.Errorf("failed to load word lists: %w", err)
	}
	ks.generator = generator
	log.Printf("Loaded %d adjectives and %d nouns", len(generator.Adjectives), len(generator.Nouns))

	// Step 2: Generate batch of unique keys
	log.Printf("Generating %d unique keys...", ks.poolSize)
	keys, err := generator.GenerateBatch(ks.poolSize)
	if err != nil {
		return fmt.Errorf("failed to generate key batch: %w", err)
	}
	log.Printf("Successfully generated %d unique keys", len(keys))

	// Step 3: Populate database
	log.Printf("Inserting keys into key_pool table...")
	err = ks.repo.Initkey_pool(ctx, keys)
	if err != nil {
		return fmt.Errorf("failed to initialize key_pool: %w", err)
	}
	log.Printf("Successfully populated key_pool with %d keys", len(keys))

	return nil
}

// GetNextKey fetches the next available key from the pool and marks it as used
// This is an atomic operation to prevent race conditions
func (ks *KeyGenerationService) GetNextKey(ctx context.Context, original_url string) (string, error) {
	key, err := ks.repo.GetNextAvailableKey(ctx, original_url)
	if err != nil {
		return "", fmt.Errorf("failed to get next key: %w", err)
	}
	log.Printf("Allocated key: %s", key)
	return key, nil
}
