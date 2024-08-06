package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func DBConn() (*sql.DB, error) {
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    host := os.Getenv("DB_HOST")
    port := 5432
    connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", user, password, dbname, host, port)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}

func InitDB() {
    db, err := DBConn()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    createTableQuery := `
    CREATE TABLE IF NOT EXISTS materials (
        uuid UUID PRIMARY KEY,
        material_type TEXT NOT NULL CHECK (material_type IN ('статья', 'видеоролик', 'презентация')),
        status TEXT NOT NULL CHECK (status IN ('архивный', 'активный')),
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `

    _, err = db.Exec(createTableQuery)
    if err != nil {
        log.Fatal(err)
    }
}
