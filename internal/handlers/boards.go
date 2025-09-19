package handlers

import (
	"forum1/internal/entity"
	"forum1/internal/models"
	"forum1/utils"
	"net/http"
	"strings"
)

var Boards = []entity.Board{
	{ID: 1, Slug: "schedule", Title: "Расписание", Description: "Обсуждаем расписание этого года"},
	{ID: 2, Slug: "games", Title: "Игры", Description: "Все о видеоиграх, консолях и ПК"},
	{ID: 3, Slug: "offtopic", Title: "Оффтопик", Description: "Свободное общение на любые темы"},
	{ID: 4, Slug: "news", Title: "Новости", Description: "Обсуждение последних новостей"},
	{ID: 5, Slug: "reviews", Title: "Рецензии", Description: "Ваши обзоры на фильмы, игры и книги"},
}

// BoardPage godoc
// @Summary Get board by slug
// @Description Возвращает доску и её посты по slug
// @Tags Boards
// @Produce json
// @Param slug query string true "Board slug"
// @Success 200 {object} entity.BoardWithPostsResponse

// @Failure 400 {object} map[string]string
// @Router /board/ [get]
func BoardPage(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		http.Error(w, "Не указана доска", http.StatusBadRequest)
		return
	}

	var board *entity.Board
	for _, b := range Boards {
		if b.Slug == slug {
			board = &b
			break
		}
	}
	if board == nil {
		http.Error(w, "Доска не найдена", http.StatusNotFound)
		return
	}

	// Здесь уже получаем только посты для нужной доски
	posts, err := models.GetPostsByBoard(board.ID)
	if err != nil {
		http.Error(w, "Ошибка при получении постов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Board entity.Board
		Posts []entity.Post
	}{
		Board: *board,
		Posts: posts, // напрямую передаём
	}

	utils.RenderTemplate(w, "board_page.html", data)
}

// BoardsListPage godoc
// @Summary Get list of boards
// @Description Возвращает список всех досок
// @Tags Boards
// @Produce json
// @Success 200 {array} entity.Board
// @Router /boards_list_page/ [get]
func BoardsListPage(w http.ResponseWriter, r *http.Request) {
	// ... existing code ...

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	var filtered []entity.Board
	if query == "" {
		filtered = Boards
	} else {
		for _, b := range Boards {
			if strings.Contains(strings.ToLower(b.Title), strings.ToLower(query)) ||
				strings.Contains(strings.ToLower(b.Description), strings.ToLower(query)) {
				filtered = append(filtered, b)
			}
		}
	}

	data := struct {
		Query  string
		Boards []entity.Board
	}{
		Query:  query,
		Boards: filtered,
	}

	utils.RenderTemplate(w, "boards_list_page.html", data)
}

func BoardsSearchPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Не указан поисковый запрос", http.StatusBadRequest)
		return
	}

	// --- Поиск по доскам ---
	var boardResults []entity.Board
	for _, b := range Boards {
		if strings.Contains(strings.ToLower(b.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(b.Description), strings.ToLower(query)) {
			boardResults = append(boardResults, b)
		}
	}

	// --- Поиск по постам ---
	posts, err := models.GetAllPosts()
	if err != nil {
		http.Error(w, "Ошибка при получении постов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var postResults []entity.Post
	for _, p := range posts {
		if strings.Contains(strings.ToLower(p.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(p.Content), strings.ToLower(query)) {
			postResults = append(postResults, p)
		}
	}

	// --- Данные для шаблона ---
	data := struct {
		Query        string
		BoardsResult []entity.Board
		PostsResult  []entity.Post
	}{
		Query:        query,
		BoardsResult: boardResults,
		PostsResult:  postResults,
	}

	utils.RenderTemplate(w, "boards_search_page.html", data)
}
