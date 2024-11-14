package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB() *sql.DB {
	// Setup Database connection
	dbUser := "root" // Ganti dengan user database MySQL Anda
	dbPass := ""     // Ganti dengan password database MySQL Anda
	dbName := "myapp_db"

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", dbUser, dbPass, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
		os.Exit(1)
	}
	return db
}
