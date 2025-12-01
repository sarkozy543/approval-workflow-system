package main

import (
	"log"
	"net/http"

	"github.com/sarkozy543/approval-workflow-system/internal/server"
)

func main() {
	srv := server.NewServer()

	log.Println("✅ Approval API is starting on :8080...")
	if err := http.ListenAndServe(":8080", srv.Router()); err != nil {
		log.Fatal("server failed:", err)
	}
}
