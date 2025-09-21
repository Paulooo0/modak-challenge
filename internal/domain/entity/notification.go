package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      NotificationType
	Message   string
	CreatedAt time.Time
}

type NotificationType string

const (
	Status    NotificationType = "status"
	News      NotificationType = "news"
	Marketing NotificationType = "marketing"
)

func IsValidNotificationType(s NotificationType) bool {
	switch s {
	case Status, News, Marketing:
		return true
	default:
		return false
	}
}
