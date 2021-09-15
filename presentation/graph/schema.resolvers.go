package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	"github.com/ageeknamedslickback/whatsapp-chat/presentation/graph/generated"
)

func (r *mutationResolver) SendMessage(ctx context.Context, to string, body string) (*domain.Message, error) {
	return r.service.SendMessage(ctx, to, body)
}

func (r *subscriptionResolver) Senders(ctx context.Context) (<-chan []*domain.Sender, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) Messages(ctx context.Context, phone string) (<-chan []*domain.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
