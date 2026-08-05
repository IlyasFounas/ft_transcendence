package handlers

import (
	// OUR PACKAGE
	"backend/internal/repositories"
)

type Handlers struct {
	Location *LocationHandler
	User     *UserHandler
	Token    *TokenHandler
}

// NewHandlers instanciates and assembles all the handlers
func NewHandlers(repo *repositories.Repositories) *Handlers {
	tokenHandler := NewTokenHandler(repo.User, repo.Token)
	return &Handlers{
		Location: NewLocationHandler(repo.Location),
		User:     NewUserHandler(repo.User, repo.Token, tokenHandler),
		Token:    tokenHandler,
	}
}
