package storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"example.com/types"

	_ "github.com/lib/pq"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountById(int) (*types.Account, error)
	GetAccountByNumber(int64) (*types.Account, error)
	GetAccountByEmail(string) (*types.Account, error)
	TransferRequest(*types.TransferRequest, int) error
	SeedDatabase() (*types.Account, error)
	SeedNewUser() (*types.Account, error)
}

type DatabaseStore struct {
	db       *sql.DB
	isSqlite bool
}

var fields = []string{
	"email",
	"first_name",
	"last_name",
	"account_number",
	"encrypted_password",
	"balance",
	"created_at",
}

func NewDatabaseStore(isSqlite bool) (*DatabaseStore, error) {
	if isSqlite {
		return NewSqliteStore()
	}

	return NewPostgresStore()
}

func NewPostgresStore() (*DatabaseStore, error) {
	connStr := "user=postgres dbname=postgres password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DatabaseStore{
		db:       db,
		isSqlite: false,
	}, nil
}

func NewSqliteStore() (*DatabaseStore, error) {

	db, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		return nil, err
	}

	return &DatabaseStore{
		db:       db,
		isSqlite: true,
	}, nil
}

func (s *DatabaseStore) Init() error {
	var migrationFile string
	if s.isSqlite {
		migrationFile = "./migration/sqlite/create.sql"
	} else {
		migrationFile = "./migration/postgres/create.sql"
	}

	return s.runMigration(migrationFile)
}

func (s *DatabaseStore) ReInit() error {
	err := s.dropAccountTable()
	if err != nil {
		return err
	}

	return s.Init()
}

func (s *DatabaseStore) dropAccountTable() error {
	query := "DROP TABLE IF EXISTS account"

	_, err := s.db.Exec(query)
	return err
}

func (s *DatabaseStore) CreateAccount(acc *types.Account) error {
	sql := "INSERT INTO account (%s) values ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
	query := fmt.Sprintf(sql, strings.Join(fields[:], ","))
	id := 0
	err := s.db.QueryRow(
		query,
		acc.Email,
		acc.FirstName,
		acc.LastName,
		acc.AccountNumber,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	).Scan(&id)

	if err != nil {
		return err
	}

	acc.ID = id

	return nil
}

func (s *DatabaseStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM account where id = $1", id)

	return err
}

func (s *DatabaseStore) UpdateAccount(*types.Account) error {
	return nil
}

func (s *DatabaseStore) GetAccountById(id int) (*types.Account, error) {
	sql := fmt.Sprintf("SELECT id,%s FROM account WHERE id = $1", strings.Join(fields[:], ","))
	rows, err := s.db.Query(sql, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Bank account %d not exists", id)
}

func (s *DatabaseStore) GetAccountByNumber(accountNumber int64) (*types.Account, error) {
	sql := fmt.Sprintf("SELECT id,%s FROM account WHERE account_number = $1", strings.Join(fields[:], ","))
	rows, err := s.db.Query(sql, accountNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("Bank account number '%d' not exists", accountNumber)
}

func (s *DatabaseStore) GetAccountByEmail(email string) (*types.Account, error) {
	sql := fmt.Sprintf("SELECT id,%s FROM account WHERE email = $1", strings.Join(fields[:], ","))
	rows, err := s.db.Query(sql, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf(" with email '%s' not exists", email)
}

func (s *DatabaseStore) withdraw(id, amount int) error {
	account, err := s.GetAccountById(id)
	if err != nil {
		return err
	}

	if account.Balance-int64(amount) < 0 {
		return fmt.Errorf("Ballance would go negative")
	}

	_, err = s.db.Exec("UPDATE account set balance = balance - $1 where id = $2", amount, id)

	return err
}

func (s *DatabaseStore) credit(accountNumber, amount int) error {
	_, err := s.db.Exec("UPDATE account set balance = balance + $1 where account_number = $2", amount, accountNumber)

	return err
}

func (s *DatabaseStore) GetAccounts() ([]*types.Account, error) {
	accounts := []*types.Account{}
	sql := fmt.Sprintf("SELECT id,%s FROM account", strings.Join(fields[:], ","))
	rows, err := s.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	if err := rows.Scan(
		&account.ID,
		&account.Email,
		&account.FirstName,
		&account.LastName,
		&account.AccountNumber,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *DatabaseStore) SeedDatabase() (*types.Account, error) {
	err := s.ReInit()
	if err != nil {
		return nil, err
	}

	return s.SeedNewUser()
}

func (s *DatabaseStore) SeedNewUser() (*types.Account, error) {
	account, err := types.NewAccount("testemail@email.com", "Test", "User", "boom")
	account.Balance = 5000
	if err != nil {
		return nil, err
	}
	err = s.CreateAccount(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *DatabaseStore) runMigration(fileName string) error {
	sql, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(string(sql))

	return err
}
