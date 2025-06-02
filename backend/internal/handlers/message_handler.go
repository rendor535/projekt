package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageRepo *repositories.MessageRepository
}

func NewMessageHandler(messageRepo *repositories.MessageRepository) *MessageHandler {
	return &MessageHandler{
		messageRepo: messageRepo,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req struct {
		ReceiverID int    `json:"receiver_id" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message := &models.ChatMessage{
		SenderID:   userID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}

	if err := h.messageRepo.Create(message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *MessageHandler) GetConversation(c *gin.Context) {
	userID := c.GetInt("user_id")
	partnerIDStr := c.Param("userID")
	limitStr := c.Query("limit")
	afterStr := c.Query("after")

	partnerID, err := strconv.Atoi(partnerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var limit int
	var after int

	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	if afterStr != "" {
		after, _ = strconv.Atoi(afterStr)
	}

	messages, err := h.messageRepo.GetConversation(userID, partnerID, limit, after)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetInt("user_id")
	senderIDStr := c.Query("sender_id")

	var count int
	var err error

	if senderIDStr != "" {
		senderID, parseErr := strconv.Atoi(senderIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
			return
		}
		count, err = h.messageRepo.GetUnreadCountFromSender(userID, senderID)
	} else {
		count, err = h.messageRepo.GetUnreadCount(userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unread count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	messageIDStr := c.Param("messageID")
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	if err := h.messageRepo.MarkAsRead(messageID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}

func (h *MessageHandler) GetConversationPartners(c *gin.Context) {
	userID := c.GetInt("user_id")

	partners, err := h.messageRepo.GetConversationPartners(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get conversation partners"})
		return
	}

	c.JSON(http.StatusOK, partners)
}
func (h *MessageHandler) MarkConversationAsRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	senderIDStr := c.Param("senderID")

	senderID, err := strconv.Atoi(senderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
		return
	}

	if err := h.messageRepo.MarkConversationAsRead(userID, senderID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark conversation as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Conversation marked as read"})
}
