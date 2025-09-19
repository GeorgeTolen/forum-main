package handlers

import (
	"database/sql"
	"net/http"

	"forum1/db"
	"forum1/internal/entity"
	"forum1/utils" // для Hash/CheckPassword если есть
)

// LoginPage godoc
// @Summary User login
// @Tags Auth
// @Accept application/x-www-form-urlencoded
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 302 {string} string "redirect to profile"
// @Failure 401 {object} map[string]string
// @Router /login_page/ [post]
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user entity.User
		err := db.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username=$1", username).
			Scan(&user.ID, &user.Username, &user.Email, &user.Password)
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if !utils.CheckPasswordHash(password, user.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "user", Value: user.Username, Path: "/"})
		http.Redirect(w, r, "/profile_page/", http.StatusSeeOther)
		return
	}
	utils.RenderTemplate(w, "login_page.html", nil)
	// render template
}

// RegisterPage godoc
// @Summary Register new user
// @Tags Auth
// @Accept application/x-www-form-urlencoded
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Success 302 {string} string "redirect to login"
// @Failure 400 {object} map[string]string
// @Router /register_page/ [post]
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Хэшируем пароль
		hashedPassword, _ := utils.HashPassword(password)

		// Сохраняем пользователя в БД
		_, err := db.DB.Exec(
			"INSERT INTO users (username, email, password) VALUES ($1, $2, $3)",
			username, email, hashedPassword,
		)
		if err != nil {
			http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// После регистрации — редирект на логин
		http.Redirect(w, r, "/login_page/", http.StatusSeeOther)
		return
	}

	// Показать форму регистрации
	utils.RenderTemplate(w, "register_page.html", nil)
}
