package middlewares

import (
	"net/http"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/database"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/helpers"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/models"
	"time"

	"github.com/gin-gonic/gin"
)

func JwtCheck() gin.HandlerFunc {

	return func(context *gin.Context) {
		auth_token, err := context.Cookie("Authorization")

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "Unauthorized. Token not found.",
			})
			return
		}

		claims, err := helpers.ParseToken(auth_token)

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  401,
				"message": "Unauthorized",
			})
			return
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  4000,
				"message": "Token Expired",
			})
			return
		}

		var user models.User
		conn, err := database.Setup()

		conn.Where("email = ?", claims["email"]).First(&user)

		if user.ID == 0 || err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  4000,
				"message": "Token Invalid",
			})
			return
		}

		context.Set("user", user)
		context.Next()
	}
}
