package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"toolKit/backend/models"
	"toolKit/backend/utils"

	"github.com/mattn/go-sqlite3"
)

// AuthHandler holds the DB connection
type AuthHandler struct {
	DB *sql.DB
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Decode JSON input
	var input models.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Validate Input (using utils)
	if !utils.ValidateEmail(input.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	if !utils.ValidateNickname(input.Nickname) {
		http.Error(w, "Invalid nickname (3-20 chars, letters/numbers/underscore)", http.StatusBadRequest)
		return
	}
	if err := utils.ValidatePassword(input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !utils.ValidateName(input.FirstName) || !utils.ValidateName(input.LastName) {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}
	if !utils.ValidateAge(input.Age) {
		http.Error(w, "Invalid age (13-120)", http.StatusBadRequest)
		return
	}
	if !utils.ValidateGender(input.Gender) {
		http.Error(w, "Invalid gender", http.StatusBadRequest)
		return
	}

	// 3. Hash Password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// 4. Generate UUID
	userID := utils.GenerateUUID()

	// 5. Insert into Database (SQLi Protected with ?)
	query := `
		INSERT INTO users (id, nickname, email, password_hash, first_name, last_name, age, gender)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = h.DB.Exec(query, userID, input.Nickname, input.Email, hashedPassword, input.FirstName, input.LastName, input.Age, input.Gender)
	if err != nil {
		// Check if it's a SQLite error
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			// Code 19 = SQLITE_CONSTRAINT (UNIQUE constraint failed)
			if sqliteErr.Code == sqlite3.ErrConstraint {
				http.Error(w, "Nickname or Email already exists", http.StatusConflict)
				return
			}
		}
		log.Println("Database error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 6. Success Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User registered successfully",
		"user_id":  userID,
		"nickname": input.Nickname,
	})
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Decode JSON input
	var input models.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Validate Credentials
	if err := utils.ValidateCredentials(input.Identifier, input.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Find User by Email OR Nickname
	query := `SELECT id, nickname, password_hash FROM users WHERE (email = ? OR nickname = ?) AND is_active = 1`
	var userID, nickname, passwordHash string
	err := h.DB.QueryRow(query, input.Identifier, input.Identifier).Scan(&userID, &nickname, &passwordHash)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// 4. Check Password
	if !utils.CheckPassword(input.Password, passwordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 5. Generate JWT Token
	token, err := utils.GenerateToken(userID, nickname)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// 6. Success Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Login successful",
		"token":    token,
		"user_id":  userID,
		"nickname": nickname,
	})
}