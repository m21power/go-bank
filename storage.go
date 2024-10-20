package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(string) error
	UpdateAccount(*Account) (*Account, error)
	GetAccount() ([]*Account, error)
	GetAccountByID(string) (*Account, error)
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
func (s *PostgresStore) DeleteAccount(id string) error {
	_, err := s.db.Exec("delete from account where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgresStore) UpdateAccount(acc *Account) (*Account, error) {
	_, err := s.db.Exec("update account set firstname=$1,lastname=$2,balance=$3,create_at=$4 where id=$5", acc.FirstName, acc.LastName, acc.Balance, acc.CreatedAt, acc.ID)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
func (s *PostgresStore) GetAccountByID(id string) (*Account, error) {
	row := s.db.QueryRow("select * from account where id=$1", id)
	acc := &Account{}
	err := row.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Balance, &acc.CreatedAt)
	if err == sql.ErrNoRows {
		return NewAccount("", ""), nil
	} else if err != nil {
		return nil, err
	}
	return acc, nil
}
func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()

}
func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
	id serial primary key,firstname varchar(50),lastname varchar(50),balance serial,create_at timestamp)`
	_, err := s.db.Exec(query)
	return err
}
