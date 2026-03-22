//go:build integration

package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func GetTestDB(t *testing.T) *sql.DB {
	t.Helper()

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "transaction_routine_db"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "pg123"
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to DB: %v", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		t.Skipf("Skipping integration test: DB not available: %v", err)
		db.Close()
		return nil
	}

	return db
}
