package handlers

import (
	"backend/internal/domain/repositories"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TrainingPlanHandler struct {
	repo repositories.TrainingPlanRepository
}

func NewTrainingPlanHandler(repo repositories.TrainingPlanRepository) *TrainingPlanHandler {
	return &TrainingPlanHandler{repo: repo}
}

// GetCalendarEvents pobiera wszystkie wydarzenia kalendarza dla użytkownika
func (h *TrainingPlanHandler) GetCalendarEvents(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Brak informacji o roli użytkownika"})
		return
	}

	fmt.Printf("🔍 GetCalendarEvents - UserID: %v, Role: %v\n", userID, role)

	var events []map[string]interface{}
	var err error

	switch role {
	case "trainer":
		events, err = h.repo.GetTrainerCalendarEvents(userID.(int))
	case "trainee":
		events, err = h.repo.GetTraineeCalendarEvents(userID.(int))
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Nieznana rola użytkownika"})
		return
	}

	if err != nil {
		fmt.Printf("❌ Błąd pobierania wydarzeń kalendarza: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas pobierania wydarzeń: " + err.Error()})
		return
	}

	fmt.Printf("✅ Znaleziono %d wydarzeń dla użytkownika %v\n", len(events), userID)
	c.JSON(http.StatusOK, events)
}

// GetTrainerPlans pobiera plany trenera
func (h *TrainingPlanHandler) GetTrainerPlans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tylko trenerzy mogą przeglądać swoje plany"})
		return
	}

	fmt.Printf("🔍 GetTrainerPlans - UserID: %v\n", userID)

	plans, err := h.repo.GetTrainerPlans(userID.(int))
	if err != nil {
		fmt.Printf("❌ Błąd pobierania planów trenera: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas pobierania planów: " + err.Error()})
		return
	}

	fmt.Printf("✅ Znaleziono %d planów dla trenera %v\n", len(plans), userID)
	c.JSON(http.StatusOK, plans)
}

// GetTraineePlans pobiera plany podopiecznego
func (h *TrainingPlanHandler) GetTraineePlans(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role != "trainee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tylko podopieczni mogą przeglądać swoje plany"})
		return
	}

	fmt.Printf("🔍 GetTraineePlans - UserID: %v\n", userID)

	plans, err := h.repo.GetTraineePlans(userID.(int))
	if err != nil {
		fmt.Printf("❌ Błąd pobierania planów podopiecznego: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas pobierania planów: " + err.Error()})
		return
	}

	fmt.Printf("✅ Znaleziono %d planów dla podopiecznego %v\n", len(plans), userID)
	c.JSON(http.StatusOK, plans)
}

// GetExercises pobiera listę dostępnych ćwiczeń
func (h *TrainingPlanHandler) GetExercises(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	fmt.Printf("🔍 GetExercises - UserID: %v\n", userID)

	exercises, err := h.repo.GetExercises()
	if err != nil {
		fmt.Printf("❌ Błąd pobierania ćwiczeń: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas pobierania ćwiczeń: " + err.Error()})
		return
	}

	fmt.Printf("✅ Znaleziono %d ćwiczeń\n", len(exercises))
	c.JSON(http.StatusOK, exercises)
}

// GetTrainees pobiera listę podopiecznych trenera
func (h *TrainingPlanHandler) GetTrainees(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tylko trenerzy mogą przeglądać podopiecznych"})
		return
	}

	fmt.Printf("🔍 GetTrainees - UserID: %v\n", userID)

	trainees, err := h.repo.GetTrainerTrainees(userID.(int))
	if err != nil {
		fmt.Printf("❌ Błąd pobierania podopiecznych: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas pobierania podopiecznych: " + err.Error()})
		return
	}

	fmt.Printf("✅ Znaleziono %d podopiecznych dla trenera %v\n", len(trainees), userID)
	c.JSON(http.StatusOK, trainees)
}

