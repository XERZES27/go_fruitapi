package routes

import (
	productController "github.com/XERZES27/go_fruitapi/controllers/admin-controllers"
	helper "github.com/XERZES27/go_fruitapi/helpers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func ProductRouter(productRouter *gin.RouterGroup, database *mongo.Database) {
	verifyToken := helper.VerifyToken("admin")
	productsCollection := database.Collection("products")
	productRouter.POST("/create", verifyToken, productController.CreateProduct(productsCollection))
	productRouter.POST("/edit", verifyToken, productController.EditProduct(productsCollection))
	productRouter.GET("/getProducts", verifyToken, productController.GetProducts(productsCollection))
	productRouter.GET("/getProductById", verifyToken, productController.GetProductById(productsCollection))
	productRouter.DELETE("/delete", verifyToken, productController.DeleteProduct(productsCollection))

}
