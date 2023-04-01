package types

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	AccountNumber     int64     `json:"accountNumber"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"-"`
}

func NewAccount(email, firstName, lastName, password string) (*Account, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		Email:             email,
		FirstName:         firstName,
		LastName:          lastName,
		AccountNumber:     int64(rand.Intn(10000000)),
		EncryptedPassword: string(encryptedPassword),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
