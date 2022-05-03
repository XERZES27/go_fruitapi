package routes

import (
	userController "github.com/XERZES27/go_fruitapi/controllers/user-controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthRoutes(userRouter *gin.RouterGroup, database *mongo.Database) {

	userCollection := database.Collection("users")
	userRouter.POST("/register", userController.Register(userCollection))
	userRouter.POST("/login", userController.Login(userCollection))

}
