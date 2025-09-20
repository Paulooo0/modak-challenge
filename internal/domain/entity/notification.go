package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID
	UserID    string
	Type      string
	Message   string
	CreatedAt time.Time
}
