package usecases

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"github.com/ageeknamedslickback/whatsapp-chat/infrastructure/db"
	"github.com/ageeknamedslickback/whatsapp-chat/infrastructure/services"
)

// MessageUsecases ..
type MessageUsecases interface {
	InboundMessages(
		message domain.Message,
	) (*domain.Message, error)
	SendMessage(
		ctx context.Context,
		to, body string,
	) (*domain.Message, error)
}

// MessagesService ..
type MessagesService struct {
	datastore db.MessageRepository
}

// NewMessageService ..
func NewMessageService(d db.MessageRepository) *MessagesService {
	return &MessagesService{
		datastore: d,
	}
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

// SendMessage sends free-form messages to user-initiated messages
func (s *MessagesService) SendMessage(
	ctx context.Context,
	to, body string,
) (*domain.Message, error) {
	phone := fmt.Sprintf("whatsapp:%s", to)
	msg, err := s.datastore.GetLastMessage(phone)
	if err != nil {
		return nil, fmt.Errorf("failed to get last message sent: %w", err)
	}
	if time.Since(msg.TimeStamp).Hours() > 24 {
		return nil, fmt.Errorf("24 hour session exceeded")
	}

	params := url.Values{}
	params.Add("From", fmt.Sprintf("whatsapp:%s", os.Getenv("TWILIO_NUMBER")))
	params.Add("To", phone)
	params.Add("Body", body)

	var message domain.Message
	if err := services.MakePostRequest(params, &message); err != nil {
		return nil, fmt.Errorf("failed to make twilio's post request: %w", err)
	}

	if err := s.datastore.CreateMessage(message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return &message, nil
}
