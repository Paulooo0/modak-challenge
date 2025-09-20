package ports

import (
	"context"
	"time"

	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
)

type NotificationRepository interface {
	Create(ctx context.Context, n entity.Notification) (entity.Notification, error)
	CountInTimeWindow(ctx context.Context, userID, notifType string, window time.Time) (int, error)
}
