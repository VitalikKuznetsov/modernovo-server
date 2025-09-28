package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func StartServer() {
	router := mux.NewRouter()

	staticPath, err := getStaticPath()
	if err != nil {
		log.Fatal("Failed to get static path:", err)
	}

	SetupRoutes(router, staticPath)

	handler := corsMiddleware(router)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func getStaticPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	staticPath := filepath.Join(wd, "..", "static")

	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		return "", err
	}

	indexPath := filepath.Join(staticPath, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return "", err
	}

	return staticPath, nil
}
