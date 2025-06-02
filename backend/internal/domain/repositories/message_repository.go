package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
	"fmt"
	"time"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.ChatMessage) error {
	query := `
		INSERT INTO messages (sender_id, receiver_id, content, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(query,
		message.SenderID,
		message.ReceiverID,
		message.Content,
		message.IsRead,
		time.Now()).Scan(&message.ID)
}

func (r *MessageRepository) GetConversation(userID, partnerID int, limit int, after int) ([]*models.ChatMessage, error) {
	var query string
	var args []interface{}

	baseCondition := "WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)"

	if after > 0 {
		baseCondition += " AND id > $3"
		args = []interface{}{userID, partnerID, after}
	} else {
		args = []interface{}{userID, partnerID}
	}

	if limit > 0 && after == 0 {
		// Pobierz ostatnie X wiadomości
		query = fmt.Sprintf(`
            SELECT * FROM (
                SELECT id, sender_id, receiver_id, content, is_read, created_at
                FROM messages
                %s
                ORDER BY created_at DESC
                LIMIT $%d
            ) sub ORDER BY created_at ASC
        `, baseCondition, len(args)+1)
		args = append(args, limit)
	} else {
		// Pobierz wszystkie lub nowe wiadomości
		query = fmt.Sprintf(`
            SELECT id, sender_id, receiver_id, content, is_read, created_at
            FROM messages
            %s
            ORDER BY created_at ASC
        `, baseCondition)

		if limit > 0 && after > 0 {
			query += fmt.Sprintf(" LIMIT $%d", len(args)+1)
			args = append(args, limit)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		message := &models.ChatMessage{}
		err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Content,
			&message.IsRead,
			&message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *MessageRepository) GetUnreadCount(userID int) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND is_read = false`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *MessageRepository) GetUnreadCountFromSender(userID, senderID int) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE receiver_id = $1 AND sender_id = $2 AND is_read = false`

	var count int
	err := r.db.QueryRow(query, userID, senderID).Scan(&count)
	return count, err
}

func (r *MessageRepository) MarkAsRead(messageID int) error {
	query := `UPDATE messages SET is_read = true WHERE id = $1`
	_, err := r.db.Exec(query, messageID)
	return err
}

func (r *MessageRepository) GetConversationPartners(userID int) ([]*models.User, error) {
	query := `
		SELECT DISTINCT u.id, u.username, u.role, u.is_online
		FROM users u
		INNER JOIN trainer_trainee tt ON 
			(u.id = tt.trainer_id AND tt.trainee_id = $1) OR 
			(u.id = tt.trainee_id AND tt.trainer_id = $1)
		WHERE u.id != $1
		ORDER BY u.username
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var partners []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.IsOnline)
		if err != nil {
			return nil, err
		}
		partners = append(partners, user)
	}

	return partners, nil
}
func (r *MessageRepository) MarkConversationAsRead(receiverID, senderID int) error {
	query := `UPDATE messages SET is_read = true WHERE receiver_id = $1 AND sender_id = $2 AND is_read = false`
	_, err := r.db.Exec(query, receiverID, senderID)
	return err
}
