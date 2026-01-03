package main

import (
	"log"
	"net/http"

	"github.com/sarkozy543/approval-workflow-system/internal/db"
	"github.com/sarkozy543/approval-workflow-system/internal/server"
)

func main() {
	// 1) DB bağlantısını başlat
	database, err := db.New()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer database.Close()

	// 2) Server'ı oluştur (DB ile birlikte)
	srv := server.NewServer(database)

	// 3) HTTP server'ı başlat
	log.Println("✅ Approval API is starting on :8080...")
	if err := http.ListenAndServe(":8080", srv.Router()); err != nil {
		log.Fatal("server failed:", err)
	}
}
