package models

// Post represents a post in the database
type Post struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Nickname  string `json:"nickname"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Category  string `json:"category"`
	Views     string `json:"views"`
	CreatedAt string `json:"created_at"`
}

// CreatePostInput: Data coming from frontend to create a post
type CreatePostInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

// PostReaction represents a reaction (like, love, etc.)
type PostReaction struct {
	ID        string `json:"id"`
	UderID    string `json:"user_id"`
	PostID    string `json:"post_id"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}
