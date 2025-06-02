package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	DB *sql.DB
}

func NewNotificationHandler(db *sql.DB) *NotificationHandler {
	return &NotificationHandler{DB: db}
}

// GetUnreadReportsCount - liczba nieprzeczytanych raportów od konkretnego podopiecznego
func (h *NotificationHandler) GetUnreadReportsCount(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	traineeID, err := strconv.Atoi(c.Param("traineeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	// Najpierw sprawdź czy trainee należy do trenera
	var relationshipExists int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM trainer_trainee 
		WHERE trainer_id = $1 AND trainee_id = $2
	`, trainerID, traineeID).Scan(&relationshipExists)

	if err != nil || relationshipExists == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var count int
	// Poprawka: używamy client_id zamiast user_id
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM progress_reports pr
		WHERE pr.client_id = $1 
		AND pr.created_at > COALESCE((
			SELECT last_reports_check FROM trainer_trainee 
			WHERE trainer_id = $2 AND trainee_id = $1
		), '1970-01-01')
	`, traineeID, trainerID).Scan(&count)

	if err != nil {
		log.Printf("Error getting unread reports count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetUnreadMediaCount - liczba nieprzeczytanych plików multimedialnych
func (h *NotificationHandler) GetUnreadMediaCount(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	traineeID, err := strconv.Atoi(c.Param("traineeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	// Sprawdź dostęp
	var relationshipExists int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM trainer_trainee 
		WHERE trainer_id = $1 AND trainee_id = $2
	`, trainerID, traineeID).Scan(&relationshipExists)

	if err != nil || relationshipExists == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var count int
	// Sprawdź czy tabela user_media istnieje recipent albo uploader
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM user_media um
		WHERE um.recipient_id = $1
		AND um.created_at > COALESCE((
			SELECT last_media_check FROM trainer_trainee 
			WHERE trainer_id = $2 AND trainee_id = $1
		), '1970-01-01')
	`, traineeID, trainerID).Scan(&count)

	if err != nil {
		log.Printf("Error getting unread media count (table may not exist): %v", err)
		// Jeśli tabela user_media nie istnieje, zwróć 0
		c.JSON(http.StatusOK, gin.H{"count": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetPendingDietFeedbackCount - liczba oczekujących opinii o diecie
func (h *NotificationHandler) GetPendingDietFeedbackCount(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	traineeID, err := strconv.Atoi(c.Param("traineeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	// Sprawdź dostęp
	var relationshipExists int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM trainer_trainee 
		WHERE trainer_id = $1 AND trainee_id = $2
	`, trainerID, traineeID).Scan(&relationshipExists)

	if err != nil || relationshipExists == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var count int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM diet_plans dp
		WHERE dp.client_id = $1 AND dp.trainer_id = $2
		AND COALESCE(dp.feedback_request, false) = true 
		AND COALESCE(dp.feedback_given, false) = false
	`, traineeID, trainerID).Scan(&count)

	if err != nil {
		log.Printf("Error getting diet feedback count (columns may not exist): %v", err)
		// Zwróć 0 jeśli kolumny nie istnieją
		c.JSON(http.StatusOK, gin.H{"count": 0})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetTrainerPendingReportsCount - wszystkie oczekujące raporty dla trenera
func (h *NotificationHandler) GetTrainerPendingReportsCount(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)

	var count int
	// Poprawka: używamy client_id zamiast user_id
	err := h.DB.QueryRow(`
		SELECT COUNT(*) FROM progress_reports pr
		JOIN trainer_trainee tt ON pr.client_id = tt.trainee_id
		WHERE tt.trainer_id = $1 
		AND pr.created_at > COALESCE(tt.last_reports_check, '1970-01-01')
	`, trainerID).Scan(&count)

	if err != nil {
		log.Printf("Error getting trainer pending reports count: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending reports count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
