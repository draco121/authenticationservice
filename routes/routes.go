package routes

import (
	"github.com/draco121/authenticationservice/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(controllers controllers.Controllers, router *gin.Engine) {
	v1 := router.Group("/v1")
	v1.POST("/login", controllers.Login)
	v1.GET("/authenticate", controllers.Authenticate)
	v1.POST("/refresh", controllers.RefreshLogin)
	v1.POST("/logout", controllers.Logout)
}
