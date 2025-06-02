package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
)

type ProgressReportRepository struct {
	db *sql.DB
}

func NewProgressReportRepository(db *sql.DB) *ProgressReportRepository {
	return &ProgressReportRepository{db: db}
}

func (r *ProgressReportRepository) Create(report *models.ProgressReport) error {
	query := `
		INSERT INTO progress_reports (client_id, trainer_id, weight, body_fat, measurements, notes, photo_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id, created_at
	`

	// Handle NULL trainer_id for self-reports
	var trainerID interface{}
	if report.TrainerID != nil {
		trainerID = *report.TrainerID
	} else {
		trainerID = nil
	}

	return r.db.QueryRow(query, report.ClientID, trainerID, report.Weight, report.BodyFat,
		report.Measurements, report.Notes, report.PhotoURL).Scan(&report.ID, &report.CreatedAt)
}

func (r *ProgressReportRepository) GetByClientID(clientID int) ([]models.ProgressReportWithTrainer, error) {
	query := `
		SELECT pr.id, pr.client_id, pr.trainer_id, pr.weight, pr.body_fat, pr.measurements,
		       pr.notes, pr.photo_url, pr.created_at, COALESCE(u.username, 'Self-report') as username
		FROM progress_reports pr
		LEFT JOIN users u ON pr.trainer_id = u.id
		WHERE pr.client_id = $1
		ORDER BY pr.created_at DESC
	`
	rows, err := r.db.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []models.ProgressReportWithTrainer
	for rows.Next() {
		var report models.ProgressReportWithTrainer
		err := rows.Scan(&report.ID, &report.ClientID, &report.TrainerID, &report.Weight, &report.BodyFat,
			&report.Measurements, &report.Notes, &report.PhotoURL, &report.CreatedAt, &report.TrainerUsername)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

// Nowa metoda: Statystyki użytkownika
// Update your GetUserStats method in ProgressReportRepository to handle empty data
func (r *ProgressReportRepository) GetUserStats(userID int) (*models.UserStats, error) {
	stats := &models.UserStats{}

	// Check if user has any reports first
	var reportCount int
	countQuery := `SELECT COUNT(*) FROM progress_reports WHERE client_id = $1`
	err := r.db.QueryRow(countQuery, userID).Scan(&reportCount)
	if err != nil {
		return nil, err
	}

	stats.TotalReports = reportCount

	// If no reports, return basic stats
	if reportCount == 0 {
		return stats, nil
	}

	// Basic statistics for users with reports
	query := `
		SELECT 
			COALESCE(MIN(weight), 0) as min_weight,
			COALESCE(MAX(weight), 0) as max_weight,
			COALESCE(AVG(weight), 0) as avg_weight,
			COALESCE(MIN(body_fat), 0) as min_body_fat,
			COALESCE(MAX(body_fat), 0) as max_body_fat,
			COALESCE(AVG(body_fat), 0) as avg_body_fat
		FROM progress_reports 
		WHERE client_id = $1 AND weight IS NOT NULL
	`

	var minBodyFat, maxBodyFat, avgBodyFat sql.NullFloat64
	err = r.db.QueryRow(query, userID).Scan(
		&stats.MinWeight, &stats.MaxWeight, &stats.AvgWeight,
		&minBodyFat, &maxBodyFat, &avgBodyFat)
	if err != nil {
		return nil, err
	}

	if minBodyFat.Valid {
		stats.MinBodyFat = minBodyFat.Float64
	}
	if maxBodyFat.Valid {
		stats.MaxBodyFat = maxBodyFat.Float64
	}
	if avgBodyFat.Valid {
		stats.AvgBodyFat = avgBodyFat.Float64
	}

	// Latest weight and body fat
	latestQuery := `
		SELECT COALESCE(weight, 0), body_fat, created_at 
		FROM progress_reports 
		WHERE client_id = $1 AND weight IS NOT NULL 
		ORDER BY created_at DESC LIMIT 1
	`
	var latestBodyFat sql.NullFloat64
	err = r.db.QueryRow(latestQuery, userID).Scan(&stats.LatestWeight, &latestBodyFat, &stats.LatestDate)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if latestBodyFat.Valid {
		stats.LatestBodyFat = latestBodyFat.Float64
	}

	// First weight
	firstQuery := `
		SELECT COALESCE(weight, 0), created_at 
		FROM progress_reports 
		WHERE client_id = $1 AND weight IS NOT NULL 
		ORDER BY created_at ASC LIMIT 1
	`
	err = r.db.QueryRow(firstQuery, userID).Scan(&stats.FirstWeight, &stats.FirstDate)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Calculate weight change
	if stats.FirstWeight > 0 && stats.LatestWeight > 0 {
		stats.WeightChange = stats.LatestWeight - stats.FirstWeight
	}

	return stats, nil
}

// Nowa metoda: Dane do wykresów
func (r *ProgressReportRepository) GetChartData(userID int) (*models.ChartData, error) {
	query := `
		SELECT weight, body_fat, created_at
		FROM progress_reports
		WHERE client_id = $1 AND weight IS NOT NULL
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chartData := &models.ChartData{
		WeightData:  []models.ChartPoint{},
		BodyFatData: []models.ChartPoint{},
	}

	for rows.Next() {
		var weight float64
		var bodyFat sql.NullFloat64
		var date string

		err := rows.Scan(&weight, &bodyFat, &date)
		if err != nil {
			return nil, err
		}

		chartData.WeightData = append(chartData.WeightData, models.ChartPoint{
			Date:  date,
			Value: weight,
		})

		if bodyFat.Valid {
			chartData.BodyFatData = append(chartData.BodyFatData, models.ChartPoint{
				Date:  date,
				Value: bodyFat.Float64,
			})
		}
	}

	return chartData, nil
}

// Nowa metoda: Statystyki trenera
func (r *ProgressReportRepository) GetTrainerStats(trainerID int) (*models.TrainerStats, error) {
	stats := &models.TrainerStats{
		TraineeStats: []models.TraineeStatsItem{},
	}

	// Only count reports where trainer_id is NOT NULL and equals the trainer
	overviewQuery := `
		SELECT 
			COUNT(DISTINCT tt.trainee_id) as total_trainees,
			COUNT(pr.id) as total_reports,
			AVG(pr.weight) as avg_weight_all_trainees
		FROM trainer_trainee tt
		LEFT JOIN progress_reports pr ON tt.trainee_id = pr.client_id AND pr.trainer_id = $1
		WHERE tt.trainer_id = $1
	`

	var avgWeight sql.NullFloat64
	err := r.db.QueryRow(overviewQuery, trainerID).Scan(&stats.TotalTrainees, &stats.TotalReports, &avgWeight)
	if err != nil {
		return nil, err
	}
	if avgWeight.Valid {
		stats.AvgWeight = avgWeight.Float64
	}

	// Get stats for each trainee (including self-reports AND trainer reports)
	traineeQuery := `
		SELECT 
			u.id, u.username,
			COUNT(pr.id) as report_count,
			MAX(pr.weight) as latest_weight,
			MIN(pr.weight) as first_weight,
			MAX(pr.created_at) as last_report_date
		FROM trainer_trainee tt
		JOIN users u ON tt.trainee_id = u.id
		LEFT JOIN progress_reports pr ON tt.trainee_id = pr.client_id
		WHERE tt.trainer_id = $1
		GROUP BY u.id, u.username
		ORDER BY last_report_date DESC NULLS LAST
	`

	rows, err := r.db.Query(traineeQuery, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.TraineeStatsItem
		var latestWeight, firstWeight sql.NullFloat64
		var lastReportDate sql.NullString

		err := rows.Scan(&item.TraineeID, &item.TraineeName, &item.ReportCount,
			&latestWeight, &firstWeight, &lastReportDate)
		if err != nil {
			return nil, err
		}

		if latestWeight.Valid {
			item.LatestWeight = latestWeight.Float64
		}
		if firstWeight.Valid {
			item.FirstWeight = firstWeight.Float64
			if latestWeight.Valid {
				item.WeightChange = item.LatestWeight - item.FirstWeight
			}
		}
		if lastReportDate.Valid {
			item.LastReportDate = lastReportDate.String
		}

		stats.TraineeStats = append(stats.TraineeStats, item)
	}

	return stats, nil
}

// Nowa metoda: Raporty konkretnego podopiecznego dla trenera
func (r *ProgressReportRepository) GetTraineeReports(trainerID, traineeID int) ([]models.ProgressReportWithTrainer, error) {
	// Sprawdź czy trener ma dostęp do tego podopiecznego
	checkQuery := `SELECT 1 FROM trainer_trainee WHERE trainer_id = $1 AND trainee_id = $2`
	var exists int
	err := r.db.QueryRow(checkQuery, trainerID, traineeID).Scan(&exists)
	if err != nil {
		return nil, err
	}

	// Pobierz raporty
	return r.GetByClientID(traineeID)
}

func (r *ProgressReportRepository) GetLatestWeight(clientID int) (float64, error) {
	query := `SELECT weight FROM progress_reports WHERE client_id = $1 AND weight IS NOT NULL ORDER BY created_at DESC LIMIT 1`
	var weight float64
	err := r.db.QueryRow(query, clientID).Scan(&weight)
	return weight, err
}
func (r *ProgressReportRepository) VerifyTrainerAccess(trainerID, traineeID int) (bool, error) {
	query := `SELECT 1 FROM trainer_trainee WHERE trainer_id = $1 AND trainee_id = $2`
	var exists int
	err := r.db.QueryRow(query, trainerID, traineeID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}
