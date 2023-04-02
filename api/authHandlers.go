package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/types"
	"golang.org/x/crypto/bcrypt"
)

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	account, err := s.Store.GetAccountByEmail(req.Email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.EncryptedPassword), []byte(req.Password))
	if err != nil {
		return fmt.Errorf("unauthenticated")
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	resp := &types.LoginRespose{
		Token: tokenString,
	}

	return WriteJSON(w, http.StatusOK, resp)
}
