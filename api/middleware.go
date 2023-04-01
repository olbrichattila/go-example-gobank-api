package api

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func withJWTAuth(hendlerFunc apiMiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			hendlerFunc(w, r, 0)
			return
		}
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		claim := token.Claims.(jwt.MapClaims)
		accountId := int(claim["accountId"].(float64))

		hendlerFunc(w, r, accountId)
	}
}
