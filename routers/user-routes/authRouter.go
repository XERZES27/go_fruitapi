package main

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthRoutes(rg *gin.RouterGroup, database *mongo.Database) {
	
	userCollection := database.Collection("users")

}