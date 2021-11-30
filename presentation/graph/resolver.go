package graph

import (
	"fmt"
	"sync"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"github.com/ageeknamedslickback/whatsapp-chat/usecases"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	service usecases.MessagesService

	// All messages since launching the GraphQL endpoint
	ChatMessages []*domain.Message
	// All active subscriptions
	Observers map[string]chan []*domain.Message
	// All message senders
	ChatSenders []*domain.Sender
	// All senders active subscribtions
	SenderObservers map[string]chan []*domain.Sender
}

func NewResolver(s usecases.MessagesService) *Resolver {
	return &Resolver{
		service:         s,
		Observers:       make(map[string]chan []*domain.Message),
		SenderObservers: make(map[string]chan []*domain.Sender),
	}
}

func (res *Resolver) AddMessagesToChannel(phone, profileName string) error {
	// get all the messages for a person
	messages, err := res.service.GetMessages(phone)
	if err != nil {
		return fmt.Errorf("failed to get %s messages: %v", phone, err)
	}
	sender := domain.Sender{
		PhoneNumber:         phone,
		ProfileName:         profileName,
		UnreadMessagesCount: 1, // todo: get unread messages count and read them
		Messages:            messages,
	}

	res.ChatSenders = append(res.ChatSenders, &sender)
	var mu sync.Mutex
	mu.Lock()
	for _, observer := range res.SenderObservers {
		observer <- res.ChatSenders
	}
	mu.Unlock()

	return nil
}
