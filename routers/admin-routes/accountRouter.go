package routes

import (
	adminController "github.com/XERZES27/go_fruitapi/controllers/admin-controllers"
	helper "github.com/XERZES27/go_fruitapi/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AccountRouter(adminRouter *gin.RouterGroup, database *mongo.Database) {
	verifyToken := helper.VerifyToken("admin")
	userCollection := database.Collection("users")
	adminRouter.GET("/getUserById", verifyToken, adminController.GetUserById(userCollection))
	adminRouter.GET("/getUsers", verifyToken, adminController.GetUsers(userCollection))
	adminRouter.POST("/setUserStatus", verifyToken, adminController.SetUserStatus(userCollection))

}
