package controllers

import (
	"context"
	"net/http"

	model "github.com/XERZES27/go_fruitapi/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

func GetUserById(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.GetQuery("id")
		if !exists {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data"})
			return
		}
		Id, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "Invalid Id"})
			return
		}

		var result bson.M
		err1 := userCollection.FindOne(context.TODO(), bson.M{"_id": Id}).Decode(&result)
		if err1 != nil {
			if err1 == mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Failed to find user"})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err1})
				return
			}
		}
		c.IndentedJSON(http.StatusOK, gin.H{"data": result})

	}
}

type UserStatus struct {
	Id string `json:"id" binding:"required" bson:"id" form:"id" validate:"required,id"`
}

func validateId(fl validator.FieldLevel) bool {
	return primitive.IsValidObjectID(fl.Field().String())

}

func SetUserStatus(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userStatus UserStatus
		e1 := c.ShouldBind(&userStatus)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": e1.Error()})
			return
		}
		validate := validator.New()
		validate.RegisterValidation("id", validateId)
		e2 := validate.Struct(userStatus)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema", "reason": e2.Error()})
			return
		}
		Id, _ := primitive.ObjectIDFromHex(userStatus.Id)
		arrayMatch := []bson.M{{"$set": bson.M{"disabled": bson.M{"$not": "$disabled"}}}}
		UpdateResult, err := userCollection.UpdateByID(context.TODO(), Id,
			arrayMatch)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "Internal Error"})
			return
		}
		if UpdateResult.MatchedCount == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Could not find user"})
			return

		} else if UpdateResult.ModifiedCount == 0 {

			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Did not update user"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": "Update Successful"})

	}
}

func GetUsers(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		lastId, exists := c.GetQuery("lastId")
		query := bson.M{}
		if exists {
			LastId, err := primitive.ObjectIDFromHex(lastId)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Id"})
				return
			}
			query = bson.M{"$match": bson.M{"_id": bson.M{"$gt": LastId}}}
		}
		arrayMatch := []bson.M{
			{"$match": query},
			{"$limit": 20},
			{"$sort": bson.M{"name": 1}},
		}
		showInfoCursor, err := userCollection.Aggregate(context.TODO(), arrayMatch,
			&opts.AggregateOptions{Collation: &opts.Collation{Locale: "en",
				Strength: 2}})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		var users []model.User
		if err = showInfoCursor.All(context.TODO(), &users); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"data": users})

	}
}
