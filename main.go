package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	adminRoutes "github.com/XERZES27/go_fruitapi/routers/admin-routes"
	userRoutes "github.com/XERZES27/go_fruitapi/routers/user-routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	router = gin.Default()
)

// Run will start the server
func Run(database *mongo.Database) {
	getRoutes(database)
	router.Run("localhost:5000")
}

func getRoutes(database *mongo.Database) {
	userRouter := router.Group("/user")
	userRoutes.Routes(userRouter, database)

	adminRouter := router.Group("/admin")
	adminRoutes.Routes(adminRouter, database)
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	uri := os.Getenv("DB_CONNECT")
	client, err := mongo.Connect(context.TODO(), options.Client().SetConnectTimeout(15*time.Second).ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Println("closing connection")
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	database := client.Database("myFirstDatabase")


	Run(database)

}
