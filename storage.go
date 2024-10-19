package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccount() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=mesay password=mesay dbname=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `insert into account (firstname,lastname,balance,create_at) values($1,$2,$3,$4)`
	_, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	return nil

}
func (s *PostgresStore) GetAccount() ([]*Account, error) {
	rows, err := s.db.Query("select * from account")
	if err != nil {
		return nil, err
	}
	account := []*Account{}
	for rows.Next() {
		acc := &Account{}
		err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Balance, &acc.CreatedAt)
		if err != nil {
			return nil, err
		}
		account = append(account, acc)
	}
	return account, nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()

}
func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
	id serial primary key,firstname varchar(50),lastname varchar(50), balance serial,create_at timestamp)`
	_, err := s.db.Exec(query)
	return err
}
