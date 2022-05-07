package helpers

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	publicUserKey   []byte
	privateUserKey  []byte
	publicAdminKey  []byte
	privateAdminKey []byte
)

func GetKeys(privKey *[]byte, pubKey *[]byte, accountType string) error {
	var privateKey *[]byte
	var publicKey *[]byte
	if accountType == "user" {
		privateKey = &publicUserKey
		publicKey = &privateUserKey

	} else {
		privateKey = &publicAdminKey
		publicKey = &privateAdminKey

	}
	if len(*publicKey) == 0 && len(*privateKey) == 0 {
		fmt.Println("from disk")
		wd, err := os.Getwd()
		if err != nil {

			return err
		}
		priv, err := os.ReadFile(wd + "/keys/" + accountType + "/privkey.pem")
		if err != nil {

			return err
		}
		pub, err := os.ReadFile(wd + "/keys/" + accountType + "/pubkey.pem")
		if err != nil {
			return err
		}
		if accountType == "user" {
			publicUserKey = priv
			privateUserKey = pub

		} else {
			publicAdminKey = priv
			privateAdminKey = pub

		}
		*privKey = priv
		*pubKey = pub

		return nil
	} else {
		fmt.Println("from memory")
		*privKey = *privateKey
		*pubKey = *publicKey
		return nil
	}

}

func VerifyToken(accountType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pubKey []byte
		var privKey []byte
		keyErr := GetKeys(&privKey, &pubKey, accountType)
		if keyErr != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			c.Abort()
			return
		}

		block, _ := pem.Decode(pubKey)
		if block == nil {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "failed to parse PEM block containing the public key"})
			c.Abort()
			return
		}
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "failed to parse DER encoded public key: " + err.Error()})
			c.Abort()
			return
		}

		token := c.GetHeader("X-Access-Token")

		if accountType == "user" {
			claims := &CustomClaims{}
			tkn, err1 := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {

				return pub, nil
			})
			if err1 != nil {
				if err1 == jwt.ErrSignatureInvalid {
					c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
					c.Abort()
				}
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid token", "reason": err.Error()})

				c.Abort()
				return
			}
			if !tkn.Valid {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}
			c.Set("CompanyName", claims.CompanyName)
			c.Set("PhoneNumber", claims.PhoneNumber)
			c.Set("Id", claims.Id)

			c.Next()
		} else {
			claims := &jwt.StandardClaims{}
			tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {

				return pub, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
					c.Abort()
				}
				c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid token", "reason": err.Error()})

				c.Abort()
				return
			}
			if !tkn.Valid {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}
			c.Set("Id", claims.Id)
			c.Next()
		}

	}
}
