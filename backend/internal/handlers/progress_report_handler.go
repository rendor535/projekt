package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type ProgressReportHandler struct {
	Repo *repositories.ProgressReportRepository
}

func NewProgressReportHandler(repo *repositories.ProgressReportRepository) *ProgressReportHandler {
	return &ProgressReportHandler{Repo: repo}
}

func (h *ProgressReportHandler) Create(c *gin.Context) {
	log.Printf("=== CREATE HANDLER START ===")

	// Safe context extraction
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("ERROR: user_id not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("Creating report for user ID: %v (type: %T)", userID, userID)

	var report models.ProgressReport
	if err := c.ShouldBindJSON(&report); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Report data received: %+v", report)

	// Convert userID to int safely
	var clientID int
	switch v := userID.(type) {
	case int:
		clientID = v
	case float64:
		clientID = int(v)
	default:
		log.Printf("ERROR: Invalid user_id type: %T", v)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	report.ClientID = clientID
	report.TrainerID = nil // Self-report

	log.Printf("Final report data before save: %+v", report)

	if err := h.Repo.Create(&report); err != nil {
		log.Printf("Error creating report in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report", "details": err.Error()})
		return
	}

	log.Printf("Successfully created report with ID: %d", report.ID)
	c.JSON(http.StatusCreated, report)
}

func (h *ProgressReportHandler) GetMyReports(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	reports, err := h.Repo.GetByClientID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// Nowy endpoint: Statystyki użytkownika
func (h *ProgressReportHandler) GetMyStats(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	stats, err := h.Repo.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Nowy endpoint: Dane do wykresów
func (h *ProgressReportHandler) GetChartData(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	chartData, err := h.Repo.GetChartData(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chart data"})
		return
	}

	c.JSON(http.StatusOK, chartData)
}

// Nowy endpoint dla trenera: Statystyki wszystkich podopiecznych
func (h *ProgressReportHandler) GetTrainerStats(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)

	stats, err := h.Repo.GetTrainerStats(trainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trainer stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Nowy endpoint dla trenera: Raporty konkretnego podopiecznego
func (h *ProgressReportHandler) GetTraineeReports(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	traineeID, err := strconv.Atoi(c.Param("traineeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	reports, err := h.Repo.GetTraineeReports(trainerID, traineeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trainee reports"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

func (h *ProgressReportHandler) GetLatestWeight(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	weight, err := h.Repo.GetLatestWeight(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No weight data found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weight": weight})
}
func (h *ProgressReportHandler) GetTraineeChartData(c *gin.Context) {
	trainerID := c.MustGet("user_id").(int)
	traineeID, err := strconv.Atoi(c.Param("traineeId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	// Verify trainer has access to this trainee
	hasAccess, err := h.Repo.VerifyTrainerAccess(trainerID, traineeID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	chartData, err := h.Repo.GetChartData(traineeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chart data"})
		return
	}

	c.JSON(http.StatusOK, chartData)
}
