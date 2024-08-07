package http

import (
	v1 "github.com/dzhordano/avito-bootcamp2024/internal/delivery/http/v1"
	"github.com/dzhordano/avito-bootcamp2024/internal/service"
	"github.com/dzhordano/avito-bootcamp2024/pkg/auth"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Handler struct {
	services      *service.Services
	tokensManager auth.TokensManager
}

func NewHandler(services *service.Services, tokensManager auth.TokensManager) *Handler {
	return &Handler{
		services:      services,
		tokensManager: tokensManager,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	handlerV1 := v1.NewHandler(h.services, h.tokensManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
