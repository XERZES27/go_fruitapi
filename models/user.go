package main

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"id" binding:"omitempty"  bson:"_id,omitempty" `
	የኩባኒኛ_ስም string             `json:"የኩባኒኛ_ስም" binding:"required" bson:"የኩባኒኛ_ስም" form:"የኩባኒኛ_ስም" validate:"required,min=6,max=255"`
	ኣይነት     [3]string          `json:"ኣይነት" binding:"required" bson:"የኩባኒኛ_ስም" form:"የኩባኒኛ_ስም"`
	አድራሻ     string             `json:"አድራሻ" binding:"required" bson:"አድራሻ" form:"አድራሻ" validate:"required,min=3,max=25"`
	ክፍለ_ከተማ  string             `json:"ክፍለ_ከተማ" binding:"required" bson:"ክፍለ_ከተማ" form:"ክፍለ_ከተማ" validate:"required,min=3,max=25"`
	የሰፈር_ስም  string             `json:"የሰፈር_ስም" binding:"required" bson:"የሰፈር_ስም" form:"የሰፈር_ስም" validate:"required,min=3,max=25"`
	ልዩ_ስም    string             `json:"ልዩ_ስም" binding:"required" bson:"ልዩ_ስም" form:"ልዩ_ስም" validate:"required,min=3,max=25"`
	ወረዳ      string             `json:"ወረዳ" binding:"required" bson:"ወረዳ" form:"ወረዳ" validate:"required,min=3,max=25"`
	ቤት_ቁጥር   int32              `json:"ቤት_ቁጥር" binding:"required" bson:"ቤት_ቁጥር" form:"ቤት_ቁጥር" validate:"required"`
	ፎቶ       string             `json:"ፎቶ" binding:"required" bson:"ፎቶ" form:"ፎቶ" validate:"required,startswith="https://"`
	ላይሰንስ_ፎቶ string             `json:"ላይሰንስ_ፎቶ" binding:"required" bson:"ላይሰንስ_ፎቶ" form:"ላይሰንስ_ፎቶ" validate:"required,startswith="https://"`
}

func Print(){
	fmt.Println("here")
}
