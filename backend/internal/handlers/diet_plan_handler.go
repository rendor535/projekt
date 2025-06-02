package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DietPlanHandler struct {
	Repo *repositories.DietPlanRepository
}

func NewDietPlanHandler(repo *repositories.DietPlanRepository) *DietPlanHandler {
	return &DietPlanHandler{Repo: repo}
}

// ===================
// TRAINER ENDPOINTS
// ===================

func (h *DietPlanHandler) CreateDietPlan(c *gin.Context) {
	role := c.GetString("role")
	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can create diet plans"})
		return
	}

	trainerID, _ := c.MustGet("user_id").(int)

	var req models.CreateDietPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.Repo.CreateWithMeals(&req, trainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diet plan: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *DietPlanHandler) GetTrainerDietPlans(c *gin.Context) {
	role := c.GetString("role")
	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can access this endpoint"})
		return
	}

	trainerID, _ := c.MustGet("user_id").(int)

	// Opcjonalne filtry
	var clientID *int
	var isActive *bool

	if clientIDStr := c.Query("client_id"); clientIDStr != "" {
		if id, err := strconv.Atoi(clientIDStr); err == nil {
			clientID = &id
		}
	}

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &active
		}
	}

	plans, err := h.Repo.GetTrainerDietPlans(trainerID, clientID, isActive)
	if err != nil {
		fmt.Printf("GetTrainerDietPlans error: %v\n", err) // <--- DODAJ TO
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get diet plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *DietPlanHandler) UpdateDietPlan(c *gin.Context) {
	role := c.GetString("role")
	userID, _ := c.MustGet("user_id").(int)
	planID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var req models.UpdateDietPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Repo.Update(planID, &req, userID, role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update diet plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diet plan updated successfully"})
}

func (h *DietPlanHandler) DeleteDietPlan(c *gin.Context) {
	role := c.GetString("role")
	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can delete diet plans"})
		return
	}

	trainerID, _ := c.MustGet("user_id").(int)
	planID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	if err := h.Repo.Delete(planID, trainerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete diet plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diet plan deleted successfully"})
}

// ===================
// TRAINEE ENDPOINTS
// ===================

func (h *DietPlanHandler) GetTraineeDietPlans(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)

	// Opcjonalny filtr dla aktywnych planów
	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &active
		}
	}

	plans, err := h.Repo.GetTraineeDietPlans(userID, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get diet plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// ===================
// SHARED ENDPOINTS
// ===================

func (h *DietPlanHandler) GetDietPlanDetails(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)
	role := c.GetString("role")
	planID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Sprawdź uprawnienia do tego diet planu
	var checkQuery string
	if role == "trainer" {
		checkQuery = "SELECT trainer_id FROM diet_plans WHERE id = $1"
	} else {
		checkQuery = "SELECT client_id FROM diet_plans WHERE id = $1"
	}

	var ownerID int
	err = h.Repo.DB.QueryRow(checkQuery, planID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diet plan not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	meals, err := h.Repo.GetDietPlanMeals(planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get diet plan details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"meals": meals})
}

// ===================
// LEGACY ENDPOINTS (dla kompatybilności)
// ===================

func (h *DietPlanHandler) Create(c *gin.Context) {
	role := c.GetString("role")
	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can create diet plans"})
		return
	}

	trainerID, _ := c.MustGet("user_id").(int)

	var plan models.DietPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan.TrainerID = trainerID

	if err := h.Repo.Create(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create diet plan"})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

func (h *DietPlanHandler) GetMyDietPlans(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)

	plans, err := h.Repo.GetByClientID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get diet plans"})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *DietPlanHandler) GetActiveDietPlan(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)

	plan, err := h.Repo.GetActivePlan(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active diet plan found"})
		return
	}

	c.JSON(http.StatusOK, plan)
}

func (h *DietPlanHandler) ActivateDietPlan(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)
	planID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	if err := h.Repo.ActivatePlan(planID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate diet plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diet plan activated successfully"})
}
