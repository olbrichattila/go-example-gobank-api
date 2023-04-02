package main

import (
	"testing"

	"example.com/types"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	_, err := types.NewAccount("email@email.com", "a", "b", "bumm")

	assert.Nil(t, err)
}
