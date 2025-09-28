package database

import (
	"Modernovo/GoBackend/WorkingWithTheStore/models"
	"database/sql"
	"errors"
)

type StoreDB struct {
	db *sql.DB
}

func ConnectToMyDB(connectSring string) (*StoreDB, error) {
	db, err := sql.Open("postgres", connectSring)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &StoreDB{db: db}, nil
}

func (db *StoreDB) Close() {
	db.db.Close()
}

func (db *StoreDB) AddProduct(id int, name, description string, price float64, url string) error {
	query := "INSERT INTO products (Name, Description, Price, ImageURL, Category, Stock) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := db.db.Exec(query, name, description, price, url)
	if err != nil {
		return errors.New("tralalelo tralala") // rework
	}
	return nil
}

func (db *StoreDB) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE ID = $1"
	result, err := db.db.Exec(query, id)
	if err != nil {
		return errors.New("failed to delete product")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to get rows affected")
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (db *StoreDB) GetInfoOfStoreDB(id int) (models.Product, error) {
	var p models.Product
	query := "SELECT Name, Description, Price, ImageURL FROM products WHERE ID = $1"
	row := db.db.QueryRow(query)
	err := row.Scan(
		&p.Name,
		&p.Description,
		&p.Price,
		&p.ImageURL,
	)
	p.ID = id

	if err == sql.ErrNoRows {
		return p, errors.New("user not found")
	}
	return p, nil
}
