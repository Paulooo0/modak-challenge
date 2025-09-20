package notification

import (
	"context"
	"log"
	"net/http"

	"github.com/Paulooo0/modak-challenge/internal/config/errs"
	"github.com/Paulooo0/modak-challenge/internal/domain/entity"
	"github.com/Paulooo0/modak-challenge/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	uc *usecase.NotificationUseCase
}

func NewNotificationHandler(uc *usecase.NotificationUseCase) *NotificationHandler {
	return &NotificationHandler{uc: uc}
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req struct {
		UserID  uuid.UUID `json:"user_id" binding:"required"`
		Type    string    `json:"type" binding:"required"`
		Message string    `json:"message" binding:"required"`
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
		log.Println(err)
		switch err {
		case errs.ErrRateLimitExceeded:
			c.JSON(http.StatusTooManyRequests, gin.H{"error": errs.ErrRateLimitExceeded.Error()})
			return
		case errs.ErrInvalidNotification:
			c.JSON(http.StatusBadRequest, gin.H{"error": errs.ErrInvalidNotification.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"status": "sent"})
}
