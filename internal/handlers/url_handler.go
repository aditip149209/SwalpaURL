package handlers

import "net/http"

type shortenRequest struct {
	URL string
}

type shortenResponse struct {
	ShortCode   string
	OriginalURL string
	Success     bool
}

func Shorten(w http.ResponseWriter, r *http.Request) {

}

func GetURL(w http.ResponseWriter, r *http.Request) {

}
