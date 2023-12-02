package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/mounis-bhat/go-bank/pkg/lib"
	"github.com/mounis-bhat/go-bank/types"
)

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	validatedAccount, err := lib.GetAccountAndValidate(r)
	if err != nil {
		lib.WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fmt.Println(validatedAccount.Username)

	account, err := s.dbConnection.GetAccountByUsername(validatedAccount.Username)
	if err != nil {

		lib.WriteJSON(w, http.StatusNotFound, "Account not found")
		return
	}

	lib.WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	validatedAccount, err := lib.GetAccountAndValidate(r)
	if err != nil {
		lib.WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !strings.Contains(validatedAccount.Roles, "admin") {
		lib.WriteJSON(w, http.StatusForbidden, "Forbidden")
		return
	}

	accounts, err := s.dbConnection.GetAccounts()
	if err != nil {
		fmt.Println(err)
		lib.WriteJSON(w, http.StatusNotFound, "Account not found")
		return
	}

	lib.WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
	validatedAccount, err := lib.GetAccountAndValidate(r)
	if err != nil {
		lib.WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := s.dbConnection.DeleteAccount(validatedAccount.ID); err != nil {
		lib.WriteJSON(w, http.StatusNotFound, "Account not found")
		return
	}

	lib.WriteJSON(w, http.StatusOK, "Account deleted")
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	var request types.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err)
	}

	if request.FirstName == "" || request.LastName == "" || request.Username == "" || request.Password == "" || request.Roles == nil {
		lib.WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	isValid := lib.IsValidPassword(request.Password)

	if !isValid {
		lib.WriteJSON(w, http.StatusBadRequest, "The password must be at least 8 characters long, contain at least one uppercase/lowercase letter, at least one number and at least one special character")
		return
	}

	hashedAndSaltedPassword, err := lib.HashAndSaltPassword(request.Password)
	if err != nil {
		lib.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error hash")
		return
	}

	account := types.NewAccount(request.FirstName, request.LastName, request.Username, hashedAndSaltedPassword, request.Roles)

	accountID, err := s.dbConnection.CreateAccount(account)
	if err != nil {
		lib.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error create")
		return
	}

	account.ID = accountID

	token, err := lib.GenerateJWTToken(account)
	if err != nil {
		lib.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error generate")
		return
	}

	lib.WriteJSON(w, http.StatusOK, token)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var request types.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println(err)
	}

	if request.Username == "" || request.Password == "" {
		lib.WriteJSON(w, http.StatusBadRequest, "Access denied")
		return
	}

	account, err := s.dbConnection.GetAccountByUsername(request.Username)
	if err != nil {
		lib.WriteJSON(w, http.StatusNotFound, "Access denied")
		return
	}

	isValid := lib.ComparePasswords(account.Password, request.Password)

	if !isValid {
		lib.WriteJSON(w, http.StatusUnauthorized, "Access denied")
		return
	}

	token, err := lib.GenerateJWTToken(account)
	if err != nil {
		lib.WriteJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	lib.WriteJSON(w, http.StatusOK, token)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) {
	var request types.UpdateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		lib.WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if request.FirstName == "" && request.LastName == "" && request.Balance == "" {
		lib.WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	validatedAccount, err := lib.GetAccountAndValidate(r)
	if err != nil {
		lib.WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	account, err := s.dbConnection.GetAccountByUsername(validatedAccount.Username)
	if err != nil {
		lib.WriteJSON(w, http.StatusNotFound, "Account not found")
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

	if err := s.dbConnection.UpdateAccount(account); err != nil {
		fmt.Println(err)
	}

	lib.WriteJSON(w, http.StatusOK, "Account updated")
}

func (s *APIServer) handleTransferMoney(w http.ResponseWriter, r *http.Request) {
	var request types.TransferMoneyRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		lib.WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	if request.ToAccountID == 0 || request.Amount == 0 {
		lib.WriteJSON(w, http.StatusBadRequest, "Invalid body")
		return
	}

	validatedAccount, err := lib.GetAccountAndValidate(r)
	if err != nil {
		lib.WriteJSON(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if validatedAccount.ID == request.ToAccountID {
		lib.WriteJSON(w, http.StatusBadRequest, "Cannot transfer money to the same account")
		return
	}

	if err := s.dbConnection.TransferMoney(validatedAccount.ID, request.ToAccountID, request.Amount); err != nil {
		lib.WriteJSON(w, http.StatusBadRequest, "Cannot transfer money")
		return
	}

	lib.WriteJSON(w, http.StatusOK, "Money transferred")
}
