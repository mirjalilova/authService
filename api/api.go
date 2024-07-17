package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/mirjalilova/authService/api/handlers"
)

// @title Authentication Service API
// @version 1.0
// @description API for Authentication Service
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func Engine(handler *handlers.Handlers) *gin.Engine {
	router := gin.Default()

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/register", handler.RegisterUser)
	router.POST("/login", handler.LoginUser)
	router.POST("/forgot-password", handler.ForgotPassword)
	router.POST("/reset-password", handler.ResetPassword)

	return router
}
