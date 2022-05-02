package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(rg *gin.RouterGroup, database *mongo.Database) {
	authRouter := rg.Group("/auth")
	AuthRoutes(authRouter, database)

}
