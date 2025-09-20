package db

import (
	"context"
	"testing"
	"time"

	"github.com/Paulooo0/modak-challenge/internal/adapters/db/sqlc"
	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockQueries struct{ mock.Mock }

func (m *mockQueries) CreateNotification(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.Notification, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.Notification), args.Error(1)
}

func (m *mockQueries) CountNotificationsInTimeWindow(ctx context.Context, arg sqlc.CountNotificationsInTimeWindowParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func TestNotificationRepositoryCreate(t *testing.T) {
	uid := uuid.New()
	createdAt := time.Now().UTC().Round(time.Second)

	mq := new(mockQueries)
	repo := NewNotificationRepository(mq)

	input := entity.Notification{UserID: uid, Type: "status", Message: "test"}
	out := sqlc.Notification{ID: uuid.New(), UserID: uid, Type: "status", Message: "test", CreatedAt: createdAt}

	mq.On("CreateNotification", mock.Anything, sqlc.CreateNotificationParams{
		UserID:  uid,
		Type:    "status",
		Message: "test",
	}).Return(out, nil)

	saved, err := repo.Create(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, uid, saved.UserID)
	require.Equal(t, "status", saved.Type)
	require.Equal(t, "test", saved.Message)
	require.WithinDuration(t, createdAt, saved.CreatedAt, time.Second)

	mq.AssertExpectations(t)
}

func TestNotificationRepositoryCountInTimeWindow(t *testing.T) {
	uid := uuid.New()
	since := time.Now().Add(-time.Minute)

	mq := new(mockQueries)
	repo := NewNotificationRepository(mq)

	mq.On("CountNotificationsInTimeWindow", mock.Anything, sqlc.CountNotificationsInTimeWindowParams{
		UserID:    uid,
		Type:      "status",
		CreatedAt: since,
	}).Return(int64(42), nil)

	count, err := repo.CountInTimeWindow(context.Background(), uid, "status", since)
	require.NoError(t, err)
	require.Equal(t, 42, count)

	mq.AssertExpectations(t)
}
