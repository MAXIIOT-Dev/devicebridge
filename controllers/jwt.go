package controllers

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/maxiiot/vbaseBridge/storage"
)

var (
	TokenInvalid error  = errors.New("Couldn't handle this token:")
	jwtsecret    []byte = []byte("maxiiot123")
	jwtTTL              = time.Hour * 24
)

// SetJWTSecret set jwt secret
func SetJWTSecret(secret string) {
	jwtsecret = []byte(secret)
}

// CreateToken generate jwt
func CreateToken(user storage.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "maxiiot",
		"username": user.UserName,
		"exp":      time.Now().Add(jwtTTL).Unix(),
	})

	return token.SignedString(jwtsecret)
}

// ParseToken parse jwt
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtsecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// JWTAuth jwt uath middleware
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			Response(c, http.StatusUnauthorized, 0, 0, "请求未携带token，无权限访问", nil)
			c.Abort()
			return
		}

		claims, err := ParseToken(token)
		if err != nil {
			Response(c, http.StatusUnauthorized, 0, 0, "token过期", nil)
			c.Abort()
			return
		}

		if username, ok := claims["username"]; ok {
			c.Set("username", username)
		}
		c.Next()
	}
}
