package handler

import (
	"encoding/json"
	"forum1/internal/service"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	_, err := h.service.Register(r.Context(), username, email, password)
	if err != nil {
		http.Error(w, "Ошибка регистрации", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	u, err := h.service.Login(r.Context(), username, password)
	if err != nil {
		http.Error(w, "Неверные данные", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "user", Value: u.Username, Path: "/", HttpOnly: true})
	if r.Header.Get("Accept") == "application/json" {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
