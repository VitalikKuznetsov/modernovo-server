package server

import (
	sd "Modernovo/GoBackend/WorkingWithTheStore/database"
	sm "Modernovo/GoBackend/WorkingWithTheStore/models"
	um "Modernovo/GoBackend/WorkingWithTheUsers/Models"
	ud "Modernovo/GoBackend/WorkingWithTheUsers/DataBase"
	"Modernovo/GoBackend/weather"
	"Modernovo/container"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func getCurrentUser(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			token = cookie.Value
		}
	}

	if token == "" {
		return "", nil
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		return "", err
	}
	defer userDB.Close()

	email, err := userDB.GetUserByToken(token)
	if err != nil {
		return "", nil
	}

	return email, nil
}

func SetupRoutes(router *mux.Router, staticPath string) {
	fs := http.FileServer(http.Dir(staticPath))

	router.HandleFunc("/api/weather", handleGetWeather).Methods("GET")

	router.HandleFunc("/api/register", handleRegister).Methods("POST")
	router.HandleFunc("/api/login", handleLogin).Methods("POST")
	router.HandleFunc("/api/logout", handleLogout).Methods("POST")
	router.HandleFunc("/api/profile", handleAddInfo).Methods("PUT")
	router.HandleFunc("/api/profile", handleGetInfo).Methods("GET")

	router.HandleFunc("/api/products", handleGetProducts).Methods("GET")
	router.HandleFunc("/api/products/{id}", handleGetProduct).Methods("GET")
	router.HandleFunc("/api/products/{id}/detail", handleGetProductDetail).Methods("GET")

	router.HandleFunc("/api/favorites", handleGetFavorites).Methods("GET")
	router.HandleFunc("/api/favorites", handleAddToFavorites).Methods("POST")
	router.HandleFunc("/api/favorites", handleRemoveFromFavorites).Methods("DELETE")
	router.HandleFunc("/api/favorites/check", handleCheckFavorite).Methods("GET")

	router.HandleFunc("/api/cart", handleGetCart).Methods("GET")
	router.HandleFunc("/api/cart", handleAddToCart).Methods("POST")
	router.HandleFunc("/api/cart", handleUpdateCart).Methods("PUT")
	router.HandleFunc("/api/cart", handleRemoveFromCart).Methods("DELETE")
	router.HandleFunc("/api/cart/clear", handleClearCart).Methods("POST")

	router.HandleFunc("/api/admin/products", handleGetAdminProducts).Methods("GET")
	router.HandleFunc("/api/admin/products", handleCreateProduct).Methods("POST")
	router.HandleFunc("/api/admin/products/{id}", handleUpdateProduct).Methods("PUT")
	router.HandleFunc("/api/admin/products/{id}", handleDeleteProduct).Methods("DELETE")

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

	token, err := mydb.CreateSession(uj.Email)
	if err != nil {
		log.Printf("Session creation failed: %v", err)
		sendError(w, "Login failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	response := um.AuthResponse{
		Token: token,
		Email: uj.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		cookie, err := r.Cookie("session_token")
		if err == nil {
			token = cookie.Value
		}
	}

	if token != "" {
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		mydb, err := ud.ConnectToUserDB(container.STR)
		if err == nil {
			defer mydb.Close()
			mydb.DeleteSession(token)
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
	})

	sendSuccess(w, "Logout successful")
}

func handleAddInfo(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var u um.InfoForUser
	err = json.NewDecoder(r.Body).Decode(&u)
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

	err = mydb.AddInfoToDB(u.Name, u.PhoneNumber, userEmail)
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
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	mydb, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer mydb.Close()

	userInfo, err := mydb.GetInfoOfDB(userEmail)
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

func handleAddToFavorites(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req um.FavoriteRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	err = userDB.AddToFavorites(userEmail, req.ProductID)
	if err != nil {
		if err.Error() == "product already in favorites" {
			sendError(w, "Product already in favorites", http.StatusConflict)
		} else {
			log.Printf("Failed to add to favorites: %v", err)
			sendError(w, "Failed to add to favorites", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Product added to favorites")
}

func handleRemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req um.FavoriteRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	err = userDB.RemoveFromFavorites(userEmail, req.ProductID)
	if err != nil {
		if err.Error() == "product not found in favorites" {
			sendError(w, "Product not found in favorites", http.StatusNotFound)
		} else {
			log.Printf("Failed to remove from favorites: %v", err)
			sendError(w, "Failed to remove from favorites", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Product removed from favorites")
}

func handleGetFavorites(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	productIDs, err := userDB.GetUserFavorites(userEmail)
	if err != nil {
		log.Printf("Failed to get favorites: %v", err)
		sendError(w, "Failed to get favorites", http.StatusInternalServerError)
		return
	}

	storeDB, err := sd.ConnectToStoreDB(container.STR)
	if err != nil {
		log.Printf("Store database connection failed: %v", err)
		sendError(w, "Store database connection failed", http.StatusInternalServerError)
		return
	}
	defer storeDB.Close()

	var favorites []sm.Product
	for _, productID := range productIDs {
		product, err := storeDB.GetProduct(productID)
		if err == nil {
			favorites = append(favorites, product)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

func handleCheckFavorite(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	productIDStr := r.URL.Query().Get("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		sendError(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	isFavorite, err := userDB.IsProductInFavorites(userEmail, productID)
	if err != nil {
		log.Printf("Failed to check favorite status: %v", err)
		sendError(w, "Failed to check favorite status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"isFavorite": isFavorite})
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

func handleGetProductDetail(w http.ResponseWriter, r *http.Request) {
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

	productDetail := sm.ProductDetail{
		Product:        product,
		AdditionalInfo: "Дополнительная информация о товаре",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productDetail)
}

func handleAddToCart(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req um.CartRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		sendError(w, "Quantity must be positive", http.StatusBadRequest)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	err = userDB.AddToCart(userEmail, req.ProductID, req.Quantity)
	if err != nil {
		if err.Error() == "product not found" {
			sendError(w, "Product not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to add to cart: %v", err)
			sendError(w, "Failed to add to cart", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Product added to cart")
}

func handleGetCart(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	cartItems, err := userDB.GetCart(userEmail)
	if err != nil {
		log.Printf("Failed to get cart: %v", err)
		sendError(w, "Failed to get cart", http.StatusInternalServerError)
		return
	}

	var total float64
	var itemCount int

	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
		itemCount += item.Quantity
	}

	response := um.CartResponse{
		Items:     cartItems,
		Total:     total,
		ItemCount: itemCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleUpdateCart(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req um.CartRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	err = userDB.UpdateCartItem(userEmail, req.ProductID, req.Quantity)
	if err != nil {
		log.Printf("Failed to update cart: %v", err)
		sendError(w, "Failed to update cart", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, "Cart updated successfully")
}

func handleRemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	var req um.CartRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	err = userDB.RemoveFromCart(userEmail, req.ProductID)
	if err != nil {
		log.Printf("Failed to remove from cart: %v", err)
		sendError(w, "Failed to remove from cart", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, "Product removed from cart")
}

func handleClearCart(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	userDB, err := ud.ConnectToUserDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer userDB.Close()

	err = userDB.ClearCart(userEmail)
	if err != nil {
		log.Printf("Failed to clear cart: %v", err)
		sendError(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	sendSuccess(w, "Cart cleared successfully")
}

func isAdmin(userEmail string) bool {
	return userEmail == "admin@mail.ru"
}

func handleGetAdminProducts(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !isAdmin(userEmail) {
		sendError(w, "Admin access required", http.StatusForbidden)
		return
	}

	storeDB, err := sd.ConnectToStoreDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer storeDB.Close()

	products, err := storeDB.GetAllProductsAdmin()
	if err != nil {
		log.Printf("Failed to get products: %v", err)
		sendError(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !isAdmin(userEmail) {
		sendError(w, "Admin access required", http.StatusForbidden)
		return
	}

	var req um.AdminProductRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Price <= 0 {
		sendError(w, "Name and price are required", http.StatusBadRequest)
		return
	}

	storeDB, err := sd.ConnectToStoreDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer storeDB.Close()

	productID, err := storeDB.CreateProduct(
		req.Name,
		req.Description,
		req.Price,
		req.ImageURL,
		req.ImageURLs,
	)
	if err != nil {
		log.Printf("Failed to create product: %v", err)
		sendError(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Product created successfully",
		"product_id": productID,
	})
}

func handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !isAdmin(userEmail) {
		sendError(w, "Admin access required", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req um.AdminProductRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Price <= 0 {
		sendError(w, "Name and price are required", http.StatusBadRequest)
		return
	}

	storeDB, err := sd.ConnectToStoreDB(container.STR)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		sendError(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer storeDB.Close()

	err = storeDB.UpdateProduct(
		id,
		req.Name,
		req.Description,
		req.Price,
		req.ImageURL,
		req.ImageURLs,
	)
	if err != nil {
		if err.Error() == "product not found" {
			sendError(w, "Product not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to update product: %v", err)
			sendError(w, "Failed to update product", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Product updated successfully")
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	userEmail, err := getCurrentUser(r)
	if err != nil || userEmail == "" {
		sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	if !isAdmin(userEmail) {
		sendError(w, "Admin access required", http.StatusForbidden)
		return
	}

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

	err = storeDB.DeleteProduct(id)
	if err != nil {
		if err.Error() == "product not found" {
			sendError(w, "Product not found", http.StatusNotFound)
		} else {
			log.Printf("Failed to delete product: %v", err)
			sendError(w, "Failed to delete product", http.StatusInternalServerError)
		}
		return
	}

	sendSuccess(w, "Product deleted successfully")
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

func handleGetWeather(w http.ResponseWriter, r *http.Request) {
	weatherData, err := weather.GetWeatherData("Moscow")
	if err != nil {
		log.Printf("Weather API error: %v", err)
		sendError(w, "Failed to get weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}
