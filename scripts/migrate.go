package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// DSN'i doğrudan sabit kullanıyoruz
	dsn := "postgres://approval_user:approval_pass@localhost:5433/approval_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
	defer db.Close()

	content, err := os.ReadFile("migrations/0001_init.sql")
	if err != nil {
		log.Fatal("failed to read migration file:", err)
	}

	fmt.Println("Running migrations...")

	if _, err := db.Exec(string(content)); err != nil {
		log.Fatal("failed to execute migration:", err)
	}

	fmt.Println("✅ Migration completed successfully.")
}
