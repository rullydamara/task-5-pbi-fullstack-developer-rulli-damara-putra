package route

import (
	"task-5-pbi-fullstack-developer-rulli-damara-putra/controllers"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	// User Route
	route.POST("/users/register", controllers.Register)
	route.GET("/users/login", controllers.Login)
	route.POST("/users/logout", controllers.Logout)

	auth_route := route.Group("/")

	auth_route.Use(middlewares.JwtCheck())
	{
		// User Route
		auth_route.PUT("/users/:id", controllers.UpdateUser)
		auth_route.DELETE("/users/:id", controllers.DeleteUser)

		// Photo Route
		auth_route.POST("/photos", controllers.CreatePhoto)
		auth_route.GET("/photos", controllers.GetPhoto)
		auth_route.PUT("/photos/:id", middlewares.PhotoAuthorization(), controllers.UpdatePhoto)
		auth_route.DELETE("/photos/:id", middlewares.PhotoAuthorization(), controllers.DeletePhoto)
	}

	return route
}
