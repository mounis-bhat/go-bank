package main

import "time"

type CreateAccountRequest struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Roles     []string `json:"roles"`
}

func NewAccount(firstName string, lastName string, username string, password string, roles []string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Password:  password,
		Roles:     roles,
		CreatedAt: time.Now().UTC(),
	}
}

type UpdateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   string `json:"balance"`
}

type ValidateAccountRequest struct {
	ID        int    `json:"account_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   int    `json:"balance"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	IAT       int    `json:"iat"`
	EXP       int    `json:"exp"`
}

type TransferMoneyRequest struct {
	ToAccountID int `json:"to_account_id"`
	Amount      int `json:"amount"`
}

type AccountsRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	Id        int    `json:"account_id"`
	Balance   int    `json:"balance"`
}

type Account struct {
	ID        int       `json:"account_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Roles     []string  `json:"roles"`
}

type APIServer struct {
	listenAddress string
	store         Storage
}

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

type APIError struct {
	Error string `json:"error"`
}
