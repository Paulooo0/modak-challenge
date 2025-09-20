package notification

import (
	"github.com/Paulooo0/modak-challenge/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

func RegisterNotificationRoutes(r *gin.RouterGroup, uc *usecase.NotificationUseCase) {
	h := NewNotificationHandler(uc)

	api := r.Group("/notifications")
	{
		api.POST("/send", h.SendNotification)
	}
}
