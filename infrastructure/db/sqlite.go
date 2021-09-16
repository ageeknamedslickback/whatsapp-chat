package db

import (
	"fmt"
	"time"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"gorm.io/gorm"
)

// DataStore defines our database layer interface
type DataStore interface {
	CreateMessage(message domain.Message)
	GetLastMessage(phone string) (*domain.Message, error)
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
	message.TimeStamp = time.Now()
	res := r.db.Create(&message)
	if res.Error != nil {
		return fmt.Errorf("failed to create a message: %w", res.Error)
	}
	return nil
}

// GetLastMessage queries the most recently sent inbound message
func (r *MessageRepository) GetLastMessage(phone string) (*domain.Message, error) {
	var message domain.Message
	res := r.db.Where(&domain.Message{From: phone}).Last(&message)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to query last message: %w", res.Error)
	}
	return &message, nil
}
