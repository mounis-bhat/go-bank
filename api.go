package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handCreateAccount))

	fmt.Println("The API is listening on port", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handGetAccount(w, r)
	case http.MethodDelete:
		return s.handDeleteAccount(w, r)
	default:
		return fmt.Errorf("unknown method %s", r.Method)
	}

}

func (s *APIServer) handGetAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("John", "Doe")
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println((id))
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	listenAddress string
}

func NewAPIServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{
				Error: err.Error(),
			})
		}
	}
}
