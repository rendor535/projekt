package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TrainerProfileHandler struct {
	Repo *repositories.TrainerProfileRepository
}

func NewTrainerProfileHandler(repo *repositories.TrainerProfileRepository) *TrainerProfileHandler {
	return &TrainerProfileHandler{Repo: repo}
}

func (h *TrainerProfileHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	fmt.Printf("Getting trainer profile for user ID: %d\n", userID)

	profile, err := h.Repo.GetByUserID(userID)
	if err != nil {
		fmt.Printf("Trainer profile not found for user %d: %v\n", userID, err)
		c.JSON(http.StatusOK, gin.H{"profile": nil})
		return
	}

	if profile.IsEmpty() {
		fmt.Printf("Trainer profile is empty for user %d\n", userID)
		c.JSON(http.StatusOK, gin.H{"profile": nil})
		return
	}

	fmt.Printf("Trainer profile found for user %d\n", userID)
	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func (h *TrainerProfileHandler) SaveProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	fmt.Printf("Saving trainer profile for user ID: %d\n", userID)

	var input models.TrainerProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("JSON binding error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}

	fmt.Printf("Trainer profile input: %+v\n", input)
	input.UserID = userID
	input.IsActive = true // Domyślnie aktywny

	if !input.IsValid() {
		fmt.Printf("Invalid trainer profile data: %+v\n", input)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainer profile data. All required fields must be filled."})
		return
	}

	err := h.Repo.CreateOrUpdate(&input)
	if err != nil {
		fmt.Printf("Database error saving trainer profile: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}

	fmt.Printf("Trainer profile saved successfully for user %d\n", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Trainer profile saved",
		"profile": input,
	})
}

func (h *TrainerProfileHandler) CheckProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	hasProfile, err := h.Repo.HasValidProfile(userID)
	if err != nil {
		fmt.Printf("Error checking trainer profile for user %d: %v\n", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_profile": hasProfile})
}

func (h *TrainerProfileHandler) GetAllTrainers(c *gin.Context) {
	trainers, err := h.Repo.GetAllActiveTrainers()
	if err != nil {
		fmt.Printf("Error getting all trainers: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trainers": trainers})
}
func (h *TrainerProfileHandler) GetTrainerByID(c *gin.Context) {
	trainerIDStr := c.Param("id")
	trainerID, err := strconv.Atoi(trainerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainer ID"})
		return
	}

	trainer, err := h.Repo.GetByUserID(trainerID)
	if err != nil {
		fmt.Printf("Trainer not found for ID %d: %v\n", trainerID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Trainer not found"})
		return
	}

	if trainer.IsEmpty() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trainer profile is empty"})
		return
	}

	// Zwróć dane w formacie oczekiwanym przez frontend
	c.JSON(http.StatusOK, gin.H{
		"trainer": trainer,
		"profile": trainer, // dla kompatybilności
	})
}
func (h *TrainerProfileHandler) SearchTrainers(c *gin.Context) {
	trainerIDStr := c.Query("trainer_id")

	if trainerIDStr != "" {
		trainerID, err := strconv.Atoi(trainerIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainer ID"})
			return
		}

		trainer, err := h.Repo.GetByUserID(trainerID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Trainer not found"})
			return
		}

		if trainer.IsEmpty() {
			c.JSON(http.StatusOK, gin.H{"trainers": []models.TrainerProfile{}})
			return
		}

		c.JSON(http.StatusOK, []models.TrainerProfile{*trainer})
		return
	}

	// Jeśli brak trainer_id, zwróć wszystkich trenerów
	trainers, err := h.Repo.GetAllActiveTrainers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, trainers)
}
