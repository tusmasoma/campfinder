package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // This blank import is used for its init function
)

var (
	dbUser     = os.Getenv("MYSQL_ROOT_USER")
	dbPassword = os.Getenv("MYSQL_ROOT_PASSWORD")
	dbHost     = os.Getenv("MYSQL_HOST")
	dbPort     = os.Getenv("MYSQL_PORT")
	dbName     = os.Getenv("MYSQL_DB_NAME")
)

func NewDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Database connection failed: %s\n", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
