package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "123qweE"

func GenerateToken(userID int64, username, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"User_ID" : userID,
		"Username": username,
		"Role":     role, 
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	log.Println("Generated Token Claims:", token.Claims) 
	return tokenString, nil
}


func VerifyToken(token string)  (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Unexpected signing method.")
		}

		return []byte(secretKey), nil
	} )

	if err != nil {
		return nil, errors.New("Could not parse token.")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return nil, errors.New("Invalid token.")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	fmt.Println("Claims:", claims)
	if !ok {
		return nil, errors.New("Invalid token claims.")
	}

	return claims, nil
}