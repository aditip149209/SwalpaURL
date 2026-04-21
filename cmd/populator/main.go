package main

import (
	"log"

	"github.com/aditip149209/SwalpaUrl/pkg/kgs"
)

func main() {

	pathAdj := "/home/aditi/Documents/SwalpaURL/private/adjectives.txt"
	pathNoun := "/home/aditi/Documents/SwalpaURL/private/nouns.txt"

	adjectives, err := kgs.LoadWords(pathAdj)
	if err != nil {
		log.Fatalf("Could not load words: %s", err)
	}
	nouns, err := kgs.LoadWords(pathNoun)
	if err != nil {
		log.Fatalf("Could not load words: %s", err)
	}

	var keysToInsert []string

	for i := 0; i < 1000; i++ {
		key, _ := kgs.GenerateWackyName(adjectives, nouns)
		keysToInsert = append(keysToInsert, key)

	}

}
