package handlers

import (
	"database/sql"
	"forum1/db"
	"forum1/internal/entity"
	"forum1/utils"
	"net/http"
)

// Профиль пользователя
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Redirect(w, r, "/login_page/", http.StatusSeeOther)
		return
	}

	username := cookie.Value

	// Загружаем пользователя из БД
	var user entity.User
	err = db.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username=$1", username).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/login_page/", http.StatusSeeOther)
		return
	} else if err != nil {
		http.Error(w, "Ошибка загрузки профиля: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		newEmail := r.FormValue("email")
		newPassword := r.FormValue("password")

		if newEmail != "" {
			user.Email = newEmail
		}

		if newPassword != "" {
			hashedPassword, _ := utils.HashPassword(newPassword)
			user.Password = hashedPassword
		}

		// Обновляем в БД
		_, err := db.DB.Exec("UPDATE users SET email=$1, password=$2 WHERE id=$3",
			user.Email, user.Password, user.ID)
		if err != nil {
			http.Error(w, "Ошибка при обновлении профиля: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/profile_page/", http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, "profile_page.html", user)
}

// Страница редактирования профиля
func EditProfilePage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "edit_profile_page.html", nil)
}
