package ws

import (
	"context"
	"encoding/json"
	"forum1/internal/service"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMsgSize = 4096
)

type Client struct {
	UserID int64
	Conn   *websocket.Conn
	Send   chan []byte
}

type WSMessage struct {
	Type       string `json:"type"`
	ReceiverID int64  `json:"receiver_id"`
	Content    string `json:"content"`
}

type WSResponse struct {
	Type      string `json:"type"`
	MessageID int64  `json:"message_id,omitempty"`
	SenderID  int64  `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (c *Client) ReadPump(hub *Hub, chatService service.ChatService) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMsgSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		if msg.Type == "message" && msg.Content != "" && msg.ReceiverID > 0 {
			saved, err := chatService.SendMessage(context.Background(), c.UserID, msg.ReceiverID, msg.Content)
			if err != nil {
				errResp, _ := json.Marshal(WSResponse{Type: "error", Content: err.Error(), SenderID: c.UserID})
				c.Send <- errResp
				continue
			}

			resp := WSResponse{
				Type:      "message",
				MessageID: saved.ID,
				SenderID:  c.UserID,
				Content:   saved.Content,
				CreatedAt: saved.CreatedAt.Format(time.RFC3339),
			}
			payload, _ := json.Marshal(resp)

			hub.SendToUser(msg.ReceiverID, payload)
			hub.SendToUser(c.UserID, payload)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
