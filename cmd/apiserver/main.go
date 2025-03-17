package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	identitybackend "github.com/maketaio/apiserver/internal/identity/backend"
	identityrouter "github.com/maketaio/apiserver/internal/identity/router"
	identitystorage "github.com/maketaio/apiserver/internal/identity/storage"
	"github.com/maketaio/apiserver/pkg/httpserver"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// TODO: use a CLI arg or flag instead to pass database URL
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	identityStorage := identitystorage.NewPostgres(pool)
	identityBackend := identitybackend.NewDefault(identityStorage)
	identityRouter := identityrouter.New(identityBackend)

	s := httpserver.New(&httpserver.Options{
		Addr:             ":8080",
		MaxBodySize:      "1MB",
		CompressionLevel: 5,
		Logger:           logger,
	})

	s.AddRouter("/identity", identityRouter)

	if err := s.Start(); err != nil {
		logger.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}
