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

type Realtime interface {
	Online(userID string, conn *websocket.Conn) error
	Offline(userID string, conn *websocket.Conn) error
	SendTo(userID string, data any) error
}

type LiveChat struct {
	realtimeSvc Realtime
	messageSvc  module.MessageService
	authSvc     auth.Jwt
}

func NewLiveChat(realtimeSvc Realtime, messageSvc module.MessageService, authSvc auth.Jwt) *LiveChat {
	return &LiveChat{
		realtimeSvc: realtimeSvc,
		messageSvc:  messageSvc,
		authSvc:     authSvc,
	}
}

func (cr *LiveChat) Listen(userID string, conn *websocket.Conn) {
	// Register user
	authenticated := false

	defer cr.realtimeSvc.Offline(userID, conn)

	for {
		var wsPayload model.Envelope
		if err := conn.ReadJSON(&wsPayload); err != nil {
			log.Errorf("error reading websocket payload: %v", err)
			return
		}

		switch wsPayload.Type {
		case "send_message":
			if !authenticated {
				conn.WriteJSON(map[string]string{"type": "ERROR", "error": "unauthorized"})
				continue
			}

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

			// ส่ง message ให้ receiver แบบ realtime
			if err := cr.realtimeSvc.SendTo(msg.ReceiverID, msg); err != nil {
				log.Errorf("failed to send to receiver id {%s}: %v", msg.ReceiverID, err)
				conn.WriteJSON(map[string]string{"type": "MSG_FAIL", "error": "failed to send message"})
				continue
			}

		case "authorization":
			// เช็คว่า authenticate ไปรึยัง
			if authenticated {
				conn.WriteJSON(map[string]string{"type": "ERROR", "error": "already authenticated"})
				conn.Close()
				return
			}

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
			claims, err := cr.authSvc.VerifyToken(tokenStr)
			if err != nil {
				authenticated = false
				log.Errorf("[Websocket]: invalid token or token expire %v", err)
				conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "token expire"})
				conn.Close()
				return
			}

			// สร้าง claims เเละ set user online
			userID := claims.Subject
			if userID == "" {
				log.Errorf("[Websocket]: userID is nil {%s}", userID)
				conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "invalid token"})
				conn.Close()
				return
			}

			// check roles ที่ allow ทั้งหมดเทียบกับ userRole
			var allowedRoles = []string{"USER", "SELLER", "ADMIN"}
			if cr.authSvc.CheckAllowRoles(userID, allowedRoles) {
				log.Info("[Websocket]: Role pass")
				authenticated = true
				cr.realtimeSvc.Online(userID, conn)
				continue
			}

			// role ไม่ตรงปิด connection
			log.Error("[Websocket]: check role not pass")
			conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "role not pass"})
			conn.Close()
			return

		// case "Refresh_token":
		// 	if !authenticated {
		// 		log.Error("[Chat_Refresh_Token]: err invalid format", )
		// 		conn.WriteJSON(map[string]string{"type": "AUTH_FAIL", "error": "invalid format"})
		// 	}

		default:
			log.Warnf("unknown message type: %s", wsPayload.Type)
		}
	}
}
