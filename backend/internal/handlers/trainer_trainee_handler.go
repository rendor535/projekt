package handlers

import (
	"backend/internal/domain/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TrainerTraineeHandler struct {
	repo *repositories.TrainerTraineeRepository
}

func NewTrainerTraineeHandler(repo *repositories.TrainerTraineeRepository) *TrainerTraineeHandler {
	return &TrainerTraineeHandler{repo: repo}
}

func (h *TrainerTraineeHandler) AssignTrainee(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	role := c.MustGet("role").(string)

	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can assign trainees"})
		return
	}

	var req struct {
		TraineeID int `json:"trainee_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.AssignTrainee(trainerID, req.TraineeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign trainee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trainee assigned successfully"})
}

func (h *TrainerTraineeHandler) GetTrainees(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	role := c.MustGet("role").(string)

	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can view trainees"})
		return
	}

	// Add debug logs
	fmt.Printf("Fetching trainees for trainer ID: %d\n", trainerID)

	trainees, err := h.repo.GetTrainerTrainees(trainerID)
	if err != nil {
		// Log the specific error
		fmt.Printf("Error fetching trainees: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trainees: " + err.Error()})
		return
	}

	fmt.Printf("Found %d trainees\n", len(trainees))
	c.JSON(http.StatusOK, trainees)
}
func (h *TrainerTraineeHandler) RemoveTrainee(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	role := c.MustGet("role").(string)

	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can remove trainees"})
		return
	}

	traineeID, err := strconv.Atoi(c.Param("traineeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	if err := h.repo.RemoveTrainee(trainerID, traineeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove trainee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trainee removed successfully"})
}

func (h *TrainerTraineeHandler) GetMyTrainer(c *gin.Context) {
	traineeID := c.MustGet("user_id").(int)

	trainer, err := h.repo.GetTraineeTrainer(traineeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trainer"})
		return
	}

	if trainer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No trainer assigned"})
		return
	}

	c.JSON(http.StatusOK, trainer)
}
