package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID       primitive.ObjectID `json:"id" binding:"omitempty"  bson:"_id,omitempty" `
	Name     string             `json:"name" binding:"required" bson:"name" form:"name" validate:"required,min=2,max=25"`
	Price    map[string]float64 `json:"price" binding:"omitempty" bson:"price" form:"name" valdiate:"omitempty,price"`
	Photo    string             `json:"photo" binding:"required" bson:"photo" form:"photo" validate:"required,startswith=https://"`
	Unit     string             `json:"unit" binding:"required" bson:"unit" form:"unit" validate:"required,min=1,max=25"`
	Disabled bool               `json:"disabled" binding:"omitempty" bson:"disabled" form:"disabled" `
}

type EditProduct struct {
	ID       primitive.ObjectID `json:"id" binding:"required"   `
	Name     string             `json:"name" binding:"omitempty" bson:"name,omitempty" form:"name" validate:"omitempty,min=2,max=25"`
	Price    map[string]float64 `json:"price" binding:"omitempty" bson:"price,omitempty" form:"name" valdiate:"omitempty,price"`
	Photo    string             `json:"photo" binding:"omitempty" bson:"photo,omitempty" form:"photo" validate:"omitempty,startswith=https://"`
	Unit     string             `json:"unit" binding:"omitempty" bson:"unit,omitempty" form:"unit" validate:"omitempty,min=1,max=25"`
	Disabled bool               `json:"disabled" binding:"omitempty" bson:"disabled,omitempty" form:"disabled" validate:"omitempty"`
}
