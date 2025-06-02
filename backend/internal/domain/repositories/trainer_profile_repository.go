package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
	"encoding/json"
	"time"
)

type TrainerProfileRepository struct {
	DB *sql.DB
}

func NewTrainerProfileRepository(db *sql.DB) *TrainerProfileRepository {
	return &TrainerProfileRepository{DB: db}
}

func (r *TrainerProfileRepository) GetByUserID(userID int) (*models.TrainerProfile, error) {
	row := r.DB.QueryRow(`
		SELECT user_id, first_name, last_name, bio, specializations, experience, 
		       certifications, price_per_hour, languages, contact_email, contact_phone, 
		       location, avatar, is_active, created_at, updated_at 
		FROM trainer_profiles 
		WHERE user_id = $1`, userID)

	var tp models.TrainerProfile
	var specializationsJSON, certificationsJSON, languagesJSON string
	var avatar sql.NullString

	err := row.Scan(&tp.UserID, &tp.FirstName, &tp.LastName, &tp.Bio,
		&specializationsJSON, &tp.Experience, &certificationsJSON,
		&tp.PricePerHour, &languagesJSON, &tp.ContactEmail, &tp.ContactPhone,
		&tp.Location, &avatar, &tp.IsActive, &tp.CreatedAt, &tp.UpdatedAt)

	if err != nil {
		return nil, err
	}

	// Parse JSON arrays
	json.Unmarshal([]byte(specializationsJSON), &tp.Specializations)
	json.Unmarshal([]byte(certificationsJSON), &tp.Certifications)
	json.Unmarshal([]byte(languagesJSON), &tp.Languages)

	if avatar.Valid {
		tp.Avatar = avatar.String
	}

	return &tp, nil
}

func (r *TrainerProfileRepository) CreateOrUpdate(tp *models.TrainerProfile) error {
	now := time.Now()

	// Convert slices to JSON
	specializationsJSON, _ := json.Marshal(tp.Specializations)
	certificationsJSON, _ := json.Marshal(tp.Certifications)
	languagesJSON, _ := json.Marshal(tp.Languages)

	_, err := r.DB.Exec(`
		INSERT INTO trainer_profiles (
			user_id, first_name, last_name, bio, specializations, experience, 
			certifications, price_per_hour, languages, contact_email, contact_phone, 
			location, avatar, is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		ON CONFLICT (user_id)
		DO UPDATE SET 
			first_name = $2, 
			last_name = $3, 
			bio = $4, 
			specializations = $5, 
			experience = $6, 
			certifications = $7, 
			price_per_hour = $8, 
			languages = $9, 
			contact_email = $10, 
			contact_phone = $11, 
			location = $12, 
			avatar = $13, 
			is_active = $14, 
			updated_at = $16
	`, tp.UserID, tp.FirstName, tp.LastName, tp.Bio, string(specializationsJSON),
		tp.Experience, string(certificationsJSON), tp.PricePerHour, string(languagesJSON),
		tp.ContactEmail, tp.ContactPhone, tp.Location, tp.Avatar, tp.IsActive, now, now)

	return err
}

func (r *TrainerProfileRepository) HasValidProfile(userID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM trainer_profiles 
		WHERE user_id = $1 
		AND first_name != '' 
		AND last_name != '' 
		AND bio != '' 
		AND price_per_hour > 0
	`, userID).Scan(&count)

	return count > 0, err
}

func (r *TrainerProfileRepository) GetAllActiveTrainers() ([]*models.TrainerProfile, error) {
	rows, err := r.DB.Query(`
		SELECT user_id, first_name, last_name, bio, specializations, experience, 
		       certifications, price_per_hour, languages, contact_email, contact_phone, 
		       location, avatar, is_active, created_at, updated_at 
		FROM trainer_profiles 
		WHERE is_active = true
		ORDER BY experience DESC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trainers []*models.TrainerProfile
	for rows.Next() {
		var tp models.TrainerProfile
		var specializationsJSON, certificationsJSON, languagesJSON string
		var avatar sql.NullString

		err := rows.Scan(&tp.UserID, &tp.FirstName, &tp.LastName, &tp.Bio,
			&specializationsJSON, &tp.Experience, &certificationsJSON,
			&tp.PricePerHour, &languagesJSON, &tp.ContactEmail, &tp.ContactPhone,
			&tp.Location, &avatar, &tp.IsActive, &tp.CreatedAt, &tp.UpdatedAt)

		if err != nil {
			continue
		}

		// Parse JSON arrays
		json.Unmarshal([]byte(specializationsJSON), &tp.Specializations)
		json.Unmarshal([]byte(certificationsJSON), &tp.Certifications)
		json.Unmarshal([]byte(languagesJSON), &tp.Languages)

		if avatar.Valid {
			tp.Avatar = avatar.String
		}

		trainers = append(trainers, &tp)
	}

	return trainers, nil
}
