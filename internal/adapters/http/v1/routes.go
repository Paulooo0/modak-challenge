package v1

import (
	"github.com/Paulooo0/modak-challenge/internal/adapters/http/v1/notification"
	"github.com/Paulooo0/modak-challenge/internal/domain/useCase"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, uc *useCase.NotificationUseCase) {
	notification.RegisterNotificationRoutes(r, uc)
}
