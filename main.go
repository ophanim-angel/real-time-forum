package main

import (
	"fmt"
	"log"
	"net/http"

	"toolkit/backend/models"
)

func main() {
	dbPath := "./database/forum.db"
	db, err := models.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error initialization DB: %v", err)
	}
	defer db.Close()
	fmt.Println("Database ready and Tables created!")

	fmt.Println("Server starting at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
