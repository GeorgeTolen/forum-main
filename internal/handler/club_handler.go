package handler

import (
	"encoding/json"
	"forum1/internal/entity"
	"forum1/internal/service"
	"forum1/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ClubHandler struct {
	service service.ClubService
}

func NewClubHandler(s service.ClubService) *ClubHandler {
	return &ClubHandler{service: s}
}

// POST /clubs
func (h *ClubHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c entity.Club
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(r.Context(), &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// GET /clubs/{id}
func (h *ClubHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	c, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// GET /clubs
func (h *ClubHandler) List(w http.ResponseWriter, r *http.Request) {
	clubs, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clubs)
}

type ClubPageHandler struct {
	service service.ClubService
}

func NewClubPageHandler(s service.ClubService) *ClubPageHandler {
	return &ClubPageHandler{service: s}
}

// GET /boards/club
func (h *ClubPageHandler) ListPage(w http.ResponseWriter, r *http.Request) {
	clubs, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Render with shared layout
	// Root context is a slice to iterate in template
	utils.RenderTemplate(w, "clubs.html", clubs)
}

// GET /boards/club/{id}
func (h *ClubPageHandler) DetailPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	club, err := h.service.GetByID(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Render with shared layout
	utils.RenderTemplate(w, "club_detail.html", club)
}

// GET /clubs/new
func (h *ClubPageHandler) NewPage(w http.ResponseWriter, r *http.Request) {
	utils.RenderTemplate(w, "club_form.html", nil)
}

// POST /clubs (HTML-форма)
func (h *ClubPageHandler) CreatePage(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	club := entity.Club{
		Name:        r.FormValue("name"),
		Topic:       r.FormValue("topic"),
		Description: r.FormValue("description"),
	}

	id, err := h.service.Create(r.Context(), &club)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/clubs/"+strconv.FormatInt(id, 10), http.StatusSeeOther)
}
