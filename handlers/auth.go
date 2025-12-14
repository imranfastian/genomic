package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/imranfastian/genomic/middleware"
)

// LoginRequest represents the expected JSON body for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginHandler handles user login and JWT generation
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Replace this with real DB check or LDAP in production
	if req.Username == "admin" && req.Password == "password" {
		token, err := middleware.GenerateJWT(req.Username)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
