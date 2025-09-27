package server

import (
	"Modernovo/GoBackend/WorkingWithTheUsers/database"
	"Modernovo/GoBackend/WorkingWithTheUsers/models"
	"Modernovo/container"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, staticPath string) {
	fs := http.FileServer(http.Dir(staticPath))

	// Исправляем названия функций - должны совпадать с HTML
	router.HandleFunc("/api/register", handleRegister).Methods("POST")
	router.HandleFunc("/api/login", handleLogin).Methods("POST")
	router.HandleFunc("/api/profile", handleAddInfo).Methods("PUT")
	router.HandleFunc("/api/profile", handleGetInfo).Methods("GET")

	router.PathPrefix("/").Handler(fs)
}

// Переименовываем функции для соответствия с HTML
func handleRegister(w http.ResponseWriter, r *http.Request) {
	var uj models.UserRegOrLog
	err := json.NewDecoder(r.Body).Decode(&uj)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация
	if uj.Email == "" || uj.Password == "" {
		sendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	mydb, err := database.ConnectToMyDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mydb.Close()

	err = mydb.AddUser(uj.Email, uj.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			sendError(w, "User with this email already exists", http.StatusConflict)
		} else {
			log.Printf("Registration failed: %v", err)
			sendError(w, "Registration failed", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Registration successful")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var uj models.UserRegOrLog
	err := json.NewDecoder(r.Body).Decode(&uj)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Валидация
	if uj.Email == "" || uj.Password == "" {
		sendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	mydb, err := database.ConnectToMyDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mydb.Close()

	err = mydb.Login(uj.Email, uj.Password)
	if err != nil {
		if err.Error() == "user not found" {
			sendError(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			log.Printf("Login failed: %v", err)
			sendError(w, "Login failed", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Login successful")
}

func handleAddInfo(w http.ResponseWriter, r *http.Request) {
	var u models.InfoForUser
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	mydb, err := database.ConnectToMyDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mydb.Close()

	err = mydb.AddInfoToDB(u.Name, u.Surname, u.PhoneNumber, u.DateOfBirth, u.Email)
	if err != nil {
		if err.Error() == "user not found" {
			sendError(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Profile update failed: %v", err)
			sendError(w, "Failed to update profile", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Profile updated successfully")
}

func handleGetInfo(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		sendError(w, "Email parameter is required", http.StatusBadRequest)
		return
	}

	mydb, err := database.ConnectToMyDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mydb.Close()

	userInfo, err := mydb.GetInfoOfDB(email)
	if err != nil {
		if err.Error() == "user not found" {
			sendError(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("Get user info failed: %v", err)
			sendError(w, "Failed to get user info", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

// Вспомогательные функции
func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func sendSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}
