package http

import (
	v1 "github.com/Paulooo0/modak-challenge/internal/adapters/http/v1"
	"github.com/Paulooo0/modak-challenge/internal/domain/usecase"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(uc *usecase.NotificationUseCase) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := r.Group("/v1")
	v1.RegisterRoutes(apiV1, uc)

	return r
}
