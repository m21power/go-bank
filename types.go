package main

import "math/rand"

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Number    int64  `json:"number"`
	balance   int64  `json:"balance"`
}

func NewAccount(firstname, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(100000),
		FirstName: firstname,
		LastName:  lastname,
		Number:    int64(rand.Intn(1000000)),
	}

}
