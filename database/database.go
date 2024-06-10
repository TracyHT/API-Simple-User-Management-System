package database

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func InitDB(dataSourceName string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create 'users' table if not exists
	_, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	username TEXT NOT NULL UNIQUE,
    	firstname TEXT NOT NULL,
   		lastname TEXT NOT NULL,
    	email TEXT UNIQUE,
    	avatar TEXT,
    	phone TEXT,
    	date_of_birth DATE CHECK(date_of_birth IS NULL OR date_of_birth <= date('now')),
    	address_country TEXT,
    	address_city TEXT,
    	address_street_name TEXT,
    	address_street_address TEXT
	);
    `)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return DB, nil
}
