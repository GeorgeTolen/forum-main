package handlers

import (
	"database/sql"
	"fmt"
	"forum1/db"
	"forum1/internal/entity"
	"forum1/internal/models"
	"forum1/utils"
	"io"
	"net/http"
	"strconv"
)

// CreatePostPage godoc
// @Summary Create a new post
// @Description Создает пост (multipart/form-data поддерживается для image)
// @Tags Posts
// @Accept multipart/form-data
// @Produce json
// @Param board_id formData int true "Board ID"
// @Param title formData string true "Title"
// @Param content formData string true "Content"
// @Param image formData file false "Image file"
// @Success 302 {string} string "redirect to board"
// @Failure 400 {object} map[string]string
// @Router /create_post_page/ [post]
func CreatePostPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Redirect(w, r, "/login_page/", http.StatusSeeOther)
		return
	}
	username := cookie.Value

	var userID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userID)
	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/login_page/", http.StatusSeeOther)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при проверке пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		boardIDStr := r.FormValue("board_id")
		title := r.FormValue("title")
		content := r.FormValue("content")

		var imageBytes []byte
		r.ParseMultipartForm(10 << 20) // 10MB
		file, _, err := r.FormFile("image")
		if err == nil && file != nil {
			defer file.Close()
			buf, _ := io.ReadAll(file)
			imageBytes = buf
		}

		if boardIDStr == "" || title == "" || content == "" {
			http.Error(w, "Заполните все поля", http.StatusBadRequest)
			return
		}

		boardID, err := strconv.Atoi(boardIDStr)
		if err != nil {
			http.Error(w, "Неверный ID доски", http.StatusBadRequest)
			return
		}

		post := &entity.Post{
			BoardID:   boardID,
			Title:     title,
			Content:   content,
			AuthorID:  userID,
			ImageData: imageBytes,
		}

		err = models.CreatePost(post)
		if err != nil {
			http.Error(w, "Ошибка при сохранении поста: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var slug string
		for _, b := range Boards {
			if b.ID == boardID {
				slug = b.Slug
				break
			}
		}
		http.Redirect(w, r, fmt.Sprintf("/board/?slug=%s", slug), http.StatusSeeOther)

	}

	utils.RenderTemplate(w, "create_post_page.html", Boards)
}

// PostPage godoc
// @Summary Get a post
// @Description Возвращает пост по id
// @Tags Posts
// @Produce json
// @Param id query int true "Post ID"
// @Success 200 {object} entity.Post
// @Failure 404 {object} map[string]string
// @Router /post_page/ [get]
func PostPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID поста", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	post, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if cookie, cerr := r.Cookie("user"); cerr == nil {
		var uid int
		if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&uid); err == nil {
			_, _ = db.DB.Exec(`INSERT INTO post_views (post_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, id, uid)
		}
	}
	var views int
	_ = db.DB.QueryRow(`SELECT COUNT(*) FROM post_views WHERE post_id=$1`, id).Scan(&views)

	comments, _ := models.GetCommentsByPost(id)
	post.Comments = comments

	_ = db.DB.QueryRow(`SELECT COALESCE(SUM(CASE WHEN value=1 THEN 1 ELSE 0 END),0) AS likes,
		COALESCE(SUM(CASE WHEN value=-1 THEN 1 ELSE 0 END),0) AS dislikes
		FROM post_votes WHERE post_id=$1`, id).Scan(&post.Likes, &post.Dislikes)

	utils.RenderTemplate(w, "post_page.html", post)
}

// GetPostsByBoard (если экспортирован как handler или используется в модели)
// @Summary Get posts by board
// @Tags Posts
// @Produce json
// @Param board_id query int true "Board ID"
// @Success 200 {array} entity.Post
// @Router /posts/by_board [get]
func GetPostsByBoard(boardID int) ([]entity.Post, error) {
	rows, err := db.DB.Query(`
		SELECT id, board_id, title, content, author_id,
		       created_at, updated_at, image_url, link_url
		FROM posts
		WHERE board_id = $1
		ORDER BY created_at DESC
	`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var p entity.Post
		if err := rows.Scan(
			&p.ID, &p.BoardID, &p.Title, &p.Content,
			&p.AuthorID, &p.CreatedAt, &p.UpdatedAt,
			&p.ImageURL, &p.LinkURL,
		); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func PostImage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID поста", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}
	post, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении поста: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(post.ImageData) == 0 {
		http.NotFound(w, r)
		return
	}
	contentType := http.DetectContentType(post.ImageData)
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(post.ImageData)
}

func EditPostPage(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Не указан ID поста", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	post, err := models.GetPostByID(id)
	if err == sql.ErrNoRows {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Ошибка при получении поста: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		if title != "" {
			post.Title = title
		}
		if content != "" {
			post.Content = content
		}

		err = models.UpdatePost(post)
		if err != nil {
			http.Error(w, "Ошибка при обновлении поста: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post_page?id=%d", post.ID), http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, "edit_post_page.html", post)
}

func VotePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("post_id")
	valStr := r.FormValue("value")
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	id, _ := strconv.Atoi(idStr)
	value, _ := strconv.Atoi(valStr)
	_, err = db.DB.Exec(`INSERT INTO post_votes (post_id, user_id, value) VALUES ($1,$2,$3) ON CONFLICT (post_id,user_id) DO UPDATE SET value=EXCLUDED.value`, id, userID, value)
	if err != nil {
		http.Error(w, "Ошибка голосования: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post_page?id="+idStr, http.StatusSeeOther)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	postIDStr := r.FormValue("post_id")
	content := r.FormValue("content")
	postID, _ := strconv.Atoi(postIDStr)
	c := &entity.Comment{PostID: postID, AuthorID: userID, Content: content}
	if err := models.CreateComment(c); err != nil {
		http.Error(w, "Ошибка добавления комментария: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post_page?id="+postIDStr, http.StatusSeeOther)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	postIDStr := r.FormValue("post_id")
	commentIDStr := r.FormValue("comment_id")
	postID, _ := strconv.Atoi(postIDStr)
	commentID, _ := strconv.Atoi(commentIDStr)
	post, err := models.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}
	if post.AuthorID == userID {
		_ = models.ForceDeleteComment(commentID)
	} else {
		_ = models.DeleteComment(commentID, userID)
	}
	http.Redirect(w, r, "/post_page?id="+postIDStr, http.StatusSeeOther)
}

func VoteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		http.Error(w, "Требуется вход", http.StatusUnauthorized)
		return
	}
	var userID int
	if err := db.DB.QueryRow("SELECT id FROM users WHERE username=$1", cookie.Value).Scan(&userID); err != nil {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}
	postIDStr := r.FormValue("post_id")
	commentIDStr := r.FormValue("comment_id")
	valueStr := r.FormValue("value")
	commentID, _ := strconv.Atoi(commentIDStr)
	value, _ := strconv.Atoi(valueStr)
	if err := models.SetCommentVote(commentID, userID, value); err != nil {
		http.Error(w, "Ошибка голосования: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/post_page?id="+postIDStr, http.StatusSeeOther)
}
