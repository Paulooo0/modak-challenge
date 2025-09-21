package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Paulooo0/modak-challenge/internal/config/errs"
	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/domain/usecase"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, n entity.Notification) (entity.Notification, error) {
	args := m.Called(ctx, n)
	return args.Get(0).(entity.Notification), args.Error(1)
}

func (m *MockRepo) CountInTimeWindow(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, since time.Time) (int, error) {
	args := m.Called(ctx, userID, notifType, since)
	return args.Get(0).(int), args.Error(1)
}

type MockGateway struct {
	mock.Mock
}

func (m *MockGateway) Send(n entity.Notification) error {
	args := m.Called(n)
	return args.Error(0)
}

func TestSendNotificationSuccess(t *testing.T) {
	repo := new(MockRepo)
	gw := new(MockGateway)
	userID := uuid.New()

	repo.On("CountInTimeWindow", mock.Anything, userID, entity.Status, mock.AnythingOfType("time.Time")).Return(0, nil)
	created := entity.Notification{UserID: userID, Type: entity.Status, Message: "hello"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("entity.Notification")).Return(created, nil)
	gw.On("Send", mock.MatchedBy(func(n entity.Notification) bool {
		return n.UserID == userID && n.Type == entity.Status && n.Message == "hello"
	})).Return(nil)

	svc := usecase.NewNotificationUseCase(repo, gw, entity.DefaultRateLimits)

	err := svc.Send(context.Background(), entity.Notification{
		UserID:  userID,
		Type:    entity.Status,
		Message: "hello",
	})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	gw.AssertExpectations(t)
}

func TestNotificationRateLimitExceeded(t *testing.T) {
	repo := new(MockRepo)
	gw := new(MockGateway)
	userID := uuid.New()

	repo.On("CountInTimeWindow", mock.Anything, userID, entity.Status, mock.Anything).Return(2, nil)

	svc := usecase.NewNotificationUseCase(repo, gw, entity.DefaultRateLimits)

	err := svc.Send(context.Background(), entity.Notification{
		UserID:  userID,
		Type:    entity.Status,
		Message: "hello",
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, errs.ErrRateLimitExceeded))
	repo.AssertExpectations(t)
	gw.AssertNotCalled(t, "Send", mock.Anything)
}

func TestSendNotificationGatewayError(t *testing.T) {
	repo := new(MockRepo)
	gw := new(MockGateway)
	userID := uuid.New()

	repo.On("CountInTimeWindow", mock.Anything, userID, entity.Status, mock.Anything).Return(0, nil)

	created := entity.Notification{UserID: userID, Type: entity.Status, Message: "hello"}
	repo.On("Create", mock.Anything, mock.AnythingOfType("entity.Notification")).Return(created, nil)

	gw.On("Send", created).Return(errors.New("gateway down"))

	svc := usecase.NewNotificationUseCase(repo, gw, entity.DefaultRateLimits)

	err := svc.Send(context.Background(), entity.Notification{
		UserID:  userID,
		Type:    entity.Status,
		Message: "hello",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gateway down")

	repo.AssertExpectations(t)
	gw.AssertExpectations(t)
}
