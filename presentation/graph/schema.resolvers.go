package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"sync"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"github.com/ageeknamedslickback/whatsapp-chat/presentation/graph/generated"
	"github.com/segmentio/ksuid"
)

var mu sync.Mutex

func (r *mutationResolver) SendMessage(ctx context.Context, to string, body string) (*domain.Message, error) {
	message, err := r.service.SendMessage(ctx, to, body)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}
	r.ChatMessages = append(r.ChatMessages, message)
	mu.Lock()
	// Notify all active subscriptions that a new message has been posted by posted. In this case we push the now
	// updated ChatMessages array to all clients that care about it.
	for _, observer := range r.Observers {
		observer <- r.ChatMessages
	}
	mu.Unlock()

	return message, nil
}

func (r *subscriptionResolver) Senders(ctx context.Context) (<-chan []*domain.Sender, error) {
	id := ksuid.New().String()
	senders := make(chan []*domain.Sender, 1)

	go func() {
		<-ctx.Done()
		mu.Lock()
		delete(r.SenderObservers, id)
		mu.Unlock()
	}()

	mu.Lock()
	// Keep a reference of the channel so that we can push changes into it when new messages are posted.
	r.SenderObservers[id] = senders
	mu.Unlock()

	r.SenderObservers[id] <- r.ChatSenders

	return senders, nil
}

func (r *subscriptionResolver) Messages(ctx context.Context, phone string) (<-chan []*domain.Message, error) {
	// Create an ID and channel for each active subscriptions. Changes will be pushed into this channel
	// When a new subscription is created, this resolver will fire first
	id := ksuid.New().String()
	msgs := make(chan []*domain.Message, 1)

	// Start a goroutine to allow for cleaning up of disconnected subscriptions
	// This go routine will only get past Done() when a client terminates the subscription. This allows us
	// to only then remove the reference from the list of ChatObservers since it is no longer needed
	go func() {
		<-ctx.Done()
		mu.Lock()
		delete(r.Observers, id)
		mu.Unlock()
	}()

	mu.Lock()
	// Keep a reference of the channel so that we can push changes into it when new messages are posted.
	r.Observers[id] = msgs
	mu.Unlock()

	// This is optional, and this allows newly subscribed clients to get a list of all the messages that have been
	// posted so far. Upon subscribing the client will be pushed the messages once, further changes are handled
	// in the PostMessage mutation.
	// todo: for consistency, get these from a persistent storage(db)
	r.Observers[id] <- r.ChatMessages

	return msgs, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
