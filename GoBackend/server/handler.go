package server

import (
	sd "Modernovo/GoBackend/WorkingWithTheStore/database"
	sm "Modernovo/GoBackend/WorkingWithTheStore/models"
	ud "Modernovo/GoBackend/WorkingWithTheUsers/database"
	um "Modernovo/GoBackend/WorkingWithTheUsers/models"
	"Modernovo/container"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, staticPath string) {
	fs := http.FileServer(http.Dir(staticPath))

	router.HandleFunc("/api/register", handleRegister).Methods("POST")
	router.HandleFunc("/api/login", handleLogin).Methods("POST")
	router.HandleFunc("/api/profile", handleAddInfo).Methods("PUT")
	router.HandleFunc("/api/profile", handleGetInfo).Methods("GET")

	router.HandleFunc("/api/products", handleGetProducts).Methods("GET")
	router.HandleFunc("/api/products/{id}", handleGetProduct).Methods("GET")

	router.PathPrefix("/").Handler(fs)

	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/",
		http.FileServer(http.Dir(filepath.Join(staticPath, "images")))))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var uj um.UserRegOrLog
	err := json.NewDecoder(r.Body).Decode(&uj)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if uj.Email == "" || uj.Password == "" {
		sendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	mydb, err := ud.ConnectToUserDB(container.STR)
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
	var uj um.UserRegOrLog
	err := json.NewDecoder(r.Body).Decode(&uj)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if uj.Email == "" || uj.Password == "" {
		sendError(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	mydb, err := ud.ConnectToUserDB(container.STR)
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
	var u um.InfoForUser
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	mydb, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mydb.Close()

	err = mydb.AddInfoToDB(u.Name, u.PhoneNumber, u.Email)
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

	mydb, err := ud.ConnectToUserDB(container.STR)
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

func handleGetProducts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	storeDB, err := sd.ConnectToStoreDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer storeDB.Close()

	products, err := storeDB.GetAllProducts(limit, offset)
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		sendError(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	total, err := storeDB.GetProductsCount()
	if err != nil {
		log.Printf("Failed to get products count: %v", err)
		sendError(w, "Failed to get products count", http.StatusInternalServerError)
		return
	}

	response := sm.ProductList{
		Products: products,
		Total:    total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	storeDB, err := sd.ConnectToStoreDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer storeDB.Close()

	product, err := storeDB.GetProduct(id)
	if err != nil {
		if err.Error() == "product not found" {
			sendError(w, "Product not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to get product: %v", err)
			sendError(w, "Failed to get product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

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
