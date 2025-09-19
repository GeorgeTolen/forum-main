package handler

import (
	"database/sql"
	"fmt"
	"forum1/db"
	"forum1/internal/entity"
	"forum1/internal/models"
	"forum1/internal/service"
	"forum1/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type PageHandler struct {
	posts    service.PostService
	boards   service.BoardService
	comments service.CommentService
}

// WithComments allows injecting CommentService fluently after construction
func (h *PageHandler) WithComments(c service.CommentService) *PageHandler {
	h.comments = c
	return h
}

func NewPageHandler(p service.PostService, b service.BoardService) *PageHandler {
	// Backwards-compatible constructor; comments can be injected later if needed
	return &PageHandler{posts: p, boards: b}
}

func (h *PageHandler) HomePageHTML(w http.ResponseWriter, r *http.Request) {
	// Load boards for sidebar/home
	boards, _ := h.boards.List(r.Context())
	data := map[string]interface{}{
		"Boards": boards,
	}
	utils.RenderTemplate(w, "home_page.html", data)
}

func (h *PageHandler) BoardsListPage(w http.ResponseWriter, r *http.Request) {
	boards, _ := h.boards.List(r.Context())
	data := map[string]interface{}{"Boards": boards}
	utils.RenderTemplate(w, "boards_list_page.html", data)
}

func (h *PageHandler) BoardPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]
	b, err := h.boards.GetBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	posts, _ := h.posts.GetPostsByBoard(r.Context(), int64(b.ID))
	data := map[string]interface{}{"Board": b, "Posts": posts}
	utils.RenderTemplate(w, "board_page.html", data)
}

func (h *PageHandler) PostPageHTML(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, _ := strconv.ParseInt(idStr, 10, 64)
	post, _ := h.posts.GetPostByID(r.Context(), id)

	// Load comments with like/dislike counters
	var comments []entity.Comment
	if h.comments != nil {
		cs, _ := h.comments.GetCommentsByPost(r.Context(), id)
		comments = cs
	} else {
		// fallback to models
		comments, _ = models.GetCommentsByPost(int(id))
	}
	if post != nil {
		post.Comments = comments
		// Load like/dislike counters for the post via service
		if likes, dislikes, err := h.posts.GetPostVotes(r.Context(), id); err == nil {
			post.Likes, post.Dislikes = likes, dislikes
		}
	}

	// Pass the post as root context as expected by the template
	utils.RenderTemplate(w, "post_page.html", post)
}

func (h *PageHandler) ProfilePageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "profile_page.html", map[string]interface{}{})
}

func (h *PageHandler) LoginPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "login_page.html", map[string]interface{}{})
}

func (h *PageHandler) RegisterPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "register_page.html", map[string]interface{}{})
}

func (h *PageHandler) CreatePostPageHTML(w http.ResponseWriter, r *http.Request) {
	boards, _ := h.boards.List(r.Context())
	// Template expects to range over root (.)
	utils.RenderTemplate(w, "create_post_page.html", boards)
}

func (h *PageHandler) BoardsSearchPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "boards_search_page.html", map[string]interface{}{})
}

func (h *PageHandler) SearchPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "search_page.html", map[string]interface{}{})
}

func (h *PageHandler) SettingsPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "settings_page.html", map[string]interface{}{})
}

func (h *PageHandler) MessagesPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "messages_page.html", map[string]interface{}{})
}

func (h *PageHandler) NotificationsPageHTML(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "notifications_page.html", map[string]interface{}{})
}

// Serve post image as /post/{id}/image
func (h *PageHandler) PostImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "bad id", http.StatusBadRequest)
		return
	}
	p, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(p.ImageData) == 0 {
		http.NotFound(w, r)
		return
	}
	ct := http.DetectContentType(p.ImageData)
	w.Header().Set("Content-Type", ct)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(p.ImageData)
}

// Like/Dislike post via GET links
func (h *PageHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	h.votePost(w, r, 1)
}

func (h *PageHandler) DislikePost(w http.ResponseWriter, r *http.Request) {
	h.votePost(w, r, -1)
}

// Like/Dislike comment via GET links (?post_id= for redirect)
func (h *PageHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	h.voteComment(w, r, 1)
}

func (h *PageHandler) DislikeComment(w http.ResponseWriter, r *http.Request) {
	h.voteComment(w, r, -1)
}

// Helpers
func (h *PageHandler) votePost(w http.ResponseWriter, r *http.Request, value int) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	postID, _ := strconv.Atoi(idStr)
	if err := h.posts.SetPostVote(r.Context(), int64(postID), int64(userID), value); err != nil {
		http.Error(w, "vote error", http.StatusInternalServerError)
		return
	}
	// If client expects JSON (AJAX), return new counters
	if acceptsJSON(r) {
		likes, dislikes, _ := h.posts.GetPostVotes(r.Context(), int64(postID))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"likes":%d,"dislikes":%d}`, likes, dislikes)))
		return
	}
	http.Redirect(w, r, "/post/"+idStr, http.StatusSeeOther)
}

func (h *PageHandler) voteComment(w http.ResponseWriter, r *http.Request, value int) {
	vars := mux.Vars(r)
	commentIDStr := vars["id"]
	postID := r.URL.Query().Get("post_id")
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	cid, _ := strconv.Atoi(commentIDStr)
	if h.comments != nil {
		if err := h.comments.SetCommentVote(r.Context(), int64(cid), int64(userID), value); err != nil {
			http.Error(w, "vote error", http.StatusInternalServerError)
			return
		}
	} else if _, err := db.DB.Exec(`INSERT INTO comment_votes (comment_id, user_id, value) VALUES ($1,$2,$3)
        ON CONFLICT (comment_id,user_id) DO UPDATE SET value=EXCLUDED.value`, cid, userID, value); err != nil {
		http.Error(w, "vote error", http.StatusInternalServerError)
		return
	}
	if acceptsJSON(r) {
		var likes, dislikes int
		if h.comments != nil {
			likes, dislikes, _ = h.comments.GetCommentVotes(r.Context(), int64(cid))
		} else {
			_ = db.DB.QueryRow(`SELECT COALESCE(SUM(CASE WHEN value=1 THEN 1 ELSE 0 END),0), COALESCE(SUM(CASE WHEN value=-1 THEN 1 ELSE 0 END),0) FROM comment_votes WHERE comment_id=$1`, cid).Scan(&likes, &dislikes)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"likes":%d,"dislikes":%d}`, likes, dislikes)))
		return
	}
	if postID == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

// acceptsJSON returns true if request Accept header prefers JSON
func acceptsJSON(r *http.Request) bool {
	a := r.Header.Get("Accept")
	if a == "" {
		return false
	}
	return strings.Contains(a, "application/json")
}
