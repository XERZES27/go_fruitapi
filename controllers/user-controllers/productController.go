package controllers

import (
	"context"
	"net/http"

	model "github.com/XERZES27/go_fruitapi/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

func GetProducts(productCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		arrayMatch := []bson.M{{"$match": bson.M{"disabled": false}},
			{"$project": bson.M{"disabled": 0}}, {"$sort": bson.M{"name": 1}}}
		showInfoCursor, err := productCollection.Aggregate(context.TODO(), arrayMatch,
			&opts.AggregateOptions{Collation: &opts.Collation{Locale: "en",
				Strength: 2}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		var products []model.Product
		if err = showInfoCursor.All(context.TODO(), &products); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"data": products})

	}
}
