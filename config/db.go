package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func OpenDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/ecommerce_db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
