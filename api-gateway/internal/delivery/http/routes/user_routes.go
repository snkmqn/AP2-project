package routes

import (
	"api-gateway/internal/delivery/http/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userController *controllers.UserController) {
	public := r.Group("/api")
	{
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)
	}
}
