package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(rg *gin.RouterGroup, database *mongo.Database) {
	AuthRoutes(rg.Group("/auth"), database)
	AccountRouter(rg.Group("/account"), database)
	ProductRouter(rg.Group("/product"), database)
	OrderRouter(rg.Group("/order"), database)

}
