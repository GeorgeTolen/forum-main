package handler

import (
	"encoding/json"
	"forum1/internal/repository"
	"forum1/internal/service"
	"forum1/internal/ws"
	"forum1/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatPageHandler struct {
	chat  service.ChatService
	users repository.UserRepository
	hub   *ws.Hub
}

func NewChatPageHandler(c service.ChatService, u repository.UserRepository, h *ws.Hub) *ChatPageHandler {
	return &ChatPageHandler{chat: c, users: u, hub: h}
}

func (h *ChatPageHandler) ChatListPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := h.users.GetUserByName(r.Context(), cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	conversations, _ := h.chat.GetConversations(r.Context(), user.ID)

	data := map[string]interface{}{
		"User":          user,
		"Conversations": conversations,
	}
	utils.RenderTemplate(w, "chat_list_page.html", data)
}

func (h *ChatPageHandler) ChatPage(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := h.users.GetUserByName(r.Context(), cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	friendIDStr := mux.Vars(r)["friendID"]
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid friend ID", http.StatusBadRequest)
		return
	}

	friend, err := h.users.GetUserByID(r.Context(), friendID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	messages, _ := h.chat.GetMessages(r.Context(), user.ID, friendID, 50)
	_ = h.chat.MarkAsRead(r.Context(), friendID, user.ID)

	data := map[string]interface{}{
		"User":     user,
		"Friend":   friend,
		"Messages": messages,
	}
	utils.RenderTemplate(w, "chat_room_page.html", data)
}

type ChatAPIHandler struct {
	chat  service.ChatService
	users repository.UserRepository
	hub   *ws.Hub
}

func NewChatAPIHandler(c service.ChatService, u repository.UserRepository, h *ws.Hub) *ChatAPIHandler {
	return &ChatAPIHandler{chat: c, users: u, hub: h}
}

func (h *ChatAPIHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.users.GetUserByName(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		UserID: user.ID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump(h.hub, h.chat)
}

func (h *ChatAPIHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.users.GetUserByName(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	friendIDStr := mux.Vars(r)["friendID"]
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid friend ID", http.StatusBadRequest)
		return
	}

	messages, err := h.chat.GetMessages(r.Context(), user.ID, friendID, 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *ChatAPIHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil || cookie.Value == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	user, err := h.users.GetUserByName(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	friendIDStr := mux.Vars(r)["friendID"]
	friendID, err := strconv.ParseInt(friendIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid friend ID", http.StatusBadRequest)
		return
	}

	_ = h.chat.MarkAsRead(r.Context(), friendID, user.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
