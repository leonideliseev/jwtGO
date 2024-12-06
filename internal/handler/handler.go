package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/leonideliseev/jwtGO/internal/service"
)

type Handler struct {
	serv *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{
		serv: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.GET("/tokens", h.createTokens)
		auth.GET("/refresh", h.refreshTokens)
	}

	return router
}
