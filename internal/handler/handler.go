package handler

import (
	"fmt"
	"github.com/Bakhram74/amazon.git/internal/config"
	"github.com/Bakhram74/amazon.git/internal/service"
	"github.com/Bakhram74/amazon.git/pkg/token"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tokenMaker token.Maker
	services   *service.Service
	config     config.Config
}

func NewHandler(config config.Config, services *service.Service) (*Handler, error) {

	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	return &Handler{
		services:   services,
		tokenMaker: tokenMaker,
		config:     config,
	}, nil

}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.renewAccessToken)
	}

	return router
}

func errorResponse(msg string, err error) gin.H {
	return gin.H{msg: err.Error()}
}
