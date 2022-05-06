package routes

import (
	userController "github.com/XERZES27/go_fruitapi/controllers/user-controllers"
	helper "github.com/XERZES27/go_fruitapi/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderRouter(orderRouter *gin.RouterGroup, database *mongo.Database) {
	productCollection := database.Collection("products")
	orderCollection := database.Collection("orders")
	verifyToken := helper.VerifyToken("user")
	orderRouter.POST("/createOrder", verifyToken, userController.CreateOrder(orderCollection, productCollection))
	orderRouter.POST("/cancelOrder", verifyToken, userController.CancelOrder(orderCollection))
	orderRouter.GET("/getOrders", verifyToken, userController.GetOrder(orderCollection))

}
