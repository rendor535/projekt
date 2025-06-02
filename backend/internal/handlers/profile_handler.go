package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	Repo *repositories.ProfileRepository
	DB   *sql.DB // Dodaj pole DB do handlera
}

func NewProfileHandler(repo *repositories.ProfileRepository, db *sql.DB) *ProfileHandler {
	return &ProfileHandler{
		Repo: repo,
		DB:   db, // Przekaż db w konstruktorze
	}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	profile, err := h.Repo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (h *ProfileHandler) SaveProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	var input models.Profile
	log.Printf("GOAL DEBUG: '%s', LEN: %d", input.Goal, len(input.Goal))
	log.Printf("PEŁNY PROFIL: %+v", input)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = userID

	// Walidacja danych
	if !input.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid profile data. All fields are required."})
		return
	}

	if err := h.Repo.CreateOrUpdate(&input); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile saved",
		"profile": input,
	})
}

func (h *ProfileHandler) CheckProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	hasProfile, err := h.Repo.HasValidProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_profile": hasProfile})
}

// GetTraineeProfile - pozwala trenerowi pobrać profil podopiecznego
func (h *ProfileHandler) GetTraineeProfile(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	role := c.MustGet("role").(string)

	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can view trainee profiles"})
		return
	}

	traineeID, err := strconv.Atoi(c.Param("traineeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	// Sprawdź czy trainee jest przypisany do tego trenera
	var count int
	err = h.DB.QueryRow(`
		SELECT COUNT(*) FROM trainer_trainee 
		WHERE trainer_id = $1 AND trainee_id = $2
	`, trainerID, traineeID).Scan(&count)

	if err != nil {
		log.Printf("Error checking trainer-trainee relationship: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Trainee not assigned to this trainer"})
		return
	}

	// Pobierz profil
	profile, err := h.Repo.GetByUserID(traineeID)
	if err != nil {
		log.Printf("Error getting trainee profile: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Trainee profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}
