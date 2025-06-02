package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
	"fmt"
	"log"
	"regexp"
)

type UserMediaRepository struct {
	DB *sql.DB
}

func NewUserMediaRepository(db *sql.DB) *UserMediaRepository {
	return &UserMediaRepository{DB: db}
}

func (r *UserMediaRepository) Create(media *models.UserMedia) error {
	query := `
		INSERT INTO user_media (uploader_id, recipient_id, file_name, file_url, file_type, description, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id, created_at
	`
	return r.DB.QueryRow(query, media.UploaderID, media.RecipientID, media.FileName, media.FileURL,
		media.FileType, media.Description, media.IsRead).Scan(&media.ID, &media.CreatedAt)
}

func (r *UserMediaRepository) GetSharedMedia(userID int) ([]models.UserMediaWithDetails, error) {
	query := `
		SELECT um.id, um.uploader_id, um.recipient_id, um.file_name, um.file_url, um.file_type,
		       um.description, um.is_read, um.created_at, u1.username, u2.username
		FROM user_media um
		JOIN users u1 ON um.uploader_id = u1.id
		JOIN users u2 ON um.recipient_id = u2.id
		WHERE um.recipient_id = $1
		ORDER BY um.created_at DESC
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []models.UserMediaWithDetails
	for rows.Next() {
		var m models.UserMediaWithDetails
		err := rows.Scan(&m.ID, &m.UploaderID, &m.RecipientID, &m.FileName, &m.FileURL, &m.FileType,
			&m.Description, &m.IsRead, &m.CreatedAt, &m.UploaderUsername, &m.RecipientUsername)
		if err != nil {
			return nil, err
		}
		media = append(media, m)
	}
	return media, nil
}

// Nowa metoda: Pobierz media przesłane przez trenera
func (r *UserMediaRepository) GetUploadedMedia(trainerID int) ([]models.UserMediaWithDetails, error) {
	query := `
		SELECT um.id, um.uploader_id, um.recipient_id, um.file_name, um.file_url, um.file_type,
		       um.description, um.is_read, um.created_at, u1.username, u2.username
		FROM user_media um
		JOIN users u1 ON um.uploader_id = u1.id
		JOIN users u2 ON um.recipient_id = u2.id
		WHERE um.uploader_id = $1
		ORDER BY um.created_at DESC
	`
	rows, err := r.DB.Query(query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var media []models.UserMediaWithDetails
	for rows.Next() {
		var m models.UserMediaWithDetails
		err := rows.Scan(&m.ID, &m.UploaderID, &m.RecipientID, &m.FileName, &m.FileURL, &m.FileType,
			&m.Description, &m.IsRead, &m.CreatedAt, &m.UploaderUsername, &m.RecipientUsername)
		if err != nil {
			return nil, err
		}
		media = append(media, m)
	}
	return media, nil
}

// Nowa metoda: Usuń media
func (r *UserMediaRepository) DeleteMedia(mediaID, trainerID int) error {
	query := `DELETE FROM user_media WHERE id = $1 AND uploader_id = $2`
	result, err := r.DB.Exec(query, mediaID, trainerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserMediaRepository) MarkAsRead(mediaID int) error {
	query := `UPDATE user_media SET is_read = true WHERE id = $1`
	_, err := r.DB.Exec(query, mediaID)
	return err
}

func (r *UserMediaRepository) GetUnreadCount(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM user_media WHERE recipient_id = $1 AND is_read = false`
	var count int
	err := r.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *UserMediaRepository) GetTrainerTrainees(trainerID int) ([]*models.User, error) {
	query := `
		SELECT DISTINCT u.id, u.username, u.role, u.is_online
		FROM users u
		INNER JOIN trainer_trainee tt ON 
			(u.id = tt.trainer_id AND tt.trainee_id = $1) OR 
			(u.id = tt.trainee_id AND tt.trainer_id = $1)
		WHERE u.id != $1
		ORDER BY u.username
	`

	rows, err := r.DB.Query(query, trainerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.IsOnline)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
func (r *UserMediaRepository) GetUserRole(userID int) (string, error) {
	var role string
	err := r.DB.QueryRow(`SELECT role FROM users WHERE id = $1`, userID).Scan(&role)
	return role, err
}
func (r *UserMediaRepository) GetExerciseVideos() ([]models.ExerciseVideo, error) {
	log.Println("🎬 Repository: GetExerciseVideos START")

	query := `
		SELECT id, name, youtube_url, created_at
		FROM exercise_videos
		ORDER BY name
	`

	log.Printf("🎬 Wykonuję query: %s", query)
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Printf("❌ Query error: %v", err)
		return nil, fmt.Errorf("failed to query exercise videos: %w", err)
	}
	defer rows.Close()

	var videos []models.ExerciseVideo
	for rows.Next() {
		var video models.ExerciseVideo
		err := rows.Scan(&video.ID, &video.Name, &video.YoutubeURL, &video.CreatedAt)
		if err != nil {
			log.Printf("❌ Scan error: %v", err)
			return nil, fmt.Errorf("failed to scan exercise video: %w", err)
		}
		log.Printf("🎬 Znaleziony film: %+v", video)
		videos = append(videos, video)
	}

	if err = rows.Err(); err != nil {
		log.Printf("❌ Rows error: %v", err)
		return nil, fmt.Errorf("error iterating exercise videos: %w", err)
	}

	log.Printf("🎬 Repository: Zwracam %d filmów", len(videos))
	return videos, nil
}

// Dodaj film z biblioteki do użytkownika (uproszczone)
func (r *UserMediaRepository) CreateFromExerciseVideo(uploaderID, recipientID, exerciseVideoID int) error {
	// Najpierw pobierz dane filmu z biblioteki
	var exerciseVideo models.ExerciseVideo
	query := `SELECT name, youtube_url FROM exercise_videos WHERE id = $1`

	err := r.DB.QueryRow(query, exerciseVideoID).Scan(
		&exerciseVideo.Name, &exerciseVideo.YoutubeURL)
	if err != nil {
		return err
	}

	// Konwertuj YouTube URL do embed format
	youtubeID := extractYouTubeID(exerciseVideo.YoutubeURL)
	if youtubeID == "" {
		return fmt.Errorf("invalid YouTube URL")
	}
	embedURL := "https://www.youtube.com/embed/" + youtubeID

	// Dodaj do user_media
	insertQuery := `
		INSERT INTO user_media (uploader_id, recipient_id, file_name, file_url, file_type, description, is_read, created_at)
		VALUES ($1, $2, $3, $4, 'video', '', false, NOW())
	`
	_, err = r.DB.Exec(insertQuery, uploaderID, recipientID, exerciseVideo.Name, embedURL)

	return err
}

// Funkcja pomocnicza (dodaj do pliku z repository)
func extractYouTubeID(url string) string {
	patterns := []string{
		`(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)`,
		`youtube\.com\/v\/([^&\n?#]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}
