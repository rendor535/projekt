package models

import "time"

type ProgressReport struct {
	ID           int       `json:"id"`
	ClientID     int       `json:"client_id"`
	TrainerID    *int      `json:"trainer_id"` // Może być null
	Weight       *float64  `json:"weight"`
	BodyFat      *float64  `json:"body_fat"`
	Measurements string    `json:"measurements"` // JSON string
	Notes        string    `json:"notes"`
	PhotoURL     string    `json:"photo_url"`
	CreatedAt    time.Time `json:"created_at"`
}

type ProgressReportWithTrainer struct {
	ProgressReport
	TrainerUsername string `json:"trainer_username"`
}

// Nowe modele dla statystyk
type UserStats struct {
	TotalReports  int       `json:"total_reports"`
	MinWeight     float64   `json:"min_weight"`
	MaxWeight     float64   `json:"max_weight"`
	AvgWeight     float64   `json:"avg_weight"`
	LatestWeight  float64   `json:"latest_weight"`
	FirstWeight   float64   `json:"first_weight"`
	WeightChange  float64   `json:"weight_change"`
	MinBodyFat    float64   `json:"min_body_fat"`
	MaxBodyFat    float64   `json:"max_body_fat"`
	AvgBodyFat    float64   `json:"avg_body_fat"`
	LatestBodyFat float64   `json:"latest_body_fat"`
	LatestDate    time.Time `json:"latest_date"`
	FirstDate     time.Time `json:"first_date"`
}

type ChartPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type ChartData struct {
	WeightData  []ChartPoint `json:"weight_data"`
	BodyFatData []ChartPoint `json:"body_fat_data"`
}

type TraineeStatsItem struct {
	TraineeID      int     `json:"trainee_id"`
	TraineeName    string  `json:"trainee_name"`
	ReportCount    int     `json:"report_count"`
	LatestWeight   float64 `json:"latest_weight"`
	FirstWeight    float64 `json:"first_weight"`
	WeightChange   float64 `json:"weight_change"`
	LastReportDate string  `json:"last_report_date"`
}

type TrainerStats struct {
	TotalTrainees int                `json:"total_trainees"`
	TotalReports  int                `json:"total_reports"`
	AvgWeight     float64            `json:"avg_weight"`
	TraineeStats  []TraineeStatsItem `json:"trainee_stats"`
}
