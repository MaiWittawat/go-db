package realtime

import (
	"context"
	"encoding/json"
	"errors"
	"go-rebuild/internal/auth"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type ChatRealtime struct {
	realtimeSvc Realtime
	messageSvc  module.MessageService
	authSvc     auth.Jwt
}

func NewChatRealtime(realtimeSvc Realtime, messageSvc module.MessageService, authSvc auth.Jwt) *ChatRealtime {
	return &ChatRealtime{
		realtimeSvc: realtimeSvc,
		messageSvc:  messageSvc,
		authSvc:     authSvc,
	}
}

func (cr *ChatRealtime) Listen(userID string, conn *websocket.Conn) {
	// Register user
	defer cr.realtimeSvc.Offline(userID, conn)

	for {
		var wsPayload model.Envelope
		if err := conn.ReadJSON(&wsPayload); err != nil {
			log.Errorf("error reading websocket payload: %v", err)
			return
		}

		switch wsPayload.Type {
		case "send_message":
			var msg model.MessageReq
			if err := json.Unmarshal(wsPayload.Payload, &msg); err != nil {
				log.Error("invalid send_message payload")
				conn.WriteJSON(map[string]string{"type": "MSG_FAIL", "error": "failed to send message"})
				continue
			}

			// Save to DB
			if err := cr.messageSvc.Save(context.Background(), &msg); err != nil {
				log.Errorf("failed to save message: %v", err)
				conn.WriteJSON(map[string]string{"type": "MSG_FAIL", "error": "failed to send message"})
				continue
			}

			// ส่งให้ receiver แบบ realtime
			if err := cr.realtimeSvc.SendTo(msg.ReceiverID, msg); err != nil {
				log.Errorf("failed to send to receiver id {%s}: %v", msg.ReceiverID, err)
				conn.WriteJSON(map[string]string{"type": "MSG_FAIL", "error": "failed to send message"})
				continue
			}
			

		case "Authorization":
			var payload struct {
				Token string `json:"token"`
			}
			if err := json.Unmarshal(wsPayload.Payload, &payload); err != nil {
				log.Errorf("[Websocket]: err invalid format %v", err)
				conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "invalid format"})
				conn.Close()
				return
			}

			// verify token
			if !strings.HasPrefix(payload.Token, "Bearer ") {
				log.Errorf("[Websocket]: invalid token %v", errors.New("no bearer token"))
				conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "invalid token format"})
				conn.Close()
				return
			}

			tokenStr := strings.TrimPrefix(payload.Token, "Bearer ")
			if err := cr.authSvc.VerifyToken(tokenStr); err != nil {
				log.Errorf("[Websocket]: no bearer token %v", err)
				conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "unauthorized"})
				conn.Close()
				return
			}

			// set user online
			userID, err := cr.authSvc.GetUserIDFromToken(tokenStr) // implement this
			if err != nil {
				log.Errorf("[Websocket]: failed to getUserIDFromToken %v", err)
				conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "unauthorized"})
				conn.Close()
				return
			}

			var allowedRoles = []string{"USER", "SELLER", "ADMIN"}
			if cr.authSvc.CheckAllowRoles(userID, allowedRoles) {
				log.Info("[Websocket]: Role pass")
				cr.realtimeSvc.Online(userID, conn)
				return
			}

			log.Error("[Websocket]: check role not pass")
			conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "role not pass"})
			conn.Close()
			return

		default:
			log.Warnf("unknown message type: %s", wsPayload.Type)
		}
	}
}
