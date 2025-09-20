package main

import (
	"context"
	"log"

	"github.com/Paulooo0/modak-challenge/internal/adapters/db"
	"github.com/Paulooo0/modak-challenge/internal/adapters/db/sqlc"
	"github.com/Paulooo0/modak-challenge/internal/adapters/gateway"
	"github.com/Paulooo0/modak-challenge/internal/adapters/http"
	"github.com/Paulooo0/modak-challenge/internal/config"
	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/domain/useCase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	q := sqlc.New(pool)
	repo := db.NewNotificationRepository(q)
	gateway := gateway.NewFakeGateway()

	uc := useCase.NewNotificationUseCase(repo, gateway, entity.DefaultRateLimits)

	r := http.NewRouter(uc)

	log.Println("Server running on :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
