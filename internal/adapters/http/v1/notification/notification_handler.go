package notification

import (
	"context"
	"net/http"

	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/domain/useCase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	uc *useCase.NotificationUseCase
}

func NewNotificationHandler(uc *useCase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req struct {
		UserID  string `json:"user_id" binding:"required"`
		Type    string `json:"type" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n := entity.Notification{
		ID:      uuid.New(),
		UserID:  req.UserID,
		Type:    req.Type,
		Message: req.Message,
	}

	err := h.uc.Send(context.Background(), n)
	if err != nil {
		if err == useCase.ErrRateLimitExceeded {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "sent"})
}
