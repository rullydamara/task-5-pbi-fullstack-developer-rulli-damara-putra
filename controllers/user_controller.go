package controllers

import (
	"net/http"
	"strconv"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/app"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/database"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/helpers"
	"task-5-pbi-fullstack-developer-rulli-damara-putra/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

func Register(context *gin.Context) {
	conn, error := database.Setup()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Error while connection to the database",
		})
		return
	}

	var newUser app.UserData
	if error := context.BindJSON(&newUser); error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials, try again",
		})
		return
	}

	insertUser := models.User{
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: helpers.EncryptPassword(newUser.Password),
	}

	_, error = govalidator.ValidateStruct(insertUser)

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": error.Error(),
		})
		return
	}

	var checkEmail models.User
	conn.Where("email = ?", newUser.Email).First(&checkEmail)

	var checkUsername models.User
	conn.Where("username = ?", newUser.Username).First(&checkUsername)

	if checkEmail.Email != "" || checkUsername.Username != "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Username or Email already exist",
		})
		return
	}

	result := conn.Create(&insertUser)

	if result.Error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Username or Email already exist",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully registered",
	})
}

func Login(context *gin.Context) {
	conn, error := database.Setup()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Error while connection to the database",
		})

		return
	}

	var user models.User

	email := context.Query("email")
	password := context.Query("password")

	err := conn.Where("email = ?", email).First(&user).Error

	if err != nil || !helpers.CheckPassword(password, user.Password) {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	token, err := helpers.GenerateToken(user)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Error while generating token",
		})
		return
	}

	context.SetCookie("Authorization", token, 3600, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully logged in",
	})
}

func Logout(context *gin.Context) {

	_, err := context.Cookie("Authorization")

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "Unauthorized",
		})
		return
	}

	context.SetCookie("Authorization", "", -1, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully Logout",
	})
}

func UpdateUser(context *gin.Context) {

	conn, error := database.Setup()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Error while connection to the database",
		})
		return
	}

	updateID := context.Param("id")

	var newUser app.UserData
	if error := context.ShouldBindJSON(&newUser); error != nil {
		context.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  404,
			"message": "Invalid credentials",
		})
		return
	}

	userData := context.MustGet("user").(models.User)

	if strconv.Itoa(userData.ID) != updateID {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "You dont have access to update this user",
		})
		return
	}

	var user models.User
	conn.Where("id = ?", updateID).First(&user)

	var checkEmail models.User
	conn.Where("email = ?", newUser.Email).First(&checkEmail)

	var checkUsername models.User
	conn.Where("username = ?", newUser.Username).First(&checkUsername)

	if checkEmail.Email != "" || checkUsername.Username != "" {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "Username or Email already exist",
		})
		return
	}

	user.Username = newUser.Username
	user.Email = newUser.Email
	user.Password = helpers.EncryptPassword(newUser.Password)

	_, error = govalidator.ValidateStruct(user)

	if error != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": error.Error(),
		})
		return
	}

	conn.Save(&user)

	context.SetCookie("Authorization", "", -1, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully update user, please login to related user if you want continue update the user",
	})

}

func DeleteUser(context *gin.Context) {

	conn, error := database.Setup()

	if error != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Error while connection to the database",
		})
		return
	}

	deleteID := context.Param("id")

	var user models.User
	conn.Where("id = ?", deleteID).First(&user)

	userData := context.MustGet("user").(models.User)

	if strconv.Itoa(userData.ID) != deleteID {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": "You dont have access to delete this user, please login to related user if you want continue delete the user",
		})
		return
	}

	conn.Delete(&user)

	context.SetCookie("Authorization", "", -1, "", "", true, true)

	context.IndentedJSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Successfully delete this user, please register or login again to continue",
	})

}
