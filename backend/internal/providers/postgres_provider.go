package providers

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
)

var postgresProviderLock = &sync.Mutex{}
var postgresDb *sql.DB

func MakePostgres() (*sql.DB, error) {
	var err error
	if postgresDb == nil {
		user := os.Getenv("PG_USER")
		dbname := os.Getenv("PG_DB")
		password := os.Getenv("PG_PASSWOD")
		host := os.Getenv("PG_HOST")
		port := os.Getenv("PG_PORT")
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s", user, password, dbname, host, port)
		postgresDb, err = sql.Open("postgres", connStr)
	}
	return postgresDb, err
}
