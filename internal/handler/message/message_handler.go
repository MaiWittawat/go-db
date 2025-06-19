package message

import (
	"go-rebuild/internal/handler"
	"go-rebuild/internal/module"
	"go-rebuild/internal/realtime"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type messageHandler struct {
	liveChat   *realtime.LiveChat
	messageSvc module.MessageService
}

func NewMessageHandler(liveChat *realtime.LiveChat, messageSvc module.MessageService) handler.MessageHandler {
	return &messageHandler{
		liveChat: liveChat,
		messageSvc: messageSvc,
	}
}

func (h *messageHandler) Connect(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "websocket upgrade failed"})
		return
	}
	log.Println("in connect pass go to listen")
	go h.liveChat.Listen(userID, conn) // ไม่ block Gin handler
}

func (h *messageHandler) GetMessagesBetweenUser(c *gin.Context) {
	userID1 := c.Param("user_id1")
	userID2 := c.Param("user_id2")

	messages, err := h.messageSvc.GetMessagesBetweenUser(c.Request.Context(), userID1, userID2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get messages success", "data": messages})
}
