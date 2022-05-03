package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" binding:"omitempty"  bson:"_id,omitempty" `
	CompanyName  string             `json:"የኩባኒኛ_ስም" binding:"required" bson:"የኩባኒኛ_ስም" form:"የኩባኒኛ_ስም" validate:"required,min=6,max=255"`
	Type         string             `json:"ኣይነት" binding:"required" bson:"ኣይነት" form:"ኣይነት"`
	Address      string             `json:"አድራሻ" binding:"required" bson:"አድራሻ" form:"አድራሻ" validate:"required,min=3,max=25"`
	KifleKetema  string             `json:"ክፍለ_ከተማ" binding:"required" bson:"ክፍለ_ከተማ" form:"ክፍለ_ከተማ" validate:"required,min=3,max=25"`
	TownName     string             `json:"የሰፈር_ስም" binding:"required" bson:"የሰፈር_ስም" form:"የሰፈር_ስም" validate:"required,min=3,max=25"`
	SpecialName  string             `json:"ልዩ_ስም" binding:"required" bson:"ልዩ_ስም" form:"ልዩ_ስም" validate:"required,min=3,max=25"`
	Wereda       string             `json:"ወረዳ" binding:"required" bson:"ወረዳ" form:"ወረዳ" validate:"required,min=3,max=25"`
	HouseNumber  int32              `json:"ቤት_ቁጥር" binding:"required" bson:"ቤት_ቁጥር" form:"ቤት_ቁጥር" validate:"required"`
	Photo        string             `json:"ፎቶ" binding:"required" bson:"ፎቶ" form:"ፎቶ" validate:"required,startswith=https://"`
	LicensePhoto string             `json:"ላይሰንስ_ፎቶ" binding:"required" bson:"ላይሰንስ_ፎቶ" form:"ላይሰንስ_ፎቶ" validate:"required,startswith=https://"`
	Password     string             `json:"የሚስጥር_ቁጥር" binding:"required" bson:"የሚስጥር_ቁጥር" form:"የሚስጥር_ቁጥር" validate:"required,password"`
	PhoneNumber  string             `json:"ስልክ_ቁጥር" binding:"required" bson:"ስልክ_ቁጥር" form:"ስልክ_ቁጥር" validate:"required,phoneNumber"`
	Disabled     bool               `json:"disabled" binding:"omitempty" bson:"disabled" form:"disabled" `
	Date         primitive.DateTime `json:"date" binding:"omitempty" bson:"date,omitempty"`
}

type LoginUser struct {
	Password    string `json:"የሚስጥር_ቁጥር" binding:"required" bson:"የሚስጥር_ቁጥር" form:"የሚስጥር_ቁጥር" validate:"required,password"`
	PhoneNumber string `json:"ስልክ_ቁጥር" binding:"required" bson:"ስልክ_ቁጥር" form:"ስልክ_ቁጥር" validate:"required,phoneNumber"`
}
