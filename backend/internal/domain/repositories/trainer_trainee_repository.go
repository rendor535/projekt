package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
)

type TrainerTraineeRepository struct {
	DB *sql.DB
}

func NewTrainerTraineeRepository(db *sql.DB) *TrainerTraineeRepository {
	return &TrainerTraineeRepository{DB: db}
}

func (r *TrainerTraineeRepository) AssignTrainee(trainerID, traineeID int) error {
	query := `
		INSERT INTO trainer_trainee (trainer_id, trainee_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (trainer_id, trainee_id) 
		DO NOTHING
	`
	_, err := r.DB.Exec(query, trainerID, traineeID)
	return err
}

func (r *TrainerTraineeRepository) GetTrainerTrainees(trainerID int) ([]models.TraineeWithDetails, error) {
	query := `
		SELECT DISTINCT u.id, u.username, u.is_online,
		       COALESCE(MAX(m.created_at), u.created_at) as last_activity,
		       COALESCE(unread.count, 0) as unread_messages,
		       tp.name as current_plan,
		       p.goal,
		       COALESCE(ROUND(
		           CASE 
		               WHEN tp.start_date IS NOT NULL AND tp.end_date IS NOT NULL THEN
		                   ((EXTRACT(epoch FROM NOW()) - EXTRACT(epoch FROM tp.start_date::timestamp)) / 
		                    (EXTRACT(epoch FROM tp.end_date::timestamp) - EXTRACT(epoch FROM tp.start_date::timestamp))) * 100
		               ELSE 0
		           END
		       )::int, 0) as progress
		FROM trainer_trainee tt
		JOIN users u ON tt.trainee_id = u.id
		LEFT JOIN messages m ON (m.sender_id = u.id OR m.receiver_id = u.id)
		LEFT JOIN (
		    SELECT receiver_id, COUNT(*) as count 
		    FROM messages 
		    WHERE sender_id = $1 AND is_read = FALSE 
		    GROUP BY receiver_id
		) unread ON unread.receiver_id = u.id
		LEFT JOIN training_plans tp ON tp.client_id = u.id AND tp.trainer_id = $1
		LEFT JOIN profiles p ON p.user_id = u.id
		WHERE tt.trainer_id = $1
		GROUP BY u.id, u.username, u.is_online, u.created_at, unread.count, tp.name, tp.start_date, tp.end_date, p.goal
		ORDER BY last_activity DESC
	`
	rows, err := r.DB.Query(query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainees []models.TraineeWithDetails
	for rows.Next() {
		var trainee models.TraineeWithDetails
		var lastActivity sql.NullTime
		var currentPlan, goal sql.NullString
		var unreadMessages sql.NullInt64
		var progress sql.NullInt64

		err := rows.Scan(
			&trainee.ID,
			&trainee.Username,
			&trainee.IsOnline,
			&lastActivity,
			&unreadMessages,
			&currentPlan,
			&goal,
			&progress,
		)

		if err != nil {
			return nil, err
		}

		// Handle null values properly
		if lastActivity.Valid {
			trainee.LastActivity = lastActivity.Time
		}
		if unreadMessages.Valid {
			trainee.UnreadMessages = int(unreadMessages.Int64)
		}
		if currentPlan.Valid {
			planValue := currentPlan.String
			trainee.CurrentPlan = &planValue // Assign the address of the string
		}
		if goal.Valid {
			goalValue := goal.String
			trainee.Goal = &goalValue // Assign the address of the string
		}
		if progress.Valid {
			trainee.Progress = int(progress.Int64)
		}

		trainees = append(trainees, trainee)
	}

	return trainees, nil
}

func (r *TrainerTraineeRepository) RemoveTrainee(trainerID, traineeID int) error {
	query := `DELETE FROM trainer_trainee WHERE trainer_id = $1 AND trainee_id = $2`
	_, err := r.DB.Exec(query, trainerID, traineeID)
	return err
}

func (r *TrainerTraineeRepository) GetTraineeTrainer(traineeID int) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.role, u.is_online
		FROM trainer_trainee tt
		JOIN users u ON tt.trainer_id = u.id
		WHERE tt.trainee_id = $1
		LIMIT 1
	`

	var trainer models.User
	err := r.DB.QueryRow(query, traineeID).Scan(
		&trainer.ID,
		&trainer.Username,
		&trainer.Role,
		&trainer.IsOnline,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &trainer, nil
}
