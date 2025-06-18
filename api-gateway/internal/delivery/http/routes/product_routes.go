package routes

import (
	"api-gateway/internal/delivery/http/controllers"
	"api-gateway/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, productController *controllers.ProductController) {
	publicRoutes := r.Group("/api/products")
	{
		publicRoutes.GET("/:id", productController.GetProductByID)
		publicRoutes.GET("/", productController.ListProducts)
	}

	privateRoutes := r.Group("/api/products")
	privateRoutes.Use(middleware.AuthMiddleWare())
	{
		privateRoutes.POST("/", productController.CreateProduct)
		privateRoutes.PATCH("/:id", productController.UpdateProduct)
		privateRoutes.DELETE("/:id", productController.DeleteProduct)
	}
}
