package db

import (
	"context"
	"time"

	"github.com/Paulooo0/modak-challenge/internal/adapters/db/sqlc"
	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/ports"
	"github.com/google/uuid"
)

// notificationsQuerier is a minimal interface implemented by *sqlc.Queries
// that allows the repository to be tested with simple mocks.
type notificationsQuerier interface {
	CreateNotification(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.Notification, error)
	CountNotificationsInTimeWindow(ctx context.Context, arg sqlc.CountNotificationsInTimeWindowParams) (int64, error)
}

type NotificationRepository struct {
	q notificationsQuerier
}

func NewNotificationRepository(q notificationsQuerier) ports.NotificationRepository {
	return &NotificationRepository{q: q}
}

func (r *NotificationRepository) Create(ctx context.Context, n entity.Notification) (entity.Notification, error) {
	row, err := r.q.CreateNotification(ctx, sqlc.CreateNotificationParams{
		UserID:  n.UserID,
		Type:    string(n.Type),
		Message: n.Message,
	})
	if err != nil {
		return entity.Notification{}, err
	}

	return entity.Notification{
		ID:        row.ID,
		UserID:    row.UserID,
		Type:      entity.NotificationType(row.Type),
		Message:   row.Message,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *NotificationRepository) CountInTimeWindow(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, since time.Time) (int, error) {
	count, err := r.q.CountNotificationsInTimeWindow(ctx, sqlc.CountNotificationsInTimeWindowParams{
		UserID:    userID,
		Type:      string(notifType),
		CreatedAt: since,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
