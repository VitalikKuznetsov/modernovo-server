package database

import (
	"Modernovo/GoBackend/WorkingWithTheStore/models"
	"database/sql"
	"errors"
	"strings"
)

type StoreDB struct {
	db *sql.DB
}

func ConnectToStoreDB(connectSring string) (*StoreDB, error) {
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

func (db *StoreDB) GetProduct(id int) (models.Product, error) {
	var p models.Product
	query := "SELECT id, name, description, price, imageurl, image_urls FROM products WHERE id = $1"
	row := db.db.QueryRow(query, id)

	var imageURLsStr sql.NullString
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.ImageURL,
		&imageURLsStr,
	)

	if err == sql.ErrNoRows {
		return p, errors.New("product not found")
	}
	if err != nil {
		return p, err
	}

	if imageURLsStr.Valid && imageURLsStr.String != "" {
		cleaned := strings.Trim(imageURLsStr.String, "{}")
		if cleaned != "" {
			urls := strings.Split(cleaned, ",")
			p.ImageURLs = urls
		}
	}

	if len(p.ImageURLs) == 0 && p.ImageURL != "" {
		p.ImageURLs = []string{p.ImageURL}
	}

	return p, nil
}

func (db *StoreDB) GetAllProducts(limit, offset int) ([]models.Product, error) {
	query := "SELECT id, name, description, price, imageurl, image_urls FROM products ORDER BY id LIMIT $1 OFFSET $2"
	rows, err := db.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var imageURLsStr sql.NullString

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.ImageURL,
			&imageURLsStr,
		)
		if err != nil {
			return nil, err
		}

		if imageURLsStr.Valid && imageURLsStr.String != "" {
			cleaned := strings.Trim(imageURLsStr.String, "{}")
			if cleaned != "" {
				urls := strings.Split(cleaned, ",")
				p.ImageURLs = urls
			}
		}

		if len(p.ImageURLs) == 0 && p.ImageURL != "" {
			p.ImageURLs = []string{p.ImageURL}
		}

		products = append(products, p)
	}

	return products, nil
}

func (db *StoreDB) GetProductsCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM products"
	err := db.db.QueryRow(query).Scan(&count)
	return count, err
}
