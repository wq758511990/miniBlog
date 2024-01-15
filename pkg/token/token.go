package token

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Config struct {
	key         string
	identityKey string
}

var ErrMissingHeader = errors.New("the length of the `Authorization` header is zero")

var (
	config = Config{viper.GetString("jwt-secret"), "identityKey"}
	once   sync.Once
)

func Init(key string, identityKey string) {
	once.Do(func() {
		if key != "" {
			config.key = key
		}
		if identityKey != "" {
			config.identityKey = identityKey
		}
	})
}

func Parse(tokenString string, key string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	var identityKey string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		identityKey = claims[config.identityKey].(string)
	}
	return identityKey, nil
}

func Sign(identityKey string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		config.identityKey: identityKey,
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"exp":              time.Now().Add(100000 * time.Hour).Unix(),
	})
	tokenString, err = token.SignedString([]byte(config.key))

	return
}

func ParseRequest(c *gin.Context) (string, error) {
	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		return "", ErrMissingHeader
	}
	var t string
	// 从请求头中取出 token
	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, config.key)
}
