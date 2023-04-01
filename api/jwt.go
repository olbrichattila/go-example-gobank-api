package api

import (
	"fmt"
	"os"

	"example.com/types"

	"github.com/golang-jwt/jwt/v5"
)

func createJWT(account *types.Account) (string, error) {
	hmacSampleSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiresAt": 15000,
		"accountId": account.ID,
	})

	tokenString, err := token.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signin method %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
