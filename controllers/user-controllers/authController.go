package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	publicKey  []byte
	privateKey []byte
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getKeys(privKey *[]byte, pubKey *[]byte) error {
	if len(privateKey) == 0 && len(publicKey) == 0 {
		wd, err := os.Getwd()
		if err != nil {

			return err
		}
		priv, err := os.ReadFile(wd + "/keys/privkey.pem")
		if err != nil {

			return err
		}
		pub, err := os.ReadFile(wd + "/keys/pubkey.pem")
		if err != nil {
			return err
		}
		privateKey = priv
		publicKey = pub

		*privKey = privateKey
		*pubKey = publicKey
		return nil
	} else {

		*privKey = privateKey
		*pubKey = publicKey
		return nil
	}

}

func Register(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		e1 := c.ShouldBind(&user)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query parameter", "reason": e1.Error()})
			return
		}
		validate := validator.New()
		validate.RegisterValidation("password", validatePassword)

		e2 := validate.Struct(user)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Schema", "reason": e2.Error()})
			return
		}
		var result bson.M
		e3 := userCollection.FindOne(context.TODO(), bson.M{"name": user.Name}).Decode(&result)
		if e3 != nil {
			if e3 != mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": e3.Error()})
				return
			}
		}
		if len(result) > 0 {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Document Already Exists"})
			return
		}
		hashedPassword, e4 := HashPassword(user.Password)
		if e4 != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}
		user.Password = hashedPassword

		basicUserResult, err := userCollection.InsertOne(context.TODO(), user)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, gin.H{"ID": basicUserResult.InsertedID})

	}
}

func Login(userCollection *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var basicUser LoginBasicUser
		e1 := c.ShouldBind(&basicUser)
		if e1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing query parameter", "reason": e1.Error()})
			return
		}
		validate := validator.New()
		validate.RegisterValidation("password", validatePassword)
		e2 := validate.Struct(basicUser)
		if e2 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Schema"})
			return
		}
		var result bson.M
		e3 := userCollection.FindOne(context.TODO(), bson.M{"email": basicUser.Email}).Decode(&result)
		if e3 != nil {
			if e3 == mongo.ErrNoDocuments {
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Email Password Combination"})
				return
			} else {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
				return
			}
		}

		passwordMatches := CheckPasswordHash(basicUser.Password, fmt.Sprintf("%v", result["password"]))
		if !passwordMatches {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Email Password Combination"})
			return
		}

		t := jwt.New(jwt.GetSigningMethod("RS256"))

		t.Claims = jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Issuer:    "Server",
			Id:        fmt.Sprint(result["_id"]),
		}
		var pubKey []byte
		var privKey []byte
		keyErr := getKeys(&privKey, &pubKey)
		if keyErr != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}
		signKey, err1 := jwt.ParseRSAPrivateKeyFromPEM(privKey)
		if err1 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err1.Error()})
			return
		}

		re, e5 := t.SignedString(signKey)
		if e5 != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": e5.Error()})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"data": re})
	}
}
