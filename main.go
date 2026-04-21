package main

import (
	"fmt"
	"log"

	"github.com/aditip149209/SwalpaUrl/pkg/storage"
)

func main() {
	// Usually, you'd load this from a .env file
	// Format: username:password@tcp(127.0.0.1:3306)/dbname?parseTime=true
	dsn := "aditi:root@tcp(127.0.0.1:3306)/swalpaurl_db?parseTime=true"

	store, err := storage.NewStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	fmt.Println("Successfully connected to the database!")
}
