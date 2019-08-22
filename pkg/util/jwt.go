package util

import (
	"fmt"
	"github.com/selinplus/go-dingtalk/pkg/setting"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret []byte

type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(330 * 24 * time.Hour)

	claims := CustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "daerdo-medicine",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (interface{}, error) {
	secret := []byte(setting.AppSetting.JwtSecret)
	log.Println(setting.AppSetting.JwtSecret)
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected siging method:%v", token.Header["alg"])
		}
		return secret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(jwt.MapClaims); ok && tokenClaims.Valid {
			for k, v := range claims {
				log.Printf("key is : %v, val is %v \n", k, v)
			}
			//log.Printf("username is %v, Expired at %v", claims["username"], claims["expires_at"])
			return claims, nil
		}
	}

	return nil, err
}
