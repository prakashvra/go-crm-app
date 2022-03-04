package repository

import (
	"crm-app/src/model"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// AccountRepository is called by AccountService for db operations
type AccountRepository interface {
	CreateNewAccount(Account model.AccountRequest) (int64, string)
	UpdateAccount(Account model.AccountRequest) error
	GetAccounts(active bool) ([]model.Account, error)
	GetAccountById(id int64) (model.Account, error)
	GetAccountByEmail(email string) (model.Account, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (accountRepo *accountRepository) CreateNewAccount(request model.AccountRequest) (int64, string) {
	addAccountStat := `
		INSERT INTO "account" (type,first_name,last_name,email,phone,gender,status,has_orders) 
		VALUES ($1, $2, $3, $4, $5, $6,$7) RETURNING id
		`
	log.Println(accountRepo.db)
	row := accountRepo.db.QueryRow(addAccountStat, request.FirstName, request.LastName, request.Email, request.Phone, request.Gender, request.Status, request.HasOrders)
	err := row.Err().Error()
	var id int64
	if err != "" {
		log.Printf("Error while creating account in DB: %v", err)
		if strings.Contains(err, "unique constraint \"account_phone_key\"") {
			err = "Account with the given phone already exists"
		}
		return id, err
	}
	err = row.Scan(&id).Error()
	if err != "" {
		log.Printf("Error while creating account in DB: %v", err)
		return id, err
	}
	return id, ""
}

func (accountRepo *accountRepository) UpdateAccount(request model.AccountRequest) error {
	updateAccountStat := `
		UPDATE "account" 
		SET first_name = $1,last_name = $2,phone = $3,gender = $4,active = $5
		WHERE id = $6
		`
	err := accountRepo.db.QueryRow(updateAccountStat, request.FirstName, request.LastName, request.Phone, request.Gender, request.Active, request.Id).Err()

	if err != nil {
		log.Printf("Error while updating account in DB: %v", err.Error())
		return err
	}

	return nil
}

func (accountRepo *accountRepository) GetAccounts(active bool) ([]model.Account, error) {
	fmt.Printf("GetAccounts active : %v\n", active)
	accountQuery := `
		SELECT id,first_name,last_name,gender,email,phone,source,status,active,account_id FROM "account" 
		WHERE active = $1
		`
	rows, err := accountRepo.db.Query(accountQuery, active)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An Account slice to hold data from returned rows.
	var accounts []model.Account

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var acct model.Account
		var source *string
		if err := rows.Scan(&acct.Id, &acct.FirstName, &acct.LastName, &acct.Gender, &acct.Email, &acct.Phone, &source,
			&acct.Status, &acct.Active, &acct); err != nil {
			if source != nil {
				acct.Source = *source
			}
			return accounts, err
		}
		accounts = append(accounts, acct)
	}
	if err = rows.Err(); err != nil {
		return accounts, err
	}
	return accounts, nil

}

func (accountRepo *accountRepository) GetAccountById(id int64) (model.Account, error) {
	accountQuery := `
		SELECT id,first_name,last_name,gender,email,phone,source,status,active,account_id FROM "account" 
		WHERE id = ?
		`
	row := accountRepo.db.QueryRow(accountQuery, id)
	return accountRepo.scanResult(row)
}

func (accountRepo *accountRepository) GetAccountByEmail(email string) (model.Account, error) {
	accountQuery := `
		SELECT id,first_name,last_name,gender,email,phone,source,status,active,account_id FROM "account" 
		WHERE email = ?
		`
	row := accountRepo.db.QueryRow(accountQuery, email)
	return accountRepo.scanResult(row)
}

func (accountRepo *accountRepository) UpdateAccountStatus(id int, status string) error {
	updateAccountStat := `
		UPDATE "account" 
		SET status = $1
		WHERE id = $2
		`
	err := accountRepo.db.QueryRow(updateAccountStat, status, id).Err()

	if err != nil {
		log.Printf("Error while updating account status in DB: %v", err.Error())
		return err
	}

	return nil
}

func (accountRepo *accountRepository) scanResult(row *sql.Row) (model.Account, error) {
	// A Account object to hold data from returned row.
	var acct model.Account

	if row == nil {
		return acct, nil
	}

	// Loop through rows, using Scan to assign column data to struct fields.
	if err := row.Scan(&acct.Id, &acct.FirstName, &acct.LastName, &acct.Gender, &acct.Email, &acct.Phone, &acct.Source,
		&acct.Status, &acct.Active); err != nil {
		return acct, err
	}
	return acct, nil
}
