package handlers

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func GetUserIDFromContext(c *gin.Context) (int, error) {
	val, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("userID not in context")
	}
	userID, ok := val.(int)
	if !ok {
		return 0, errors.New("userID in context is not int")
	}
	return userID, nil
}

func (h *UserHandler) SetOnline(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	_, err = h.DB.Exec(`
		UPDATE users 
		SET is_online = TRUE, last_seen = NOW()
		WHERE id = $1
		`, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "db error"})
		return
	}
	//c.JSON(200, gin.H{"success": true})
}

func (h *UserHandler) SetOffline(c *gin.Context) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	_, err = h.DB.Exec("UPDATE users SET is_online = FALSE, last_seen = NOW() WHERE id = $1", userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "db error"})
		return
	}
	//c.JSON(200, gin.H{"success": true})
}
