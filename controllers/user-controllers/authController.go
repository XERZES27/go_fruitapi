package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	helper "github.com/XERZES27/go_fruitapi/helpers"
	model "github.com/XERZES27/go_fruitapi/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 7 {
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

func validatePhoneNumber(f1 validator.FieldLevel) bool {
	phoneNumber := f1.Field().String()
	isValid := true
	switch len(phoneNumber) {
	case 13:
		if strings.HasPrefix(phoneNumber, "+2519") {
			_, err := strconv.Atoi(phoneNumber[5:])
			isValid = err == nil
		}
	case 10:
		if strings.HasPrefix(phoneNumber, "09") {
			_, err := strconv.Atoi(phoneNumber[2:])
			isValid = err == nil
		}

	default:
		isValid = false
	}
	return isValid
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func formatPhoneNumber(phoneNumber string) string {
	if len(phoneNumber) == 13 {
		return "0" + phoneNumber[4:]
	}
	return phoneNumber
}

func Register(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		e1 := c.ShouldBind(&user)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": e1.Error()})
			return
		}

		validate := validator.New()
		validate.RegisterValidation("password", validatePassword)
		validate.RegisterValidation("phoneNumber", validatePhoneNumber)
		user.PhoneNumber = formatPhoneNumber(user.PhoneNumber)
		user.Date = primitive.NewDateTimeFromTime(time.Now().UTC())
		user.Disabled = false
		e2 := validate.Struct(user)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema", "reason": e2.Error()})
			return
		}
		var result bson.M
		e3 := userCollection.FindOne(context.TODO(), bson.M{"ስልክ_ቁጥር": user.PhoneNumber}).Decode(&result)
		if e3 != nil {
			if e3 != mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": e3.Error()})
				return
			}
		}
		if len(result) > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Document Already Exists"})
			return
		}
		hashedPassword, e4 := HashPassword(user.Password)
		if e4 != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "internal error"})
			return
		}
		user.Password = hashedPassword

		basicUserResult, err := userCollection.InsertOne(context.TODO(), user)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"ID": basicUserResult.InsertedID})

	}
}

func Login(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.LoginUser
		e1 := c.ShouldBind(&user)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Data", "reason": e1.Error()})
			return
		}
		validate := validator.New()
		validate.RegisterValidation("password", validatePassword)
		validate.RegisterValidation("phoneNumber", validatePhoneNumber)
		user.PhoneNumber = formatPhoneNumber(user.PhoneNumber)
		e2 := validate.Struct(user)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Schema"})
			return
		}
		var result model.User
		e3 := userCollection.FindOne(context.TODO(), bson.M{"ስልክ_ቁጥር": user.PhoneNumber}).Decode(&result)
		if e3 != nil {
			if e3 == mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Email Password Combination"})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"status": "internal error"})
				return
			}
		}

		passwordMatches := CheckPasswordHash(user.Password, fmt.Sprintf("%v", result.Password))
		if !passwordMatches {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Email Password Combination"})
			return
		}

		t := jwt.New(jwt.GetSigningMethod("RS256"))

		t.Claims = &helper.CustomClaims{
			Issuer:      "Server",
			IssuedAt:    time.Now().Unix(),
			Id:          result.ID.Hex(),
			CompanyName: fmt.Sprintf(result.CompanyName),
			PhoneNumber: fmt.Sprintf(result.PhoneNumber)}

		var pubKey []byte
		var privKey []byte
		keyErr := helper.GetKeys(&privKey, &pubKey, "user")
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
