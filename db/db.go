package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal((err))
	}

	return db, err
}

func InitDB(cfg mysql.Config) (*sql.DB, error) {
	db, err := NewMySQLStorage(cfg)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}
