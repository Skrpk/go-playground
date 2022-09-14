package middlewares

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/big"
	"net/http"
)

type JWK struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

func MakeJwtVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		publicKeyPath := "public_key.json"
		token := c.GetHeader("Authorization")

		isValid, err := verifyToken(token, publicKeyPath)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if !isValid {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func verifyToken(tokenString, publicKeyPath string) (bool, error) {
	keyData, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return false, err
	}

	jwk := new(JWK)
	err = json.Unmarshal(keyData, jwk)
	if err != nil {
		return false, err
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key := convertKey(jwk.Keys[1].E, jwk.Keys[1].N)
		return key, nil
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func convertKey(rawE, rawN string) *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey
}
