package db

import (
	"context"
	"time"

	"github.com/Paulooo0/modak-challenge/internal/adapters/db/sqlc"
	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/ports"
	"github.com/google/uuid"
)

type NotificationRepository struct {
	q *sqlc.Queries
}

func NewNotificationRepository(q *sqlc.Queries) ports.NotificationRepository {
	return &NotificationRepository{q: q}
}

func (r *NotificationRepository) Create(ctx context.Context, n entity.Notification) (entity.Notification, error) {
	row, err := r.q.CreateNotification(ctx, sqlc.CreateNotificationParams{
		UserID:  n.UserID,
		Type:    n.Type,
		Message: n.Message,
	})
	if err != nil {
		return entity.Notification{}, err
	}

	return entity.Notification{
		ID:        row.ID,
		UserID:    row.UserID,
		Type:      row.Type,
		Message:   row.Message,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *NotificationRepository) CountInTimeWindow(ctx context.Context, userID uuid.UUID, notifType string, since time.Time) (int, error) {
	count, err := r.q.CountNotificationsInTimeWindow(ctx, sqlc.CountNotificationsInTimeWindowParams{
		UserID:    userID,
		Type:      notifType,
		CreatedAt: since,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
