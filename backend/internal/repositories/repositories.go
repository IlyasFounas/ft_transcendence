package repositories

import (
	"github.com/uptrace/bun"
)

type Repositories struct {
	Location *LocationRepository
	User     *UserRepository
	Token    *TokenRepository
}

// NewRepositories instanciates and assembles all the repositories
func NewRepositories(db *bun.DB) (*Repositories, error) {
	locationRepo, err := NewLocationRepository(db)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Location: locationRepo,
		User:     NewUserRepository(db),
		Token:    NewTokenRepository(db),
	}, nil
}
