package api

import (
	"encoding/json"
	"net/http"

	"example.com/types"
)

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request, u int) error {
	transferRequest := new(types.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	err := s.Store.TransferRequest(transferRequest, u)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, transferRequest)
}
