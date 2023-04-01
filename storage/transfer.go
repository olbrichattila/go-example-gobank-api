package storage

import (
	"example.com/types"
)

func (s *DatabaseStore) TransferRequest(transferRequest *types.TransferRequest, u int) error {

	_, err := s.GetAccountByNumber(int64(transferRequest.ToAccount))
	if err != nil {
		return err
	}

	if err := s.withdraw(u, transferRequest.Amount); err != nil {
		return err
	}

	return s.credit(transferRequest.ToAccount, transferRequest.Amount)
}