// CreatePlan tworzy nowy plan treningowy
func (h *TrainingPlanHandler) CreatePlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tylko trenerzy mogą tworzyć plany"})
		return
	}

	var planData map[string]interface{}
	if err := c.ShouldBindJSON(&planData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe dane: " + err.Error()})
		return
	}

	// Dodaj ID trenera do danych planu
	planData["trainer_id"] = userID.(int)

	fmt.Printf("🔍 CreatePlan - UserID: %v, Data: %+v\n", userID, planData)

	planID, err := h.repo.CreateTrainingPlan(planData)
	if err != nil {
		fmt.Printf("❌ Błąd tworzenia planu: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas tworzenia planu: " + err.Error()})
		return
	}

	fmt.Printf("✅ Plan utworzony z ID: %d\n", planID)
	c.JSON(http.StatusCreated, gin.H{"id": planID, "message": "Plan treningowy został utworzony"})
}

// UpdatePlan aktualizuje istniejący plan treningowy - NOWA METODA! 🆕
func (h *TrainingPlanHandler) UpdatePlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tylko trenerzy mogą edytować plany"})
		return
	}

	planIDStr := c.Param("id")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowy ID planu"})
		return
	}

	var planData map[string]interface{}
	if err := c.ShouldBindJSON(&planData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowe dane: " + err.Error()})
		return
	}

	// Sprawdź czy plan należy do tego trenera
	//details, err := h.repo.GetPlanDetails(planID, userID.(int))
	if err != nil {
		fmt.Printf("❌ Błąd sprawdzania uprawnień do planu: %v\n", err)
		c.JSON(http.StatusForbidden, gin.H{"error": "Plan nie znaleziony lub brak uprawnień"})
		return
	}

	// Dodaj ID trenera do danych (zabezpieczenie)
	planData["trainer_id"] = userID.(int)

	fmt.Printf("🔍 UpdatePlan - UserID: %v, PlanID: %d, Data: %+v\n", userID, planID, planData)

	err = h.repo.UpdateTrainingPlan(planID, planData)
	if err != nil {
		fmt.Printf("❌ Błąd aktualizacji planu: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas aktualizacji planu: " + err.Error()})
		return
	}

	fmt.Printf("✅ Plan %d zaktualizowany\n", planID)
	c.JSON(http.StatusOK, gin.H{"message": "Plan treningowy został zaktualizowany"})
}

// DeletePlan usuwa plan treningowy
func (h *TrainingPlanHandler) DeletePlan(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	role, exists := c.Get("role")
	if !exists || role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Tylko trenerzy mogą usuwać plany"})
		return
	}

	planIDStr := c.Param("id")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowy ID planu"})
		return
	}

	fmt.Printf("🔍 DeletePlan - UserID: %v, PlanID: %d\n", userID, planID)

	err = h.repo.DeleteTrainingPlan(planID, userID.(int))
	if err != nil {
		fmt.Printf("❌ Błąd usuwania planu: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas usuwania planu: " + err.Error()})
		return
	}

	fmt.Printf("✅ Plan %d usunięty\n", planID)
	c.JSON(http.StatusOK, gin.H{"message": "Plan treningowy został usunięty"})
}

// GetPlanDetails pobiera szczegóły planu
func (h *TrainingPlanHandler) GetPlanDetails(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Użytkownik nie uwierzytelniony"})
		return
	}

	planIDStr := c.Param("id")
	planID, err := strconv.Atoi(planIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nieprawidłowy ID planu"})
		return
	}

	fmt.Printf("🔍 GetPlanDetails - UserID: %v, PlanID: %d\n", userID, planID)

	details, err := h.repo.GetPlanDetails(planID, userID.(int))
	if err != nil {
		fmt.Printf("❌ Błąd pobierania szczegółów planu: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Błąd podczas pobierania szczegółów: " + err.Error()})
		return
	}

	fmt.Printf("✅ Szczegóły planu %d pobrane\n", planID)
	c.JSON(http.StatusOK, details)
}
