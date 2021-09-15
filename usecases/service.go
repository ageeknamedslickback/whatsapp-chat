package usecases

import (
	"fmt"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"github.com/ageeknamedslickback/whatsapp-chat/infrastructure/db"
)

// MessageUsecases ..
type MessageUsecases interface {
	InboundMessages(
		message domain.Message,
	) (*domain.Message, error)
}

// MessagesService ..
type MessagesService struct {
	datastore db.MessageRepository
}

// NewMessageService ..
func NewMessageService(d db.MessageRepository) *MessagesService {
	return &MessagesService{datastore: d}
}

// InboundMessages recieves and saves user-initiated messages
func (s *MessagesService) InboundMessages(
	message domain.Message,
) (*domain.Message, error) {
	if err := s.datastore.CreateMessage(message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return &message, nil
}
