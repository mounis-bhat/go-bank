package main

import (
	"fmt"
	"math/rand"
)

type Account struct {
	ID            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Balance       int    `json:"balance"`
	AccountNumber string `json:"account_number"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:            rand.Intn(100),
		FirstName:     firstName,
		LastName:      lastName,
		Balance:       1000,
		AccountNumber: fmt.Sprintf("%d-%d-%d", rand.Intn(100), rand.Intn(100), rand.Intn(100)),
	}
}
