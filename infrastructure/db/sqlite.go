package db

import (
	"fmt"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"gorm.io/gorm"
)

// DataStore defines our database layer interface
type DataStore interface {
	CreateMessage(message domain.Message)
}

// MessageRepository ..
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository ..
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// CreateMessage saves/creates a message in sqlite database
func (r *MessageRepository) CreateMessage(message domain.Message) error {
	res := r.db.Create(&message)
	if res.Error != nil {
		return fmt.Errorf("failed to create a message: %w", res.Error)
	}
	return nil
}
