package app

import (
	"authgrpc/internal/app/grpcapp"
	"authgrpc/internal/config"
	"authgrpc/internal/services/auth"
	"authgrpc/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	db config.DBConfig,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.New(db)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
