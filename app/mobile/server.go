package mobile

import (
	"fmt"
	"net/http"

	"github.com/alfon/pokemon-app/app"
)

// NewServer creates an HTTP server with all REST routes registered.
func NewServer(a *app.App, port int) *http.Server {
	mux := http.NewServeMux()
	RegisterRoutes(mux, a)
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: corsMiddleware(mux),
	}
}

// corsMiddleware adds CORS headers for WebView/localhost access.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
