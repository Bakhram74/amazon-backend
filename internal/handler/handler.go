package handler

import (
	"fmt"
	"github.com/Bakhram74/amazon-backend.git/internal/config"
	"github.com/Bakhram74/amazon-backend.git/internal/service"
	"github.com/Bakhram74/amazon-backend.git/pkg/token"

	"github.com/gin-contrib/cors"

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
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"https://localhost:3000"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("OPTIONS")

	// Register the middleware
	router.Use(cors.New(corsConfig))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.renewAccessToken)
	}
	api := router.Group("/api")
	{
		user := api.Group("/").Use(authMiddleware(h.tokenMaker))
		{
			user.GET("/profile", h.getProfile)
			//user.POST("/profile", h.updateProfile)
			//user.PATCH("/profile/favorites/:productId", h.getProfile)
		}

	}

	return router
}
