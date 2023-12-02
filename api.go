package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", s.handleGetAccount).Methods(http.MethodGet)
	router.HandleFunc("/account", s.handleUpdateAccount).Methods(http.MethodPatch)
	router.HandleFunc("/account", s.handleCreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/account", s.handleDeleteAccount).Methods(http.MethodDelete)
	router.HandleFunc("/accounts", s.handleGetAccounts).Methods(http.MethodGet)
	router.HandleFunc("/transfer", s.handleTransferMoney).Methods(http.MethodPost)
	router.HandleFunc("/login", s.handleLogin).Methods(http.MethodPost)

	fmt.Println("Server is running 🚀")
	http.ListenAndServe(s.listenAddress, router)
}
