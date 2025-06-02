package models

import "time"

type TraineeDetails struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Status         string    `json:"status"`
	LastActivity   time.Time `json:"last_activity"`
	UnreadMessages int       `json:"unread_messages"`
	CurrentPlan    string    `json:"current_plan"`
	Goal           string    `json:"goal"`
	Progress       int       `json:"progress"`
}
