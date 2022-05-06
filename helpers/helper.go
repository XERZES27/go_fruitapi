package helpers

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	publicKey  []byte
	privateKey []byte
)

func GetKeys(privKey *[]byte, pubKey *[]byte, accountType string) error {
	if len(privateKey) == 0 && len(publicKey) == 0 {
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
		claims := &CustomClaims{}

		tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {

			return pub, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
			}
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid token"})

			c.Abort()
			return
		}
		if !tkn.Valid {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}
		if accountType == "user" {
			c.Set("CompanyName", claims.CompanyName)
			c.Set("PhoneNumber", claims.PhoneNumber)
			c.Set("Id", claims.Id)
		}
		c.Next()
	}
}
