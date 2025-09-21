package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockRepo struct{ mock.Mock }

func (m *MockRepo) Create(ctx context.Context, n entity.Notification) (entity.Notification, error) {
	args := m.Called(ctx, n)
	return args.Get(0).(entity.Notification), args.Error(1)
}

func (m *MockRepo) CountInTimeWindow(ctx context.Context, userID uuid.UUID, notifType entity.NotificationType, since time.Time) (int, error) {
	args := m.Called(ctx, userID, notifType, since)
	return args.Get(0).(int), args.Error(1)
}

type MockGateway struct{ mock.Mock }

func (m *MockGateway) Send(n entity.Notification) error {
	args := m.Called(n)
	return args.Error(0)
}

func buildHandler(repo *MockRepo, gw *MockGateway, rules map[entity.NotificationType]entity.RateLimit) *NotificationHandler {
	uc := usecase.NewNotificationUseCase(repo, gw, rules)
	return NewNotificationHandler(uc)
}

const (
	pathSend          = "/v1/notifications/send"
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

type sendPayload struct {
	UserID  uuid.UUID `json:"user_id"`
	Type    string    `json:"type"`
	Message string    `json:"message"`
}

func newJSONRequest(t testing.TB, method, path string, v any) *http.Request {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	body := bytes.NewBuffer(b)
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(headerContentType, contentTypeJSON)
	return req
}

func TestSendNotificationSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepo)
	gw := new(MockGateway)
	rules := map[entity.NotificationType]entity.RateLimit{entity.Status: {Limit: 10, Interval: time.Minute}}
	h := buildHandler(repo, gw, rules)

	var captured entity.Notification
	repo.On("CountInTimeWindow", mock.Anything, mock.AnythingOfType("uuid.UUID"), entity.Status, mock.AnythingOfType("time.Time")).Return(0, nil)
	repo.On("Create", mock.Anything, mock.MatchedBy(func(n entity.Notification) bool { captured = n; return true })).Return(captured, nil)
	gw.On("Send", mock.AnythingOfType("entity.Notification")).Return(nil)

	r := gin.New()
	w := httptest.NewRecorder()
	r.POST(pathSend, h.SendNotification)

	req := newJSONRequest(t, http.MethodPost, pathSend, sendPayload{UserID: uuid.New(), Type: string(entity.Status), Message: "hello"})

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	repo.AssertExpectations(t)
	gw.AssertExpectations(t)
}

func TestSendNotificationRateLimited(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepo)
	gw := new(MockGateway)
	rules := map[entity.NotificationType]entity.RateLimit{entity.Status: {Limit: 1, Interval: time.Minute}}
	h := buildHandler(repo, gw, rules)

	repo.On("CountInTimeWindow", mock.Anything, mock.AnythingOfType("uuid.UUID"), entity.Status, mock.AnythingOfType("time.Time")).Return(1, nil)

	r := gin.New()
	w := httptest.NewRecorder()
	r.POST(pathSend, h.SendNotification)

	req := newJSONRequest(t, http.MethodPost, pathSend, sendPayload{UserID: uuid.New(), Type: string(entity.Status), Message: "hello"})

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestSendNotificationInvalidType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepo)
	gw := new(MockGateway)
	rules := map[entity.NotificationType]entity.RateLimit{"invalidType": {Limit: 10, Interval: time.Minute}}
	h := buildHandler(repo, gw, rules)

	r := gin.New()
	w := httptest.NewRecorder()
	r.POST(pathSend, h.SendNotification)

	req := newJSONRequest(t, http.MethodPost, pathSend, sendPayload{UserID: uuid.New(), Type: "invalidType", Message: "hello"})

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSendNotificationInternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := new(MockRepo)
	gw := new(MockGateway)
	rules := map[entity.NotificationType]entity.RateLimit{entity.Status: {Limit: 10, Interval: time.Minute}}
	h := buildHandler(repo, gw, rules)

	boom := errors.New("db error")
	repo.On("CountInTimeWindow", mock.Anything, mock.AnythingOfType("uuid.UUID"), entity.Status, mock.AnythingOfType("time.Time")).Return(0, nil)
	repo.On("Create", mock.Anything, mock.AnythingOfType("entity.Notification")).Return(entity.Notification{}, boom)

	r := gin.New()
	w := httptest.NewRecorder()
	r.POST(pathSend, h.SendNotification)

	req := newJSONRequest(t, http.MethodPost, pathSend, sendPayload{UserID: uuid.New(), Type: string(entity.Status), Message: "hello"})

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
