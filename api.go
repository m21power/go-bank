package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
	router.HandleFunc("/account/delete/{id}", makeHTTPHandleFunc(s.handleDeleteAccount))
	router.HandleFunc("/account/update/{id}", makeHTTPHandleFunc(s.handleUpdateAccount))

	log.Println("JSON Api running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)

}

type ApiError struct {
	Error string
}

type DeleteSuccess struct {
	Message string
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

//The "s *APIServer" part means the method has a receiver s, which is a pointer to an APIServer.
// This means the method can modify the server's fields if needed.

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)

}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)["id"] // we are getting id from the request
	// accessing account from the db
	account, err := s.store.GetAccountByID(vars)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account, err := s.store.GetAccount()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAcccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAcccountReq); err != nil {
		return err
	}
	account := NewAccount(createAcccountReq.FirstName, createAcccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	Id := mux.Vars(r)["id"]
	s.store.DeleteAccount(Id)
	return WriteJSON(w, http.StatusOK, DeleteSuccess{Message: "Account Deleted"})
}
func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	account, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}
	updateAccountReq := new(UpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(updateAccountReq); err != nil {
		return err
	}
	account.FirstName = updateAccountReq.FirstName
	account.LastName = updateAccountReq.LastName
	account.Balance = updateAccountReq.Balance
	account.CreatedAt = updateAccountReq.CreatedAt
	account.Number = updateAccountReq.Number
	account, err = s.store.UpdateAccount(account)
	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apifunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apifunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error here
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
