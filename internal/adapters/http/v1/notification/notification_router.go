package notification

import (
	"github.com/Paulooo0/modak-challenge/internal/domain/useCase"
	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(r *gin.RouterGroup, uc *useCase.NotificationUseCase) {
	h := NewNotificationHandler(uc)

	api := r.Group("/notifications")
	{
		api.POST("/send", h.SendNotification)
	}
}
