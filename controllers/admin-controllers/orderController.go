package controllers

import (
	"context"
	"net/http"

	model "github.com/XERZES27/go_fruitapi/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetOrders(orderCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query bson.M
		lastId, exists := c.GetQuery("lastId")
		if exists {
			LastId, err := primitive.ObjectIDFromHex(lastId)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Id"})
				return
			}
			query = bson.M{"_id": bson.M{"$gt": LastId}}
		} else {
			query = bson.M{}
		}

		arrayMatch := []bson.M{
			{"$match": query},
			{"$limit": 20},
			{"$sort": bson.M{"date": -1}},
		}
		showInfoCursor, err := orderCollection.Aggregate(context.TODO(), arrayMatch)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		var orders []model.Order
		if err = showInfoCursor.All(context.TODO(), &orders); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"stauts": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"data": orders})

	}
}

type cancelBody struct {
	Id string `json:"id" binding:"required"`
}

func CancelOrder(orderCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body cancelBody
		e1 := c.ShouldBind(&body)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": e1.Error()})
			return
		}
		Id, err := primitive.ObjectIDFromHex(body.Id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Id"})
			return
		}


		UpdateResult, err := orderCollection.UpdateByID(context.TODO(), Id, bson.M{"$set": bson.M{"canceled": true}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		if UpdateResult.MatchedCount == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Could not find order"})
			return
		} else if UpdateResult.ModifiedCount == 0 {

			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Did not update order"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": "Update Successful"})
	}
}
