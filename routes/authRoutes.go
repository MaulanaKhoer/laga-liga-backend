package routes

import (
	"github.com/gin-gonic/gin"
	"laga-liga-backend/controllers"
)

// setupAuthRoutes mendaftarkan semua route yang berkaitan dengan autentikasi.
// Semua route di sini bersifat PUBLIC — tidak membutuhkan token JWT.
func setupAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
}
