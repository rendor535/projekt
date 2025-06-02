package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
	"time"
)

type ProfileRepository struct {
	DB *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{DB: db}
}

func (r *ProfileRepository) GetByUserID(userID int) (*models.Profile, error) {
	row := r.DB.QueryRow(`
		SELECT user_id, age, height, weight, gender, activity_level, goal, created_at, updated_at 
		FROM profiles 
		WHERE user_id = $1`, userID)

	var p models.Profile
	var createdAt, updatedAt time.Time

	err := row.Scan(&p.UserID, &p.Age, &p.Height, &p.Weight, &p.Gender, &p.ActivityLevel, &p.Goal, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	p.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	p.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	return &p, nil
}

func (r *ProfileRepository) CreateOrUpdate(p *models.Profile) error {
	now := time.Now()

	_, err := r.DB.Exec(`
		INSERT INTO profiles (user_id, age, height, weight, gender, activity_level, goal, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id)
		DO UPDATE SET 
			age = $2, 
			height = $3, 
			weight = $4, 
			gender = $5, 
			activity_level = $6, 
			goal = $7, 
			updated_at = $9
	`, p.UserID, p.Age, p.Height, p.Weight, p.Gender, p.ActivityLevel, p.Goal, now, now)

	return err
}

func (r *ProfileRepository) HasValidProfile(userID int) (bool, error) {
	var count int
	err := r.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM profiles 
		WHERE user_id = $1 
		AND age > 0 
		AND height > 0 
		AND weight > 0 
		AND gender != '' 
		AND activity_level > 0 
		AND goal != ''
	`, userID).Scan(&count)

	return count > 0, err
}
