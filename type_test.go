package main

import (
	"fmt"
	"testing"

	"example.com/types"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	res, err := types.NewAccount("email@email.com", "a", "b", "bumm")

	assert.Nil(t, err)

	fmt.Println(res.EncryptedPassword)
}
