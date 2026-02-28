package main

import (
	"fmt"
	"log"
	"net/http"
	"toolkit/backend/api"
	"toolkit/backend/models"
	"toolkit/backend/websocket"
)

func main() {
	
	db, err := models.InitDB("./database/forum.db")
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	
	hub := websocket.NewHub(db)
	go hub.Run()

	
	http.HandleFunc("/api/register", api.RegisterHandler(db))
	http.HandleFunc("/api/login", api.LoginHandler(db))
	http.HandleFunc("/api/logout", api.AuthMiddleware(db, api.LogoutHandler(db)))
	http.HandleFunc("/api/me", api.AuthMiddleware(db, api.MeHandler(db)))

	
	
	http.HandleFunc("/api/posts", api.GetPostsHandler(db))
	http.HandleFunc("/api/posts/create", api.AuthMiddleware(db, api.CreatePostHandler(db)))

	
	http.HandleFunc("/api/comments", api.GetCommentsHandler(db))
	http.HandleFunc("/api/comments/create", api.AuthMiddleware(db, api.CreateCommentHandler(db)))
	http.HandleFunc("/api/react", api.AuthMiddleware(db, api.ReactionHandler(db)))

	
	http.HandleFunc("/api/users", api.AuthMiddleware(db, api.GetUsersHandler(db)))

	
	http.HandleFunc("/ws", api.AuthMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	}))

	
	
	fs := http.FileServer(http.Dir("./frontend/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			
			http.ServeFile(w, r, "./frontend/index.html")
			return
		}
		http.ServeFile(w, r, "./frontend/index.html")
	})

	
	port := ":8080"
	fmt.Printf("Server started on http:localhost:8080")
	fmt.Println("Press Ctrl+C to stop")

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
