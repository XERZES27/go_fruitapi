package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID       primitive.ObjectID `json:"id" binding:"omitempty"  bson:"_id,omitempty" `
	Name     string             `json:"name" binding:"required" bson:"name" form:"name" validate:"required,eq=ADMIN_ZERO"`
	Password string             `json:"password" binding:"required" bson:"password" form:"password" validate:"required,password"`
	Date     primitive.DateTime `json:"date" binding:"omitempty" bson:"date,omitempty"`
}
type LoginAdmin struct {
	Name     string `json:"name" binding:"required" validate:"required,eq=ADMIN_ZERO"`
	Password string `json:"password" binding:"required"  validate:"required,password"`
}
