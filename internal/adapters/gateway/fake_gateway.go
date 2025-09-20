package gateway

import (
	"fmt"

	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/ports"
)

type FakeGateway struct{}

func NewFakeGateway() ports.NotificationGateway {
	return &FakeGateway{}
}

func (g *FakeGateway) Send(n entity.Notification) error {
	fmt.Printf("📩 sending %s notification to %s: %s\n", n.Type, n.UserID, n.Message)
	return nil
}
