package handlers

import (
	"backend/internal/domain/models"
	"backend/internal/domain/repositories"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type TrainerRequestHandler struct {
	Repo *repositories.TrainerRequestRepository
}

func NewTrainerRequestHandler(repo *repositories.TrainerRequestRepository) *TrainerRequestHandler {
	return &TrainerRequestHandler{Repo: repo}
}

func (h *TrainerRequestHandler) SearchTrainers(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query must be at least 2 characters"})
		return
	}

	trainers, err := h.Repo.SearchTrainers(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search trainers"})
		return
	}

	c.JSON(http.StatusOK, trainers)
}
func (h *TrainerRequestHandler) SendRequest(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)
	role := c.GetString("role")

	if role != "trainee" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only clients can send trainer requests"})
		return
	}

	var req models.TrainerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a pending request already exists
	existingRequests, err := h.Repo.GetByClientIDAndTrainerID(userID, req.TrainerID)
	if err == nil && existingRequests != nil {
		// If there's an existing request, check its status
		if existingRequests.Status == "pending" {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Masz już aktywną prośbę do tego trenera",
				"request": existingRequests,
			})
			return
		} else if existingRequests.Status == "rejected" || existingRequests.Status == "cancelled" {
			// If it was rejected or cancelled, delete it to make way for a new request
			if err := h.Repo.Delete(existingRequests.ID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Nie można usunąć wcześniejszej prośby"})
				return
			}
		}
	}

	req.ClientID = userID
	req.Status = "pending"
	fmt.Printf("[DEBUG] Tworzenie TrainerRequest: client_id=%v, trainer_id=%v, userID from context=%v\n", req.ClientID, req.TrainerID, userID)

	if err := h.Repo.Create(&req); err != nil {
		fmt.Println("Błąd przy tworzeniu TrainerRequest:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}

	c.JSON(http.StatusCreated, req)
}
func (h *TrainerRequestHandler) GetMyRequests(c *gin.Context) {
	userID, _ := c.MustGet("user_id").(int)
	fmt.Printf("Getting requests for user ID: %d\n", userID)

	// Use a safer SQL query with error handling
	query := `
        SELECT tr.id, tr.client_id, tr.trainer_id, tr.status, tr.message, tr.created_at, tr.updated_at,
               u.username as trainer_username
        FROM trainer_requests tr
        JOIN users u ON tr.trainer_id = u.id
        WHERE tr.client_id = $1
    `

	rows, err := h.Repo.DB.Query(query, userID)
	if err != nil {
		fmt.Printf("Error querying trainer requests: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}
	defer rows.Close()

	var requests []gin.H
	for rows.Next() {
		var id, clientID, trainerID int
		var status, message, trainerUsername string
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&id, &clientID, &trainerID, &status, &message, &createdAt, &updatedAt, &trainerUsername); err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		requests = append(requests, gin.H{
			"id":               id,
			"client_id":        clientID,
			"trainer_id":       trainerID,
			"status":           status,
			"message":          message,
			"created_at":       createdAt,
			"updated_at":       updatedAt,
			"trainer_username": trainerUsername,
		})
	}

	c.JSON(http.StatusOK, requests)
}
func (h *TrainerRequestHandler) CancelRequest(c *gin.Context) {
	requestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	if err := h.Repo.Delete(requestID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request cancelled successfully"})
}

func (h *TrainerRequestHandler) GetMyTrainer(c *gin.Context) {
	traineeID, _ := c.MustGet("user_id").(int)
	fmt.Printf("Looking up trainer for trainee ID: %d\n", traineeID)

	// First try to get trainer from trainer_trainee relationship
	trainerRepo := repositories.NewTrainerTraineeRepository(h.Repo.DB)
	trainer, err := trainerRepo.GetTraineeTrainer(traineeID)

	if err != nil {
		fmt.Printf("Error getting trainer from trainer_trainee: %v\n", err)
		// Only return error if it's not "no rows"
		if err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
	} else if trainer != nil {
		// Found trainer in trainer_trainee
		c.JSON(http.StatusOK, trainer)
		return
	}

	// As fallback, check approved requests
	fmt.Printf("No trainer found in trainer_trainee, checking approved requests\n")

	approvedRequest, err := h.Repo.GetApprovedRequest(traineeID)
	if err != nil {
		fmt.Printf("Error getting approved request: %v\n", err)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No approved trainer found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trainer: " + err.Error()})
		return
	}

	// Get trainer details from user table
	var trainerUser models.User // Changed variable name from 'trainer' to 'trainerUser'
	query := `SELECT id, username, role FROM users WHERE id = $1`
	err = h.Repo.DB.QueryRow(query, approvedRequest.TrainerID).Scan(
		&trainerUser.ID, &trainerUser.Username, &trainerUser.Role)

	if err != nil {
		fmt.Printf("Error getting trainer details: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trainer details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, trainerUser)
}
func (h *TrainerRequestHandler) GetRequestsForMe(c *gin.Context) {
	trainerID, _ := c.MustGet("user_id").(int)
	role := c.GetString("role")

	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can access trainer requests"})
		return
	}

	requests, err := h.Repo.GetByTrainerID(trainerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}

// UpdateRequestStatus handles approving or rejecting a trainer request
//
//	func (h *TrainerRequestHandler) UpdateRequestStatus(c *gin.Context) {
//		requestID, err := strconv.Atoi(c.Param("id"))
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
//			return
//		}
//
//		trainerID, _ := c.MustGet("user_id").(int)
//		role := c.GetString("role")
//
//		if role != "trainer" {
//			c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can update request status"})
//			return
//		}
//
//		// Verify that this request is actually addressed to this trainer
//		request, err := h.Repo.GetByID(requestID)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get request details"})
//			return
//		}
//
//		// Make sure this trainer is the one who received the request
//		if request.TrainerID != trainerID {
//			c.JSON(http.StatusForbidden, gin.H{"error": "You can only update requests sent to you"})
//			return
//		}
//
//		var req struct {
//			Status string `json:"status" binding:"required"`
//		}
//
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//			return
//		}
//
//		if req.Status != "approved" && req.Status != "rejected" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'approved' or 'rejected'"})
//			return
//		}
//
//		if err := h.Repo.UpdateStatus(requestID, req.Status); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{"message": "Request status updated successfully"})
//	}
//
//	func (h *TrainerRequestHandler) UpdateRequestStatus(c *gin.Context) {
//		requestID, err := strconv.Atoi(c.Param("id"))
//		if err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
//			return
//		}
//
//		trainerID, _ := c.MustGet("user_id").(int)
//		role := c.GetString("role")
//
//		if role != "trainer" {
//			c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can update request status"})
//			return
//		}
//
//		// Get the request to verify ownership
//		request, err := h.Repo.GetByID(requestID)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get request details"})
//			return
//		}
//
//		// Make sure this trainer is the one who received the request
//		if request.TrainerID != trainerID {
//			c.JSON(http.StatusForbidden, gin.H{"error": "You can only update requests sent to you"})
//			return
//		}
//
//		var req struct {
//			Status string `json:"status" binding:"required"`
//		}
//
//		if err := c.ShouldBindJSON(&req); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//			return
//		}
//
//		if req.Status != "approved" && req.Status != "rejected" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'approved' or 'rejected'"})
//			return
//		}
//
//		// Start a transaction to ensure data consistency
//		tx, err := h.Repo.DB.Begin()
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
//			return
//		}
//		defer tx.Rollback() // Will be ignored if transaction is committed
//
//		// Update request status
//		if err := h.Repo.UpdateStatusTx(tx, requestID, req.Status); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
//			return
//		}
//
//		// If approved, create trainer-trainee relationship
//		if req.Status == "approved" {
//			// Insert into trainer_trainee table
//			_, err := tx.Exec(
//				`INSERT INTO trainer_trainee (trainer_id, trainee_id, created_at, updated_at)
//				VALUES ($1, $2, NOW(), NOW())`,
//				trainerID, request.ClientID)
//
//			if err != nil {
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trainer-trainee relationship"})
//				return
//			}
//		}
//
//		// Commit the transaction
//		if err := tx.Commit(); err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{"message": "Request status updated successfully", "status": req.Status})
//	}
func (h *TrainerRequestHandler) UpdateRequestStatus(c *gin.Context) {
	requestID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	trainerID, _ := c.MustGet("user_id").(int)
	role := c.GetString("role")

	if role != "trainer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only trainers can update request status"})
		return
	}

	// Get the request to verify ownership
	request, err := h.Repo.GetByID(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get request details"})
		return
	}

	// Make sure this trainer is the one who received the request
	if request.TrainerID != trainerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update requests sent to you"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Status != "approved" && req.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be 'approved' or 'rejected'"})
		return
	}

	// Start a transaction to ensure data consistency
	tx, err := h.Repo.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Update request status
	if err := h.Repo.UpdateStatusTx(tx, requestID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
		return
	}

	// If approved, create trainer-trainee relationship
	if req.Status == "approved" {
		// First, check if relationship already exists
		var existingRelationship int
		checkQuery := `SELECT COUNT(*) FROM trainer_trainee WHERE trainer_id = $1 AND trainee_id = $2`
		err := tx.QueryRow(checkQuery, trainerID, request.ClientID).Scan(&existingRelationship)

		if err != nil {
			fmt.Printf("Error checking existing relationship: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing relationship"})
			return
		}

		if existingRelationship > 0 {
			// Relationship already exists, no need to create another
			fmt.Printf("Relationship already exists between trainer %d and trainee %d\n", trainerID, request.ClientID)
		} else {
			// Insert into trainer_trainee table
			_, err := tx.Exec(
				`INSERT INTO trainer_trainee (trainer_id, trainee_id, created_at, updated_at) 
				VALUES ($1, $2, NOW(), NOW())`,
				trainerID, request.ClientID)

			if err != nil {
				fmt.Printf("Error creating trainer-trainee relationship: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trainer-trainee relationship: " + err.Error()})
				return
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Request status updated successfully", "status": req.Status})
}
