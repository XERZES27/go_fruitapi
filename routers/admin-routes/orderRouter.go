package routes

import (
	adminController "github.com/XERZES27/go_fruitapi/controllers/admin-controllers"
	helper "github.com/XERZES27/go_fruitapi/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderRouter(orderRouter *gin.RouterGroup, database *mongo.Database) {
	orderCollection := database.Collection("orders")
	verifyToken := helper.VerifyToken("admin")
	orderRouter.POST("/cancelOrder", verifyToken, adminController.CancelOrder(orderCollection))
	orderRouter.GET("/getOrders", verifyToken, adminController.GetOrders(orderCollection))

}


