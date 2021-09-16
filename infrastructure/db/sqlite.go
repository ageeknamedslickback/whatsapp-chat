package db

import (
	"fmt"
	"sort"
	"time"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"gorm.io/gorm"
)

// DataStore defines our database layer interface
type DataStore interface {
	CreateMessage(message domain.Message)
	GetLastMessage(phone string) (*domain.Message, error)
	GetMessages(phone string) ([]domain.Message, error)
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

func (r *MessageRepository) GetSentMessages(phone string) ([]domain.Message, error) {
	var messages []domain.Message
	sent := r.db.Where(&domain.Message{To: phone}).Find(&messages)
	if sent.Error != nil {
		return nil, fmt.Errorf("failed to query sent messages: %w", sent.Error)
	}

	return messages, nil
}

func (r *MessageRepository) GetReceivedMessages(phone string) ([]domain.Message, error) {
	var messages []domain.Message
	res := r.db.Where(&domain.Message{From: phone}).Find(&messages)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to query received messages: %w", res.Error)
	}
	return messages, nil
}

func (r *MessageRepository) GetMessages(phone string) ([]domain.Message, error) {
	var messages []domain.Message
	sent, err := r.GetSentMessages(phone)
	if err != nil {
		return messages, fmt.Errorf("failed to get sent messages: %w", err)
	}
	messages = append(messages, sent...)

	received, err := r.GetReceivedMessages(phone)
	if err != nil {
		return messages, fmt.Errorf("failed to get sent messages: %w", err)
	}
	messages = append(messages, received...)

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].TimeStamp.Before(messages[j].TimeStamp)
	})

	return messages, nil
}
