// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"toolkit/backend/models"
// )

// func main() {
// 	dbPath := "./database/forum.db"
// 	db, err := models.InitDB(dbPath)
// 	if err != nil {
// 		log.Fatalf(" Erreur f l-database: %v", err)
// 	}
// 	defer db.Close()
// 	fmt.Println("Database initialisée mzyan f:", dbPath)

// 	if _, err := os.Stat("./frontend/index.html"); os.IsNotExist(err) {
// 		fmt.Println("Warning: index.html makhdamch f path ./frontend/index.html")
// 	} else {
// 		fmt.Println("Frontend file (index.html) kayn.")
// 	}

// 	fs := http.FileServer(http.Dir("./frontend"))
// 	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static"))))
// 	http.Handle("/", fs)

// 	port := ":8080"
// 	fmt.Printf("Server khdam f http://localhost%s\n", port)

// 	err = http.ListenAndServe(port, nil)
// 	if err != nil {
// 		log.Fatalf("Server mabhghach y-khdem: %v", err)
// 	}
// }