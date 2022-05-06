package routes

import (
	userController "github.com/XERZES27/go_fruitapi/controllers/user-controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	helper "github.com/XERZES27/go_fruitapi/helpers"
)

func ProductRouter(productRouter *gin.RouterGroup, database *mongo.Database){
	productCollection := database.Collection("products")
	productRouter.GET("/getProducts",helper.VerifyToken("user"), userController.GetProducts(productCollection))
}