package database

import (
	"Modernovo/GoBackend/WorkingWithTheUsers/models"
	"database/sql"
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
