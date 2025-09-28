package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func StartServer() {
	router := mux.NewRouter()

	staticPath, err := getStaticPath()
	if err != nil {
		log.Fatal("Failed to get static path:", err)
	}

	SetupRoutes(router, staticPath)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: router,
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
