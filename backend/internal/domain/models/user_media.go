package models

import "time"

type UserMedia struct {
	ID          int       `json:"id"`
	UploaderID  int       `json:"uploader_id"`
	RecipientID int       `json:"recipient_id"`
	FileName    string    `json:"file_name"`
	FileURL     string    `json:"file_url"`
	FileType    string    `json:"file_type"` // video, image, document
	Description string    `json:"description"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserMediaWithDetails struct {
	UserMedia
	UploaderUsername  string `json:"uploader_username"`
	RecipientUsername string `json:"recipient_username"`
}
type ExerciseVideo struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	YoutubeURL string    `json:"youtube_url"`
	CreatedAt  time.Time `json:"created_at"`
}
