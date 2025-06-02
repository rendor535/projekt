package models

import "time"

type TrainerTrainee struct {
	ID        int       `json:"id" db:"id"`
	TrainerID int       `json:"trainer_id" db:"trainer_id"`
	TraineeID int       `json:"trainee_id" db:"trainee_id"`
	Status    string    `json:"status" db:"status"` // active, inactive, pending
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type TraineeWithDetails struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	IsOnline       bool      `json:"is_online"`
	LastActivity   time.Time `json:"last_activity"`
	UnreadMessages int       `json:"unread_messages"`
	CurrentPlan    *string   `json:"current_plan"`
	Goal           *string   `json:"goal"`
	Progress       int       `json:"progress"`
}
