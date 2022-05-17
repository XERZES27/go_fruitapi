package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	model "github.com/XERZES27/go_fruitapi/models"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

func validatePrice(priceMap map[string]float64) bool {
	if len(priceMap) == 0 {
		return false
	}
	for quantity, price := range priceMap {
		currentQuantity, err := strconv.ParseFloat(quantity, 64)

		if err != nil {
			return false
		}
		if currentQuantity <= 0 || price <= 0 {
			return false
		}
	}
	return true
}

func CreateProduct(productCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.Product

		err := c.ShouldBind(&product)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": err.Error()})
			return
		}

		validate := validator.New()
		e2 := validate.Struct(product)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema"})
			return
		}
		if !validatePrice(product.Price) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema"})
			return
		}

		productResult, e3 := productCollection.InsertOne(context.TODO(), product)
		if e3 != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"ID": productResult.InsertedID})

	}
}

func EditProduct(productCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product model.EditProduct
		err := c.ShouldBind(&product)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": err.Error()})
			return
		}
		validate := validator.New()
		e2 := validate.Struct(product)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema", "reason": e2.Error()})
			return
		}

		if product.Price != nil {
			if !validatePrice(product.Price) {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema"})
				return
			}
		}
		updateCommand := bson.M{"disabled": product.Disabled}

		if product.Name != "" {

			updateCommand["name"] = product.Name
		}
		if product.Name != "" {
			updateCommand["photo"] = product.Photo
		}
		if product.Name != "" {
			updateCommand["unit"] = product.Unit
		}
		if product.Price != nil {
			updateCommand["price"] = product.Price
		}

		UpdateResult, err := productCollection.UpdateByID(context.TODO(), product.ID, bson.M{"$set": updateCommand})

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		if UpdateResult.MatchedCount == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Could not find product"})
			return
		} else if UpdateResult.ModifiedCount == 0 {

			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Did not update product"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": "Update Successful"})
	}
}

type deleteProduct struct {
	ID primitive.ObjectID `json:"id" binding:"required"  bson:"_id,required" `
}

func DeleteProduct(productCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product deleteProduct
		err := c.ShouldBind(&product)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": err.Error()})
			return
		}
		fmt.Println(product.ID.Hex())
		if !primitive.IsValidObjectID(product.ID.Hex()) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data"})
			return

		}

		DeleteResult, err := productCollection.DeleteOne(context.TODO(), bson.M{"_id": product.ID})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		if DeleteResult.DeletedCount == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Did not find product"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": "Delete Successful"})

	}

}

func GetProductById(productCollection *mongo.Collection) gin.HandlerFunc {
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
		err1 := productCollection.FindOne(context.TODO(), bson.M{"_id": Id}).Decode(&result)
		if err1 != nil {
			if err1 == mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Failed to find product"})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err1})
				return
			}
		}
		c.IndentedJSON(http.StatusOK, gin.H{"data": result})

	}
}

func GetProducts(productCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		arrayMatch := []bson.M{
			{"$sort":bson.M{"name":1}},
		}
		showInfoCursor, err := productCollection.Aggregate(context.TODO(),arrayMatch ,
			&opts.AggregateOptions{Collation: &opts.Collation{Locale: "en",
				Strength: 2}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		}

		var products []model.Product
		if err = showInfoCursor.All(context.TODO(), &products); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"data": products})

	}
}
