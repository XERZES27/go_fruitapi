package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID       primitive.ObjectID `json:"id" binding:"omitempty"  bson:"_id,omitempty" `
	Name     string             `json:"name" binding:"required" bson:"name" form:"name" validate:"required,min=2,max=25"`
	Price    map[string]float64   `json:"price" binding:"required" bson:"price" form:"name" valdiate:"required,price"`
	Photo    string             `json:"photo" binding:"required" bson:"photo" form:"photo" validate:"required,startswith=https://"`
	Unit     string             `json:"unit" binding:"required" bson:"unit" form:"unit" validate:"required,min=1,max=25"`
	Disabled bool               `json:"disabled" binding:"omitempty" bson:"disabled" form:"disabled" `
}
