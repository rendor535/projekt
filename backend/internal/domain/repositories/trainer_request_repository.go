package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
	"encoding/json"
)

type TrainerRequestRepository struct {
	DB *sql.DB
}

func NewTrainerRequestRepository(db *sql.DB) *TrainerRequestRepository {
	return &TrainerRequestRepository{DB: db}
}

func (r *TrainerRequestRepository) Create(req *models.TrainerRequest) error {
	query := `
		INSERT INTO trainer_requests (client_id, trainer_id, status, message, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at
	`
	return r.DB.QueryRow(query, req.ClientID, req.TrainerID, req.Status, req.Message).
		Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
}

func (r *TrainerRequestRepository) GetByClientID(clientID int) ([]models.TrainerRequestWithDetails, error) {
	query := `
		SELECT tr.id, tr.client_id, tr.trainer_id, tr.status, tr.message, tr.created_at, tr.updated_at,
		       u.username, u.specialization, u.experience
		FROM trainer_requests tr
		JOIN users u ON tr.trainer_id = u.id
		WHERE tr.client_id = $1
		ORDER BY tr.created_at DESC
	`
	rows, err := r.DB.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.TrainerRequestWithDetails
	for rows.Next() {
		var req models.TrainerRequestWithDetails
		err := rows.Scan(&req.ID, &req.ClientID, &req.TrainerID, &req.Status, &req.Message,
			&req.CreatedAt, &req.UpdatedAt, &req.TrainerUsername, &req.TrainerSpecialization, &req.TrainerExperience)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *TrainerRequestRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE trainer_requests SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.Exec(query, status, id)
	return err
}

func (r *TrainerRequestRepository) Delete(id int) error {
	query := `DELETE FROM trainer_requests WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *TrainerRequestRepository) GetApprovedRequest(clientID int) (*models.TrainerRequest, error) {
	query := `
        SELECT id, client_id, trainer_id, status, message, created_at, updated_at
        FROM trainer_requests
        WHERE client_id = $1 AND status = 'approved'
        LIMIT 1
    `

	var request models.TrainerRequest
	err := r.DB.QueryRow(query, clientID).Scan(
		&request.ID,
		&request.ClientID,
		&request.TrainerID,
		&request.Status,
		&request.Message,
		&request.CreatedAt,
		&request.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *TrainerRequestRepository) SearchTrainers(query string) ([]models.TrainerSearchResult, error) {
	sqlQuery := `
        SELECT 
            u.id, 
            u.username, 
            tp.first_name, 
            tp.last_name, 
            tp.specializations, 
            tp.experience, 
            u.is_online
        FROM users u
        JOIN trainer_profiles tp ON u.id = tp.user_id
        WHERE u.role = 'trainer'
          AND (
            u.username ILIKE $1
            OR tp.first_name ILIKE $1
            OR tp.last_name ILIKE $1
            OR (tp.first_name || ' ' || tp.last_name) ILIKE $1
          )
        ORDER BY tp.first_name, tp.last_name
    `
	rows, err := r.DB.Query(sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainers []models.TrainerSearchResult
	for rows.Next() {
		var trainer models.TrainerSearchResult
		var specsBytes []byte
		err := rows.Scan(
			&trainer.ID,
			&trainer.Username,
			&trainer.FirstName,
			&trainer.LastName,
			&specsBytes,
			&trainer.Experience,
			&trainer.IsOnline,
		)
		if err != nil {
			return nil, err
		}
		// Unmarshal JSONB
		_ = json.Unmarshal(specsBytes, &trainer.Specializations)
		trainers = append(trainers, trainer)
	}
	return trainers, nil
}
func (r *TrainerRequestRepository) GetByTrainerID(trainerID int) ([]models.TrainerRequestWithDetails, error) {
	query := `
		SELECT tr.id, tr.client_id, tr.trainer_id, tr.status, tr.message, tr.created_at, tr.updated_at,
		       u.username as client_username
		FROM trainer_requests tr
		JOIN users u ON tr.client_id = u.id
		WHERE tr.trainer_id = $1
		ORDER BY tr.created_at DESC
	`
	rows, err := r.DB.Query(query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.TrainerRequestWithDetails
	for rows.Next() {
		var req models.TrainerRequestWithDetails
		err := rows.Scan(&req.ID, &req.ClientID, &req.TrainerID, &req.Status, &req.Message,
			&req.CreatedAt, &req.UpdatedAt, &req.ClientUsername)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}
func (r *TrainerRequestRepository) GetByID(id int) (*models.TrainerRequest, error) {
	query := `
		SELECT id, client_id, trainer_id, status, message, created_at, updated_at
		FROM trainer_requests
		WHERE id = $1
	`

	var req models.TrainerRequest
	err := r.DB.QueryRow(query, id).Scan(
		&req.ID, &req.ClientID, &req.TrainerID, &req.Status,
		&req.Message, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &req, nil
}
func (r *TrainerRequestRepository) UpdateStatusTx(tx *sql.Tx, id int, status string) error {
	query := `UPDATE trainer_requests SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.Exec(query, status, id)
	return err
}
func (r *TrainerRequestRepository) GetByClientIDAndTrainerID(clientID, trainerID int) (*models.TrainerRequest, error) {
	query := `
		SELECT id, client_id, trainer_id, status, message, created_at, updated_at
		FROM trainer_requests
		WHERE client_id = $1 AND trainer_id = $2
	`

	var req models.TrainerRequest
	err := r.DB.QueryRow(query, clientID, trainerID).Scan(
		&req.ID, &req.ClientID, &req.TrainerID, &req.Status,
		&req.Message, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No request found
		}
		return nil, err
	}

	return &req, nil
}
