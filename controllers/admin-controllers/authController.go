package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"unicode"

	helper "github.com/XERZES27/go_fruitapi/helpers"
	model "github.com/XERZES27/go_fruitapi/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func validateAdminPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 25 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(adminCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var admin model.Admin
		e1 := c.ShouldBind(&admin)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": e1.Error()})
			return
		}

		validate := validator.New()
		validate.RegisterValidation("password", validateAdminPassword)
		admin.Date = primitive.NewDateTimeFromTime(time.Now().UTC())
		e2 := validate.Struct(admin)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema", "reason": e2.Error()})
			return
		}
		var result bson.M
		e3 := adminCollection.FindOne(context.TODO(), bson.M{"name": "ADMIN_ZERO"}).Decode(&result)
		if e3 != nil {
			if e3 != mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": e3.Error()})
				return
			}
		}
		if len(result) > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Admin already exists"})
			return
		}
		hashedPassword, e4 := HashPassword(admin.Password)
		if e4 != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "internal error"})
			return
		}
		admin.Password = hashedPassword

		adminResult, err := adminCollection.InsertOne(context.TODO(), admin)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"ID": adminResult.InsertedID})

	}
}

func Login(adminCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var admin model.LoginAdmin
		e1 := c.ShouldBind(&admin)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": e1.Error()})
			return
		}
		validate := validator.New()
		err := validate.RegisterValidation("password", validateAdminPassword)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": err.Error()})
			return
		}
		e2 := validate.Struct(admin)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema"})
			return
		}
		var result model.Admin
		e3 := adminCollection.FindOne(context.TODO(), bson.M{"name": admin.Name}).Decode(&result)
		if e3 != nil {
			if e3 == mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Email Password Combination"})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "internal error"})
				return
			}
		}

		passwordMatches := CheckPasswordHash(admin.Password, fmt.Sprintf("%v", result.Password))
		if !passwordMatches {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Email Password Combination"})
			return
		}

		t := jwt.New(jwt.GetSigningMethod("RS256"))

		t.Claims = jwt.StandardClaims{
			Issuer:   "Server",
			IssuedAt: time.Now().Unix(),
			Id:       result.ID.Hex()}

		var pubKey []byte
		var privKey []byte
		keyErr := helper.GetKeys(&privKey, &pubKey, "admin")
		if keyErr != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "internal error"})
			return
		}
		signKey, err1 := jwt.ParseRSAPrivateKeyFromPEM(privKey)
		if err1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": err1.Error()})
			return
		}

		re, e5 := t.SignedString(signKey)
		if e5 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": e5.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"id": result.ID, "X-Access-Token": re})
	}
}
