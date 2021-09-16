package graph

import "github.com/ageeknamedslickback/whatsapp-chat/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
//go:generate go run github.com/99designs/gqlgen

type Resolver struct {
	service usecases.MessagesService
}

func NewResolver(s usecases.MessagesService) *Resolver {
	return &Resolver{service: s}
}
