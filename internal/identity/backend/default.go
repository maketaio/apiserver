package backend

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maketaio/apiserver/internal/types"
	"github.com/maketaio/apiserver/pkg/api"
)

type Storage interface {
	CreateUser(ctx context.Context, obj *api.User, hashedPassword string) error
}

type Default struct {
	storage Storage
}

func NewDefault(storage Storage) *Default {
	return &Default{
		storage: storage,
	}
}

func (d *Default) SignUp(ctx context.Context, input *types.SignUpInput) (*api.User, error) {
	// Hash password
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &api.User{
		ID:        uuid.New().String(),
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	if err := d.storage.CreateUser(ctx, user, hashedPassword); err != nil {
		return nil, fmt.Errorf("failed to insert user to db: %w", err)
	}

	return user, nil
}
