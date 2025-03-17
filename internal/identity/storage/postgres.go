package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maketaio/apiserver/pkg/api"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: pool,
	}
}

func (s *Postgres) CreateUser(ctx context.Context, user *api.User, hashedPassword string) error {
	_, err := s.pool.Exec(
		ctx,
		`INSERT INTO users (id, first_name, last_name, email, hashed_password) VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		hashedPassword,
	)

	if err != nil {
		return err
	}

	return nil
}
