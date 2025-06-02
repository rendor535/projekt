package models

import "time"

type TrainerProfile struct {
	UserID          int       `json:"user_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Bio             string    `json:"bio"`
	Specializations []string  `json:"specializations"`
	Experience      int       `json:"experience"` // lata doświadczenia
	Certifications  []string  `json:"certifications"`
	PricePerHour    float64   `json:"price_per_hour"`
	Languages       []string  `json:"languages"`
	ContactEmail    string    `json:"contact_email"`
	ContactPhone    string    `json:"contact_phone"`
	Location        string    `json:"location"`
	Avatar          string    `json:"avatar,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (tp *TrainerProfile) IsEmpty() bool {
	return tp.FirstName == "" && tp.LastName == "" && tp.Bio == ""
}

func (tp *TrainerProfile) IsValid() bool {
	return tp.FirstName != "" && tp.LastName != "" && tp.Bio != "" &&
		len(tp.Specializations) > 0 && tp.Experience >= 0 && tp.PricePerHour > 0
}

func (tp *TrainerProfile) GetFullName() string {
	return tp.FirstName + " " + tp.LastName
}
