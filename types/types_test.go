package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account := NewAccount("Mounis", "Bhat", "mounis", "password", []string{"admin"})
	assert.Equal(t, "Mounis", account.FirstName)
	assert.Equal(t, "Bhat", account.LastName)
	assert.Equal(t, "mounis", account.Username)
	assert.Equal(t, "password", account.Password)
	assert.Equal(t, []string{"admin"}, account.Roles)
}
