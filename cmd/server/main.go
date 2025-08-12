package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/carelessly3/url-shortener/internal/shortener"
	"github.com/carelessly3/url-shortener/internal/storage"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
	Code     string `json:"code"`
}

func main() {
	// replace module path above with whatever you put in go.mod
	store := storage.NewMemoryStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/shorten", makeShortenHandler(store))
	// redirect handler: pattern "/" catches root + any other paths; we parse path
	mux.HandleFunc("/", makeRedirectHandler(store))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server running at http://localhost%s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")
	_ = srv.Close()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func makeShortenHandler(store *storage.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// only POST allowed
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req shortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// try generate unique code (simple retry loop)
		var code string
		var err error
		for i := 0; i < 5; i++ {
			code, err = shortener.GenerateCode()
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			// attempt to save, error means collision
			if saveErr := store.Save(code, req.URL, nil, false); saveErr == nil {
				break
			}
			code = ""
		}
		if code == "" {
			http.Error(w, "could not generate code, try again", http.StatusInternalServerError)
			return
		}

		resp := shortenResponse{
			ShortURL: fmt.Sprintf("http://localhost:8080/%s", code),
			Code:     code,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func makeRedirectHandler(store *storage.MemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// root path or health or api paths should be ignored here
		path := r.URL.Path
		if path == "/" || path == "/health" || len(path) < 2 || path == "/api/shorten" {
			http.NotFound(w, r)
			return
		}
		code := path[1:] // strip leading '/'
		rec, ok := store.Get(code)
		if !ok {
			http.NotFound(w, r)
			return
		}
		// Check expiry if present
		if rec.ExpiresAt != nil && rec.ExpiresAt.Before(time.Now().UTC()) {
			http.Error(w, "link expired", http.StatusGone)
			return
		}
		// increment click count (best-effort)
		store.IncrementClick(code)
		http.Redirect(w, r, rec.LongURL, http.StatusFound)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}
