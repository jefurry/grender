package main

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var (
	JWTSecretKey string = ""
)

func getKey() string {
	if JWTSecretKey != "" {
		return JWTSecretKey
	}

	key := os.Getenv("GRENDER_JWT_SECRET_KEY")
	if key != "" {
		JWTSecretKey = key
		return key
	}

	return JWTSecretKey
}

func GenToken(exp int64, id, iss, sub, aud string) (string, error) {
	key := getKey()

	t := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        id,
		NotBefore: t - exp,
		ExpiresAt: t + exp,
		IssuedAt:  t,
		Issuer:    iss,
		Audience:  aud,
		Subject:   sub,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	key := getKey()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if token == nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
