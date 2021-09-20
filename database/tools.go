package database

import (
	"database/sql"
	"fmt"
)

const (
	host     string = "db"
	port     int    = 5432
	user     string = "khong"
	password string = "khongpass"
	dbName   string = "khong"
)

func GetDB() (*sql.DB, error) {
	SourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	db, err := sql.Open("postgres", SourceName)
	return db, err
}
