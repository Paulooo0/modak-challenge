package usecase

import (
	"context"
	"time"

	"github.com/Paulooo0/modak-challenge/internal/config/errs"
	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/ports"
)

type NotificationUseCase struct {
	repo    ports.NotificationRepository
	gateway ports.NotificationGateway
	rules   map[string]entity.RateLimit
}

func NewNotificationUseCase(
	repo ports.NotificationRepository,
	gateway ports.NotificationGateway,
	rules map[string]entity.RateLimit,
) *NotificationUseCase {
	return &NotificationUseCase{
		repo:    repo,
		gateway: gateway,
		rules:   rules,
	}
}

func (s *NotificationUseCase) Send(ctx context.Context, n entity.Notification) error {
	rule, ok := s.rules[n.Type]
	if !ok {
		return errs.ErrInvalidNotification
	}

	since := time.Now().Add(-rule.Interval)
	count, err := s.repo.CountInTimeWindow(ctx, n.UserID, n.Type, since)
	if err != nil {
		return err
	}

	if count >= rule.Limit {
		return errs.ErrRateLimitExceeded
	}

	saved, err := s.repo.Create(ctx, n)
	if err != nil {
		return err
	}

	return s.gateway.Send(saved)
}
