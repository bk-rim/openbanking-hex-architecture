package repository

import (
	"database/sql"

	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitSqlite() {

	dbPath := os.Getenv("DB_PATH")
	dbInit, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	db = dbInit
	migrate()
}

func migrate() {
	sqlStmt := `
	create table if not exists payments (id INTEGER PRIMARY KEY AUTOINCREMENT, debtor_iban VARCHAR(30) NOT NULL, debtor_name VARCHAR(30) NOT NULL, creditor_iban VARCHAR(30) NOT NULL,
	creditor_name VARCHAR(30) NOT NULL, amount FLOAT(64) NOT NULL, idempotency VARCHAR(65) UNIQUE NOT NULL, status VARCHAR(30) NOT NULL);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}
}
