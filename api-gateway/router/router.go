package router

import (
	"api-gateway/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	router *gin.Engine,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	userHandler *handlers.UserHandler,
) {
	router.POST("/products", productHandler.CreateProduct)
	router.POST("/orders", orderHandler.CreateOrder)
	router.POST("/users/register", userHandler.RegisterUser)
	router.POST("/users/login", userHandler.LoginUserHandler)
	router.GET("/users/:user_id", userHandler.RetrieveProfileHandler)
}
