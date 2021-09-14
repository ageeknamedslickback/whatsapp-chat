package db

import "gorm.io/gorm"

// DataStore defines our database layer interface
type DataStore interface{}

// MessageRepository ..
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository ..
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}
