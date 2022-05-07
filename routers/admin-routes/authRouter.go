package routes

import (
	adminController "github.com/XERZES27/go_fruitapi/controllers/admin-controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthRoutes(userRouter *gin.RouterGroup, database *mongo.Database) {

	adminCollection := database.Collection("admins")
	userRouter.POST("/register", adminController.Register(adminCollection))
	userRouter.POST("/login", adminController.Login(adminCollection))

}
