package utils

import (
	"api-gateway/pkg/consts"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte(consts.JWTSecret)

type UserClaims struct {
	jwt.StandardClaims
	UserId       string
	EmailAddress string
}

func GenerateToken(userId, email string) (string, error) {
	expireTime := time.Now().Add(time.Hour * consts.TokenExpiration)
	claims := UserClaims{
		UserId:       userId,
		EmailAddress: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    consts.JWTIssuer,
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

func ParseToken(token string) (*UserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !tokenClaims.Valid {
		return nil, err
	}

	claims, ok := tokenClaims.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("failed to parse claims")
	}
	return claims, nil
}
