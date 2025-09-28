package database

import "database/sql"

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
