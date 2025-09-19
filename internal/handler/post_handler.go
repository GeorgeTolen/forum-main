package handler

import (
	"encoding/json"
	"forum1/internal/entity"
	"forum1/internal/repository"
	"forum1/internal/service"
	"io"
	"net/http"
	"strconv"
)

type PostHandler struct {
	svc   service.PostService
	users repository.UserRepository
}

func NewPostHandler(svc service.PostService, users repository.UserRepository) *PostHandler {
	return &PostHandler{svc: svc, users: users}
}

func (h *PostHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("home"))
}

func (h *PostHandler) GetPostPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("post page"))
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if ct := r.Header.Get("Content-Type"); ct != "" && (ct == "application/json" || ct[:16] == "application/json") {
		var p entity.Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		id, err := h.svc.CreatePost(r.Context(), &p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	boardID, _ := strconv.ParseInt(r.FormValue("board_id"), 10, 64)
	title := r.FormValue("title")
	content := r.FormValue("content")
	var imageData []byte
	file, _, err := r.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()
		imageData, _ = io.ReadAll(file)
	}
	// Resolve author from cookie
	c, errCookie := r.Cookie("user")
	if errCookie != nil || c.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	u, err := h.users.GetUserByName(r.Context(), c.Value)
	if err != nil || u == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	p := &entity.Post{BoardID: int(boardID), Title: title, Content: content, AuthorID: int(u.ID), ImageData: imageData}
	id, err := h.svc.CreatePost(r.Context(), p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) GetPostsJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	posts, err := h.svc.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch posts", http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(posts)
}
