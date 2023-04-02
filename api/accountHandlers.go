package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/types"
)

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request, u int) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.Store.GetAccountById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request, u int) error {
	accounts, err := s.Store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request, u int) error {
	createAccountRequest := new(types.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountRequest); err != nil {
		return err
	}
	account, err := types.NewAccount(
		createAccountRequest.Email,
		createAccountRequest.FirstName,
		createAccountRequest.LastName,
		createAccountRequest.Password,
	)

	if err != nil {
		return err
	}

	if err = s.Store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccountById(w http.ResponseWriter, r *http.Request, u int) error {
	id, err := getID(r)
	if err != nil {
		return err
	}
	if id == u {
		return fmt.Errorf("Logged in account cannot be deleted")
	}

	err = s.Store.DeleteAccount(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, nil)
}
