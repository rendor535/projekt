package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type UserMediaHandler struct {
	Repo *repositories.UserMediaRepository
	DB   *sql.DB // Dodajemy bezpośredni dostęp do bazy
}

func NewUserMediaHandler(repo *repositories.UserMediaRepository, db *sql.DB) *UserMediaHandler {
	return &UserMediaHandler{Repo: repo, DB: db}
}

// Funkcja pomocnicza do konwersji YouTube URL
func extractYouTubeID(url string) string {
	patterns := []string{
		`(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)`,
		`youtube\.com\/v\/([^&\n?#]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

func (h *UserMediaHandler) Upload(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)

	var media models.UserMedia
	if err := c.ShouldBindJSON(&media); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sprawdź czy URL to YouTube
	if media.FileType == "video" {
		youtubeID := extractYouTubeID(media.FileURL)
		if youtubeID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid YouTube URL"})
			return
		}
		media.FileURL = "https://www.youtube.com/embed/" + youtubeID
	}

	media.UploaderID = userID
	media.IsRead = false

	if err := h.Repo.Create(&media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload media"})
		return
	}

	c.JSON(http.StatusCreated, media)
}

func (h *UserMediaHandler) GetSharedMedia(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)

	media, err := h.Repo.GetSharedMedia(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get shared media"})
		return
	}

	unreadCount, _ := h.Repo.GetUnreadCount(userID)

	c.JSON(http.StatusOK, gin.H{
		"media":        media,
		"unread_count": unreadCount,
	})
}

// Pobierz wszystkich użytkowników (dla trenera) - bez UserRepository
func (h *UserMediaHandler) GetUsers(c *gin.Context) {
	trainerID := c.GetInt("user_id")

	// Sprawdź rolę (opcjonalnie, jeśli potrzebna dodatkowa walidacja)
	role, err := h.Repo.GetUserRole(trainerID)
	if err != nil || role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	users, err := h.Repo.GetTrainerTrainees(trainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserMediaHandler) GetUploadedMedia(c *gin.Context) {
	trainerID, _ := c.MustGet("user_id").(int)

	media, err := h.Repo.GetUploadedMedia(trainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get uploaded media"})
		return
	}

	c.JSON(http.StatusOK, media)
}

func (h *UserMediaHandler) DeleteMedia(c *gin.Context) {
	mediaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	trainerID, _ := c.MustGet("user_id").(int)

	if err := h.Repo.DeleteMedia(mediaID, trainerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

func (h *UserMediaHandler) MarkAsRead(c *gin.Context) {
	mediaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	if err := h.Repo.MarkAsRead(mediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}
func (h *UserMediaHandler) GetExerciseVideos(c *gin.Context) {
	log.Println("🎬 GetExerciseVideos: START")

	// Sprawdź czy użytkownik to trener
	userID := c.GetInt("user_id")
	log.Printf("🎬 User ID: %d", userID)

	role, err := h.Repo.GetUserRole(userID)
	log.Printf("🎬 User role: %s, error: %v", role, err)

	if err != nil || role != "trainer" {
		log.Printf("❌ Access denied - role: %s, error: %v", role, err)
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	log.Println("🎬 Pobieranie filmów z bazy...")
	videos, err := h.Repo.GetExerciseVideos()
	log.Printf("🎬 Liczba filmów: %d, error: %v", len(videos), err)

	if err != nil {
		log.Printf("❌ Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exercise videos"})
		return
	}

	log.Printf("🎬 Zwracam filmy: %+v", videos)
	c.JSON(http.StatusOK, videos)
}

// Dodaj film z biblioteki do użytkownika (bez zmian)
func (h *UserMediaHandler) AddExerciseVideoToUser(c *gin.Context) {
	trainerID := c.GetInt("user_id")

	// Sprawdź czy użytkownik to trener
	role, err := h.Repo.GetUserRole(trainerID)
	if err != nil || role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var request struct {
		ExerciseVideoID int `json:"exercise_video_id" binding:"required"`
		RecipientID     int `json:"recipient_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Repo.CreateFromExerciseVideo(trainerID, request.RecipientID, request.ExerciseVideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign video"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Video assigned successfully"})
}
