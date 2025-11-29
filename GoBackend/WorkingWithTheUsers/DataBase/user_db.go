package database

import (
	"Modernovo/GoBackend/WorkingWithTheUsers/models"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type UserDB struct {
	db *sql.DB
}

func (db *UserDB) AddUser(email, password string) error {
	var exists bool
	err := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE Email = $1)", email).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("user already exists")
	}

	sqlRequest := "INSERT INTO users (Email, Password) VALUES ($1, $2)"
	_, err = db.db.Exec(sqlRequest, email, password)
	return err
}

func (db *UserDB) Login(email, password string) error {
	query := "SELECT Email, Password FROM users WHERE Email = $1 AND Password = $2"

	row := db.db.QueryRow(query, email, password)

	user := &models.User{}
	err := row.Scan(
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return errors.New("user not found")
	}

	return err
}

func (db *UserDB) AddInfoToDB(name, phonenumber, email string) error {
	query := "UPDATE users SET Name = $1, PhoneNumber = $2 WHERE Email = $3"
	r, err := db.db.Exec(query, name, phonenumber, email)
	rowsAffected, _ := r.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return err
}

func (db *UserDB) GetInfoOfDB(email string) (models.InfoForUser, error) {
	query := "SELECT Name, PhoneNumber FROM users WHERE Email = $1"
	row := db.db.QueryRow(query, email)
	user := models.InfoForUser{}

	err := row.Scan(
		&user.Name,
		&user.PhoneNumber,
	)
	user.Email = email

	if err == sql.ErrNoRows {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (db *UserDB) Close() {
	db.db.Close()
}

func ConnectToUserDB(connectSring string) (*UserDB, error) {
	db, err := sql.Open("postgres", connectSring)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &UserDB{db: db}, nil
}

func (db *UserDB) AddToFavorites(userEmail string, productID int) error {
	query := "INSERT INTO favorites (user_email, product_id) VALUES ($1, $2) ON CONFLICT (user_email, product_id) DO NOTHING"
	result, err := db.db.Exec(query, userEmail, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product already in favorites")
	}

	return nil
}

func (db *UserDB) RemoveFromFavorites(userEmail string, productID int) error {
	query := "DELETE FROM favorites WHERE user_email = $1 AND product_id = $2"
	result, err := db.db.Exec(query, userEmail, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found in favorites")
	}

	return nil
}

func (db *UserDB) GetUserFavorites(userEmail string) ([]int, error) {
	query := "SELECT product_id FROM favorites WHERE user_email = $1"
	rows, err := db.db.Query(query, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productIDs []int
	for rows.Next() {
		var productID int
		if err := rows.Scan(&productID); err != nil {
			return nil, err
		}
		productIDs = append(productIDs, productID)
	}

	return productIDs, nil
}

func (db *UserDB) IsProductInFavorites(userEmail string, productID int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM favorites WHERE user_email = $1 AND product_id = $2)"
	var exists bool
	err := db.db.QueryRow(query, userEmail, productID).Scan(&exists)
	return exists, err
}

func (db *UserDB) CreateSession(email string) (string, error) {
	token := generateToken()
	query := "INSERT INTO user_sessions (token, user_email, created_at) VALUES ($1, $2, NOW())"
	_, err := db.db.Exec(query, token, email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (db *UserDB) GetUserByToken(token string) (string, error) {
	var email string
	query := "SELECT user_email FROM user_sessions WHERE token = $1 AND created_at > NOW() - INTERVAL '30 days'"
	err := db.db.QueryRow(query, token).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func (db *UserDB) DeleteSession(token string) error {
	query := "DELETE FROM user_sessions WHERE token = $1"
	_, err := db.db.Exec(query, token)
	return err
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (db *UserDB) AddToCart(userEmail string, productID, quantity int) error {
	var exists bool
	err := db.db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", productID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("product not found")
	}

	query := `
		INSERT INTO cart (user_email, product_id, quantity) 
		VALUES ($1, $2, $3) 
		ON CONFLICT (user_email, product_id) 
		DO UPDATE SET quantity = cart.quantity + EXCLUDED.quantity
	`
	_, err = db.db.Exec(query, userEmail, productID, quantity)
	return err
}

func (db *UserDB) UpdateCartItem(userEmail string, productID, quantity int) error {
	if quantity <= 0 {
		return db.RemoveFromCart(userEmail, productID)
	}

	query := "UPDATE cart SET quantity = $1 WHERE user_email = $2 AND product_id = $3"
	result, err := db.db.Exec(query, quantity, userEmail, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

func (db *UserDB) RemoveFromCart(userEmail string, productID int) error {
	query := "DELETE FROM cart WHERE user_email = $1 AND product_id = $2"
	result, err := db.db.Exec(query, userEmail, productID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

func (db *UserDB) GetCart(userEmail string) ([]models.CartItem, error) {
	query := `
		SELECT c.product_id, c.quantity, p.name, p.price, p.imageurl 
		FROM cart c 
		JOIN products p ON c.product_id = p.id 
		WHERE c.user_email = $1
		ORDER BY c.product_id
	`
	rows, err := db.db.Query(query, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ProductID, &item.Quantity, &item.Name, &item.Price, &item.ImageURL)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, item)
	}

	return cartItems, nil
}

func (db *UserDB) ClearCart(userEmail string) error {
	query := "DELETE FROM cart WHERE user_email = $1"
	_, err := db.db.Exec(query, userEmail)
	return err
}
