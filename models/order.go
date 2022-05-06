package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID                primitive.ObjectID `json:"id" binding:"omitempty"  bson:"_id,omitempty" `
	Name              string             `json:"name" binding:"required" bson:"name" form:"name" validate:"required,min=2,max=25"`
	UserId            primitive.ObjectID `json:"userId" binding:"required"  bson:"userId" validate:"required"`
	ProductQuantities map[string]float64   `json:"productQuantities" binding:"required" bson:"productQuantities" form:"productQuantities" valdiate:"required,productQuantities"`
	Price             float64              `json:"price" binding:"required" bson:"price" form:"price" validate:"required"`
	Comments          []string           `json:"comments" binding:"required" bson:"comments" form:"comments" validate:"required"`
	Canceled          bool               `json:"canceled" binding:"required" bson:"canceled" form:"canceled" `
	Date              primitive.DateTime `json:"date" binding:"omitempty" bson:"date,omitempty"`
}

type OrderJson struct {
	Order map[string]float64 `json:"order" binding:"required" bson:"order" form:"order" validate:"required"`
}
