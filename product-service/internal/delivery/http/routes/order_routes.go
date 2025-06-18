package routes

import (
	"github.com/gin-gonic/gin"
	"product-service/internal/delivery/http/controllers"
	"product-service/internal/delivery/http/middleware"
)

func RegisterOrderRoutes(r *gin.Engine, OrderController *controllers.OrderController) {
	privateRoutes := r.Group("/api/orders")
	privateRoutes.Use(middleware.AuthMiddleWare())
	{
		privateRoutes.POST("/", OrderController.CreateOrder)
		privateRoutes.GET("/:id", OrderController.GetOrderByID)
		privateRoutes.PATCH("/:id", OrderController.UpdateOrder)
		privateRoutes.GET("/user/:id", OrderController.GetOrdersByUserID)
	}
}
