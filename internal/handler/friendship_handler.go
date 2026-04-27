package handler

import (
	"encoding/json"
	"forum1/internal/entity"
	"forum1/internal/repository"
	"forum1/internal/service"
	"forum1/utils"
	"net/http"
	"strconv"
)

type FriendshipPageHandler struct {
	friendships service.FriendshipService
	users       repository.UserRepository
}

func NewFriendshipPageHandler(f service.FriendshipService, u repository.UserRepository) *FriendshipPageHandler {
	return &FriendshipPageHandler{friendships: f, users: u}
}

func (h *FriendshipPageHandler) FriendsPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		utils.RenderTemplate(w, "auth_required_page.html", map[string]interface{}{"User": nil})
		return
	}
	user, err := h.users.GetUserByName(r.Context(), cookie.Value)
	if err != nil {
		utils.RenderTemplate(w, "auth_required_page.html", map[string]interface{}{"User": nil})
		return
	}

	friends, _ := h.friendships.GetFriends(r.Context(), user.ID)
	incoming, _ := h.friendships.GetPendingIncoming(r.Context(), user.ID)
	outgoing, _ := h.friendships.GetPendingOutgoing(r.Context(), user.ID)

	data := map[string]interface{}{
		"User":     user,
		"Friends":  friends,
		"Incoming": incoming,
		"Outgoing": outgoing,
	}
	utils.RenderTemplate(w, "friends_page.html", data)
}

type FriendshipAPIHandler struct {
	friendships service.FriendshipService
	users       repository.UserRepository
}

func NewFriendshipAPIHandler(f service.FriendshipService, u repository.UserRepository) *FriendshipAPIHandler {
	return &FriendshipAPIHandler{friendships: f, users: u}
}

func (h *FriendshipAPIHandler) getCurrentUser(r *http.Request) (int64, error) {
	cookie, err := r.Cookie("user")
	if err != nil {
		return 0, err
	}
	user, err := h.users.GetUserByName(r.Context(), cookie.Value)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func (h *FriendshipAPIHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	q := r.URL.Query().Get("q")
	users, err := h.friendships.SearchUsers(r.Context(), q, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if users == nil {
		users = []entity.User{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *FriendshipAPIHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		ReceiverID int64 `json:"receiver_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.friendships.SendRequest(r.Context(), userID, req.ReceiverID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *FriendshipAPIHandler) AcceptRequest(w http.ResponseWriter, r *http.Request) {
	_, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		FriendshipID int64 `json:"friendship_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.friendships.AcceptRequest(r.Context(), req.FriendshipID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *FriendshipAPIHandler) DeclineRequest(w http.ResponseWriter, r *http.Request) {
	_, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var req struct {
		FriendshipID int64 `json:"friendship_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.friendships.DeclineRequest(r.Context(), req.FriendshipID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *FriendshipAPIHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	_, err := h.getCurrentUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.URL.Query().Get("id")
	friendshipID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.friendships.RemoveFriend(r.Context(), friendshipID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
