package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	validatedAccount, err := GetAccountAndValidate(r)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	account, err := s.store.GetAccountByUsername(validatedAccount.Username)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, "Account not found")
		return
	}

	WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	accessToken := strings.Split(token, " ")[1]
	_, err := ValidateToken(accessToken)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	accounts, err := s.store.GetAccounts()
	if err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusNotFound, "Account not found")
		return
	}

	WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	validatedAccount, err := GetAccountAndValidate(r)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := s.store.DeleteAccount(validatedAccount.ID); err != nil {
		WriteJSON(w, http.StatusNotFound, "Account not found")
		return
	}

	WriteJSON(w, http.StatusOK, "Account deleted")
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var request CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err)
	}

	if request.FirstName == "" || request.LastName == "" || request.Username == "" || request.Password == "" || request.Roles == nil {
		WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	hashedAndSaltedPassword, err := HashAndSaltPassword(request.Password)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error hash")
		return
	}

	account := NewAccount(request.FirstName, request.LastName, request.Username, hashedAndSaltedPassword, request.Roles)

	accountID, err := s.store.CreateAccount(account)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error create")
		return
	}

	account.ID = accountID

	token, err := GenerateJWTToken(account)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error generate")
		return
	}

	WriteJSON(w, http.StatusOK, token)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err)
	}

	if request.Username == "" || request.Password == "" {
		WriteJSON(w, http.StatusBadRequest, "Access denied")
		return
	}

	account, err := s.store.GetAccountByUsername(request.Username)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, "Access denied")
		return
	}

	isValid := ComparePasswords(account.Password, request.Password)

	if !isValid {
		WriteJSON(w, http.StatusUnauthorized, "Access denied")
		return
	}

	token, err := GenerateJWTToken(account)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	WriteJSON(w, http.StatusOK, token)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	var request UpdateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if request.FirstName == "" && request.LastName == "" && request.Balance == "" {
		WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	validatedAccount, err := GetAccountAndValidate(r)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	account, err := s.store.GetAccountByUsername(validatedAccount.Username)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, "Account not found")
		return
	}

	var firstName string
	var lastName string
	var balance int
	if request.FirstName == "" {
		firstName = account.FirstName
	} else {
		firstName = request.FirstName
	}
	if request.LastName == "" {
		lastName = account.LastName
	} else {
		lastName = request.LastName
	}
	if request.Balance == "" {
		balance = account.Balance
	} else {
		bal, err := strconv.Atoi(request.Balance)
		if err != nil {
			fmt.Println(err)
		}
		balance = bal
	}

	account.FirstName = firstName
	account.LastName = lastName
	account.Balance = balance

	if err := s.store.UpdateAccount(account); err != nil {
		fmt.Println(err)
	}

	WriteJSON(w, http.StatusOK, "Account updated")
}

func (s *APIServer) handleTransferMoney(w http.ResponseWriter, r *http.Request) {
	var request TransferMoneyRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if request.ToAccountID == 0 || request.Amount == 0 {
		WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	validatedAccount, err := GetAccountAndValidate(r)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if validatedAccount.ID == request.ToAccountID {
		WriteJSON(w, http.StatusBadRequest, "Cannot transfer money to the same account")
		return
	}

	if err := s.store.TransferMoney(validatedAccount.ID, request.ToAccountID, request.Amount); err != nil {
		WriteJSON(w, http.StatusBadRequest, "Cannot transfer money")
		return
	}

	WriteJSON(w, http.StatusOK, "Money transferred")
}
