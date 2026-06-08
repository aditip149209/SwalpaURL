package wacky

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// WordLists holds the adjectives and nouns needed for wacky name generation
type WordLists struct {
	Adjectives []string
	Nouns      []string
	mu         sync.RWMutex
}

// Load reads adjectives.txt and nouns.txt from the private/ directory
// Returns a WordLists struct ready for name generation
func Load(basePath string) (*WordLists, error) {
	adjPath := filepath.Join(basePath, "private", "adjectives.txt")
	nounPath := filepath.Join(basePath, "private", "nouns.txt")

	adjectives, err := readWordList(adjPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load adjectives: %w", err)
	}

	nouns, err := readWordList(nounPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load nouns: %w", err)
	}

	if len(adjectives) == 0 {
		return nil, fmt.Errorf("adjectives list is empty")
	}
	if len(nouns) == 0 {
		return nil, fmt.Errorf("nouns list is empty")
	}

	return &WordLists{
		Adjectives: adjectives,
		Nouns:      nouns,
	}, nil
}

// readWordList reads a file line-by-line and returns non-empty trimmed lines
func readWordList(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// Generate creates a single random wacky name in format "Adjective-Noun-NNNN"
// where NNNN is a random 4-digit number (1000-9999)
func (wl *WordLists) Generate() string {
	wl.mu.RLock()
	defer wl.mu.RUnlock()

	adjIdx := rand.Intn(len(wl.Adjectives))
	nounIdx := rand.Intn(len(wl.Nouns))
	number := rand.Intn(9000) + 1000 // 1000-9999

	adj := wl.Adjectives[adjIdx]
	noun := wl.Nouns[nounIdx]

	return fmt.Sprintf("%s-%s-%d", adj, noun, number)
}

// GenerateBatch generates count unique wacky names
// Uses a map to track uniqueness and retries on collision
// Returns error if unable to generate enough unique names
func (wl *WordLists) GenerateBatch(count int) ([]string, error) {
	if count <= 0 {
		return []string{}, nil
	}

	seen := make(map[string]bool)
	maxAttempts := count * 10 // Allow 10x attempts to avoid infinite loops
	attempts := 0

	for len(seen) < count && attempts < maxAttempts {
		key := wl.Generate()
		seen[key] = true
		attempts++
	}

	if len(seen) < count {
		return nil, fmt.Errorf(
			"could not generate %d unique names after %d attempts, got %d unique names",
			count, maxAttempts, len(seen),
		)
	}

	// Convert map to slice
	result := make([]string, 0, count)
	for key := range seen {
		result = append(result, key)
	}

	return result, nil
}
