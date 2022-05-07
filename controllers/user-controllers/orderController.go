package controllers

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	model "github.com/XERZES27/go_fruitapi/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func orderValidator(order *map[string]float64) bool {
	for key, element := range *order {
		if !primitive.IsValidObjectID(key) || element <= 0 {
			delete(*order, key)
		}
	}
	return len(*order) > 0
}

func calculatePrice(products *[]model.Product, order *map[string]float64, comments *[]string) (float64, error) {
	var price float64
	for _, product := range *products {
		var quantityPrice float64
		quantity := (*order)[product.ID.Hex()]
		for quantityStr, price := range product.Price {
			currentQuantity, err := strconv.ParseFloat(quantityStr, 64)

			if err != nil {
				return 0, err
			}
			if quantity > currentQuantity {
				quantityPrice = price
			} else {
				quantityPrice = price
				break
			}
		}
		price += quantityPrice * quantity
		pluralityIdentifier := "s of "
		if strings.HasSuffix(product.Unit, "s") {
			pluralityIdentifier = " of "
		}
		*comments = append(*comments, fmt.Sprint(quantity)+" "+product.Unit+pluralityIdentifier+product.Name)

	}
	return price, nil
}

func covertToObjectKeys(keys []reflect.Value, objectIdKeys *[]primitive.ObjectID) {

	for _, v := range keys {
		key, _ := primitive.ObjectIDFromHex(fmt.Sprint(v))
		*objectIdKeys = append(*objectIdKeys, key)
	}

}

func CreateOrder(orderCollection *mongo.Collection, productCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order model.OrderJson
		e1 := c.ShouldBind(&order)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Missing query parameter", "reason": e1.Error()})
			return
		}
		if !orderValidator(&order.Order) {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Orders", "reason": e1.Error()})
			return
		}

		keys := reflect.ValueOf(order.Order).MapKeys()

		var objectIdKeys []primitive.ObjectID
		covertToObjectKeys(keys, &objectIdKeys)
		if len(objectIdKeys) == 0 {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "Invalid Orders"})
			return
		}
		arrayMatch := []bson.M{{"$match": bson.M{"_id": bson.M{"$in": objectIdKeys}}}}
		showInfoCursor, err := productCollection.Aggregate(context.TODO(), arrayMatch)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		var products []model.Product
		if err = showInfoCursor.All(context.TODO(), &products); err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		if len(products) == 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Could not find associated products"})
			return
		}
		var comments []string
		price, err := calculatePrice(&products, &order.Order, &comments)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			return
		}
		CompanyName, exists := c.Get("CompanyName")
		if !exists {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "fatal token error1"})
			return
		}
		Id, _ := c.Get("Id")
		validId, err := primitive.ObjectIDFromHex(fmt.Sprint(Id))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "fatal token error2", "reason": err.Error()})
			return
		}
		ord := &model.Order{
			Name:              fmt.Sprint(CompanyName),
			UserId:            validId,
			ProductQuantities: order.Order,
			Price:             price,
			Comments:          comments,
			Canceled:          false,
			Date:              primitive.NewDateTimeFromTime(time.Now().UTC()),
		}
		orderResult, err := orderCollection.InsertOne(context.TODO(), ord)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"id": orderResult.InsertedID})

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

		} else if UpdateResult.ModifiedCount == 0 {

			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Did not update order"})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"status": "Update Successful"})

	}
}

func GetOrder(orderCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query bson.M
		userId, _ := c.Get("Id")
		UserId, err := primitive.ObjectIDFromHex(fmt.Sprint(userId))
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		lastId, exists := c.GetQuery("lastId")
		if exists {
			LastId, err := primitive.ObjectIDFromHex(lastId)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Id"})
				return
			}
			query = bson.M{"userId": UserId, "_id": bson.M{"$gt": LastId}}
		} else {
			query = bson.M{"userId": UserId}
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
