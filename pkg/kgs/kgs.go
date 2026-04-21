package kgs

// In a production environment, the available_keys table acts as a buffer.
// You don't want to generate a key only when a user asks for it because that creates a bottleneck. Instead, you have a Producer (the Populator)
// that fills the table and a Consumer (the API) that takes them out.

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
)

func LoadWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, scanner.Err()
}

func generateShortHash() (string, error) {
	b := make([]byte, 2)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func getRandomWord(words []string) (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
	if err != nil {
		return "", err
	}

	return words[n.Int64()], nil
}

func GenerateWackyName(adjectives []string, nouns []string) (string, error) {
	adj, err := getRandomWord(adjectives)
	if err != nil {
		return "", err
	}

	noun, err := getRandomWord(nouns)
	if err != nil {
		return "", err
	}

	hash, err := generateShortHash()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s-%s", adj, noun, hash), nil

}
