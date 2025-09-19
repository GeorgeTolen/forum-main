package handler

import (
	"encoding/json"
	"forum1/internal/entity"
	"forum1/internal/repository"
	"forum1/internal/service"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	svc   service.CommentService
	users repository.UserRepository
	posts service.PostService
}

func NewCommentHandler(svc service.CommentService, users repository.UserRepository) *CommentHandler {
	return &CommentHandler{svc: svc, users: users}
}

// CreateComment accepts either JSON or form (multipart/urlencoded)
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
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

	// JSON
	if ct := r.Header.Get("Content-Type"); ct != "" && (ct == "application/json" || ct[:16] == "application/json") {
		var in struct {
			PostID  int64  `json:"post_id"`
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		cmt := &entity.Comment{PostID: in.PostID, AuthorID: u.ID, Content: in.Content}
		id, err := h.svc.CreateComment(r.Context(), cmt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
		return
	}

	// Form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	postID, _ := strconv.ParseInt(r.FormValue("post_id"), 10, 64)
	content := r.FormValue("content")
	cmt := &entity.Comment{PostID: postID, AuthorID: u.ID, Content: content}
	id, err := h.svc.CreateComment(r.Context(), cmt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// If client expects JSON (AJAX), return created info
	if r.Header.Get("Accept") != "" && (r.Header.Get("Accept") == "application/json" || r.Header.Get("Accept")[:16] == "application/json") {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
		return
	}
	// Redirect back to the post page
	http.Redirect(w, r, "/post/"+strconv.FormatInt(postID, 10), http.StatusSeeOther)
}

// DeleteComment allows delete by comment author or by post author
func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	// Auth
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

	postID, _ := strconv.ParseInt(r.FormValue("post_id"), 10, 64)
	commentID, _ := strconv.ParseInt(r.FormValue("comment_id"), 10, 64)
	if commentID == 0 || postID == 0 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Load comment to verify author
	cmt, err := h.svc.GetCommentByID(r.Context(), commentID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Load post to verify post author if needed
	var isPostAuthor bool
	if h.posts != nil {
		p, err := h.posts.GetPostByID(r.Context(), postID)
		if err == nil && p != nil && int64(p.AuthorID) == u.ID {
			isPostAuthor = true
		}
	}

	if cmt.AuthorID == u.ID || isPostAuthor {
		_ = h.svc.DeleteComment(r.Context(), commentID, u.ID)
	} else {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	http.Redirect(w, r, "/post/"+strconv.FormatInt(postID, 10), http.StatusSeeOther)
}

// WithPosts enables checking post authorship
func (h *CommentHandler) WithPosts(p service.PostService) *CommentHandler {
	h.posts = p
	return h
}
