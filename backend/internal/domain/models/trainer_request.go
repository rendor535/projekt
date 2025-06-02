package models

import "time"

type TrainerSearchResult struct {
	ID              int      `json:"id"`
	Username        string   `json:"username"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	Specializations []string `json:"specializations"`
	Experience      int      `json:"experience"`
	IsOnline        bool     `json:"is_online"`
}
type TrainerRequest struct {
	ID        int       `json:"id"`
	ClientID  int       `json:"client_id"`
	TrainerID int       `json:"trainer_id"`
	Status    string    `json:"status"` // pending, approved, rejected, cancelled
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TrainerRequestWithDetails struct {
	TrainerRequest
	TrainerUsername       string `json:"trainer_username"`
	TrainerSpecialization string `json:"trainer_specialization"`
	TrainerExperience     int    `json:"trainer_experience"`
	ClientUsername        string `json:"client_username"`
}
