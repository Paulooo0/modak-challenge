package ports

import (
	"context"
	"time"

	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/google/uuid"
)

type NotificationRepository interface {
	Create(ctx context.Context, n entity.Notification) (entity.Notification, error)
	CountInTimeWindow(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, window time.Time) (int, error)
}
