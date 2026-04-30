package services

import "time"

type key_linked struct {
	ID          string
	OriginalURL string
	ShortUrl    string
	createdAt   time.Time
}

type key_pool struct {
	key_code   string
	is_used    bool
	created_at time.Time
}

func GetURL(OriginalURL string) (string, error) {
	//fetch a short key,
	//move that from key_gen db to key_linked db
	//build response and return
	return "", nil
}
