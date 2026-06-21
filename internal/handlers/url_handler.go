package handlers

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aditip149209/SwalpaUrl/internal/services"
)

//go:embed templates/*
var templateFS embed.FS

// ShortenRequest represents the incoming request to shorten a URL
type ShortenRequest struct {
	URL string `json:"url"`
}

// ShortenResponse represents the response when a URL is shortened
type ShortenResponse struct {
	ShortCode   string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	Success     bool   `json:"success"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

// URLHandler holds dependencies for URL-related HTTP handlers
type URLHandler struct {
	urlService *services.URLService
}

// NewURLHandler creates a new URLHandler instance
func NewURLHandler(urlService *services.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// Shorten handles POST /shorten requests to create a shortened URL
func (uh *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "method not allowed, use POST",
			Success: false,
		})
		return
	}

	// Limit request body to 1MB for security defense
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var targetURL string
	// Detect if the request is coming from your HTMX form frontend
	isHTMX := r.Header.Get("HX-Request") == "true"

	// 1. READ INPUT BASED ON CLIENT TYPE
	if isHTMX {
		// HTMX form submissions send data as standard form parameters
		targetURL = r.FormValue("url")
	} else {
		// Postman or native API clients send data as a raw JSON payload
		var req ShortenRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{
				Error:   "invalid JSON format or empty body",
				Success: false,
			})
			return
		}
		targetURL = req.URL
	}

	// 2. VALIDATE THE URL INPUT
	if targetURL == "" {
		if isHTMX {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`<p class="text-red-400 font-medium text-sm">URL field is required</p>`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "url field is required", Success: false})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// 3. EXECUTE THE BACKEND SERVICE CORE LOGIC
	shortCode, err := uh.urlService.GetURL(ctx, targetURL)
	if err != nil {
		log.Printf("Failed to shorten URL: %v", err)
		if isHTMX {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`<p class="text-red-400 font-medium text-sm">Error: %s</p>`, err.Error())))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "failed to shorten URL: " + err.Error(),
			Success: false,
		})
		return
	}

	// 4. RETURN THE OUTPUT COMPONENT OR DATA
	if isHTMX {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		displayURL := baseURL
		if len(displayURL) > 8 && displayURL[:8] == "https://" {
			displayURL = displayURL[8:]
		} else if len(displayURL) > 7 && displayURL[:7] == "http://" {
			displayURL = displayURL[7:]
		}

		// This gorgeous, fully responsive component box replaces the raw JSON text!
		fmt.Fprintf(w, `
            <div class="bg-gray-700/40 border border-gray-600/60 p-5 rounded-lg shadow-inner text-center mt-2 space-y-3">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold bg-green-950 text-green-400 border border-green-500/20">
                    ✓ Shortlink Generated
                </span>
                <p class="text-xs text-gray-400 font-medium break-all px-2">Target: %s</p>
                <div class="flex items-center justify-between bg-gray-900/90 p-3 rounded-lg border border-gray-700/80 gap-3">
                    <a href="%s/%s" target="_blank" class="text-md font-bold text-green-400 hover:text-green-300 hover:underline break-all tracking-wide text-left pl-1">
                        %s/%s
                    </a>
                    <button onclick="navigator.clipboard.writeText('%s/%s')" class="text-xs bg-gray-800 hover:bg-gray-700 text-gray-200 font-semibold px-3 py-2 rounded-md border border-gray-600 active:scale-95 transition-all cursor-pointer shadow">
                        Copy
                    </button>
                </div>
            </div>
        `, targetURL, baseURL, shortCode, displayURL, shortCode, baseURL, shortCode)
		return
	}

	// Fallback JSON block keeps Postman working perfectly!
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ShortenResponse{
		ShortCode:   shortCode,
		OriginalURL: targetURL,
		Success:     true,
	})
}

// GetURL handles GET /{shortCode} requests to redirect to the original URL
func (uh *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	// Enforce GET method
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "method not allowed, use GET",
			Success: false,
		})
		return
	}

	// Extract short code from URL path
	// Assumes route like GET /:shortCode
	shortCode := r.PathValue("shortCode")
	if shortCode == "" {
		// Fallback: try to get from query parameter
		shortCode = r.URL.Query().Get("code")
	}

	if shortCode == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "short code is required",
			Success: false,
		})
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Fetch original URL from service
	original_url, err := uh.urlService.Getoriginal_url(ctx, shortCode)
	if err != nil {
		log.Printf("Failed to get original URL for code %s: %v", shortCode, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "short code not found",
			Success: false,
		})
		return
	}

	// Redirect to original URL
	w.Header().Set("Location", original_url)
	w.WriteHeader(http.StatusFound)
}

// HomeHandler serves the static HTML landing page
func (uh *URLHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}
