package http

import (
	v1 "github.com/Paulooo0/modak-challenge/internal/adapters/http/v1"
	"github.com/Paulooo0/modak-challenge/internal/domain/useCase"
	"github.com/gin-gonic/gin"
)

func NewRouter(uc *useCase.NotificationUseCase) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	apiV1 := r.Group("/v1")
	v1.RegisterRoutes(apiV1, uc)

	return r
}
