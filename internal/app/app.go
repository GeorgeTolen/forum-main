package app

import (
	"fmt"
	"forum1/db"
	handler "forum1/internal/handler"
	"forum1/internal/repository"
	"forum1/internal/router"
	"forum1/internal/service"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Run() {
	// подключение к БД
	err := db.InitDB()
	if err != nil {
		fmt.Println("Ошибка подключения к базе:", err)
		return
	}
	defer db.CloseDB()

	database := db.GetDB() // получаем *sql.DB

	// слой repository
	postRepo := repository.NewPostRepository(database)
	boardRepo := repository.NewBoardRepository(database)
	commentRepo := repository.NewCommentRepository(database)
	clubRepo := repository.NewClubRepository(database)

	// слой service
	postService := service.NewPostService(postRepo)
	boardService := service.NewBoardService(boardRepo)
	commentService := service.NewCommentService(commentRepo)
	clubService := service.NewClubService(clubRepo)

	// слой handler
	userRepo := repository.NewUserRepository(database)
	postHandler := handler.NewPostHandler(postService, userRepo)
	commentHandler := handler.NewCommentHandler(commentService, userRepo).WithPosts(postService)
	pageHandler := handler.NewPageHandler(postService, boardService).WithComments(commentService)
	userHandler := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(database)))
	clubPageHandler := handler.NewClubPageHandler(clubService)
	clubAPIHandler := handler.NewClubHandler(clubService)

	// слой router
	r := router.NewRouter(postHandler)
	// HTML routes for templates
	r.HandleFunc("/", pageHandler.HomePageHTML).Methods(http.MethodGet)
	r.HandleFunc("/boards", pageHandler.BoardsListPage).Methods(http.MethodGet)
	r.HandleFunc("/board/{slug}", pageHandler.BoardPage).Methods(http.MethodGet)
	// Clubs pages
	r.HandleFunc("/clubs", clubPageHandler.ListPage).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}", pageHandler.PostPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/clubs", clubPageHandler.ListPage).Methods("GET")
	r.HandleFunc("/clubs/{id:[0-9]+}", clubPageHandler.DetailPage).Methods("GET")
	r.HandleFunc("/clubs/new", clubPageHandler.NewPage).Methods("GET")
	r.HandleFunc("/clubs", clubPageHandler.CreatePage).Methods("POST")

	// post image
	r.HandleFunc("/post/{id}/image", pageHandler.PostImage).Methods(http.MethodGet)
	// like/dislike GET endpoints
	r.HandleFunc("/post/{id}/like", pageHandler.LikePost).Methods(http.MethodGet)
	r.HandleFunc("/post/{id}/dislike", pageHandler.DislikePost).Methods(http.MethodGet)
	r.HandleFunc("/comment/{id}/like", pageHandler.LikeComment).Methods(http.MethodGet)
	r.HandleFunc("/comment/{id}/dislike", pageHandler.DislikeComment).Methods(http.MethodGet)
	r.HandleFunc("/profile/{id}", pageHandler.ProfilePageHTML).Methods(http.MethodGet)
	r.HandleFunc("/login", pageHandler.LoginPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/register", pageHandler.RegisterPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/create-post", pageHandler.CreatePostPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/boards/search", pageHandler.BoardsSearchPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/search", pageHandler.SearchPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/settings", pageHandler.SettingsPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/messages", pageHandler.MessagesPageHTML).Methods(http.MethodGet)
	r.HandleFunc("/notifications", pageHandler.NotificationsPageHTML).Methods(http.MethodGet)

	// CORS (dev permissive)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	// Logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, req)
			fmt.Printf("%s %s %s\n", req.Method, req.URL.Path, time.Since(start))
		})
	})

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API auth endpoints
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", userHandler.RegisterPage).Methods(http.MethodPost)
	api.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/comment", commentHandler.CreateComment).Methods(http.MethodPost)
	api.HandleFunc("/delete_comment", commentHandler.DeleteComment).Methods(http.MethodPost)
	// Clubs API
	api.HandleFunc("/clubs", clubAPIHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/clubs", clubAPIHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/clubs/{id}", clubAPIHandler.GetByID).Methods(http.MethodGet)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
