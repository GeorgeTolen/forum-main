package handlers

import (
	"encoding/json"
	"forum1/internal/entity"
	"forum1/internal/models"
	"forum1/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetAllPostsAPI(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Ошибка при получении постов", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostByIDAPI(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не передан id", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный id", http.StatusBadRequest)
		return
	}

	post, err := models.GetPostByID(int64(id))
	if err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// API для досок
type BoardAPIHandler struct {
	boardService service.BoardService
}

func NewBoardAPIHandler(bs service.BoardService) *BoardAPIHandler {
	return &BoardAPIHandler{boardService: bs}
}

// GET /api/boards - получить все доски
func (h *BoardAPIHandler) GetAllBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.boardService.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// GET /api/clubs/{id}/boards - получить доски клуба
func (h *BoardAPIHandler) GetClubBoards(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clubIDStr := vars["id"]
	clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid club id", http.StatusBadRequest)
		return
	}

	boards, err := h.boardService.GetByClubID(r.Context(), clubID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boards)
}

// POST /api/boards - создать доску
func (h *BoardAPIHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации (как в PostHandler)
	c, errCookie := r.Cookie("user")
	if errCookie != nil || c.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Поддержка multipart/form-data
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "invalid multipart form", http.StatusBadRequest)
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
	}

	board := &entity.Board{
		Title:       r.FormValue("title"),
		Slug:        r.FormValue("slug"),
		Description: r.FormValue("description"),
	}

	// Обработка club_id
	clubIDStr := r.FormValue("club_id")
	if clubIDStr != "" {
		clubID, err := strconv.ParseInt(clubIDStr, 10, 64)
		if err == nil {
			board.ClubID = &clubID
		}
	}

	id, err := h.boardService.Create(r.Context(), board)
	if err != nil {
		http.Error(w, "Ошибка создания доски: "+err.Error(), http.StatusBadRequest)
		return
	}

	board.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(board)
}

// GET /api/auth/check - проверка авторизации
func (h *BoardAPIHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	// Проверка авторизации (как в PostHandler)
	c, errCookie := r.Cookie("user")
	if errCookie != nil || c.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Авторизован - возвращаем информацию о пользователе
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "authorized",
		"user":   c.Value,
	})
}

// GET /api/search - поиск по всем типам контента
func (h *BoardAPIHandler) SearchAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter required", http.StatusBadRequest)
		return
	}

	results, err := models.SearchAll(query)
	if err != nil {
		http.Error(w, "search error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
