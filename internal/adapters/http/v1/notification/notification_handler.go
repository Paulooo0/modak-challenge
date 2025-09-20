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

// SendNotification godoc
// @Summary Send a notification
// @Description Sends a notification to a user respecting per-type rate limits
// @Tags notifications
// @Param request body SendNotificationRequest true "Notification payload"
// @Success 201 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/notifications/send [post]
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
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
			c.JSON(http.StatusTooManyRequests, ErrorResponse{Error: errs.ErrRateLimitExceeded.Error()})
			return
		case errs.ErrInvalidNotification:
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: errs.ErrInvalidNotification.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, StatusResponse{Status: "sent"})
}
