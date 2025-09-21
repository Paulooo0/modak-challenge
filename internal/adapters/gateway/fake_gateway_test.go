package gateway

import (
	"testing"

	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFakeGatewaySendNoError(t *testing.T) {
	g := NewFakeGateway()
	n := entity.Notification{ID: uuid.New(), UserID: uuid.New(), Type: entity.Status, Message: "hello"}
	require.NoError(t, g.Send(n))
}
