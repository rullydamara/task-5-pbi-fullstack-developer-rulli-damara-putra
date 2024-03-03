package middlewares

import (
	"net/http"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/database"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/models"

	"github.com/gin-gonic/gin"
)

func PhotoAuthorization() gin.HandlerFunc {

	return func(context *gin.Context) {

		photoId := context.Param("id")

		conn, error := database.Setup()

		if error != nil {
			context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"status":  500,
				"message": "Error while connection to the database",
			})
			return
		}

		var photo models.Photo
		conn.Where("id = ?", photoId).First(&photo)

		if photo.ID == 0 {
			context.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "Photo not found",
			})
			return
		}

		userData := context.MustGet("user").(models.User)

		if photo.UserID != userData.ID {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "You dont have access to this photo",
			})
			return
		}

		context.Next()
	}
}
