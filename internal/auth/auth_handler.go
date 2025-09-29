package auth

import (
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	access, refresh, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	resp := map[string]string{"access_token": access, "refresh_token": refresh}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
