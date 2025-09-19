package router

import (
	h "forum1/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(post *h.PostHandler, opts ...Option) *mux.Router {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/", post.HomePage).Methods(http.MethodGet)
	api.HandleFunc("/post/{id}", post.GetPostPage).Methods(http.MethodGet)
	api.HandleFunc("/post", post.CreatePost).Methods(http.MethodPost)
	api.HandleFunc("/post/{id}", post.UpdatePost).Methods(http.MethodPut)
	api.HandleFunc("/post/{id}", post.DeletePost).Methods(http.MethodDelete)
	api.HandleFunc("/posts", post.GetPostsJSON).Methods(http.MethodGet)

	// HTML pages
	// Will be added by app with PageHandler when composed

	return r
}

type Option func(r *mux.Router)
