package ports

import "github.com/Paulooo0/modak-challenge/internal/domain/entity"

type NotificationGateway interface {
	Send(n entity.Notification) error
}
