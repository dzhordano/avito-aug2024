package v1

import (
	_ "github.com/dzhordano/avito-bootcamp2024/docs"
	"github.com/dzhordano/avito-bootcamp2024/internal/service"
	"github.com/dzhordano/avito-bootcamp2024/pkg/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services      *service.Services
	tokensManager auth.TokensManager
}

func NewHandler(services *service.Services, tokensManger auth.TokensManager) *Handler {
	return &Handler{
		services:      services,
		tokensManager: tokensManger,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("")
	{
		h.initAuthRoutes(v1)
		h.initHouseRoutes(v1)
		h.initFlatRoutes(v1)
	}
}
