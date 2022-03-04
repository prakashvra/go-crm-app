package repository

import (
	"crm-app/src/model"
	"crm-app/src/util"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// CustomerRepository is called by CustomerService for db operations
type CustomerRepository interface {
	CreateNewCustomer(customer model.CustomerRequest) (int, string)
	UpdateCustomer(customer model.CustomerRequest) error
	GetCustomers(active bool) ([]model.Customer, error)
	GetCustomerById(id int64) (model.Customer, error)
	GetCustomerByEmail(email string) (model.Customer, error)
	UpdateCustomerStatus(id int64, status string) (model.Customer, error)
}

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) CustomerRepository {
	return &customerRepository{
		db: db,
	}
}

func (customerRepo *customerRepository) CreateNewCustomer(request model.CustomerRequest) (int, string) {
	id := 0
	//var custId *int
	tx, err := customerRepo.db.Begin()
	if err != nil {
		return id, err.Error()
	}
	{
		addCustomerStat := `
		INSERT INTO customer (first_name,last_name,email,phone,gender,status,has_orders,active,created_by) 
		VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9) RETURNING id
		`
		stmt, err := tx.Prepare(addCustomerStat)
		if err != nil {
			return id, err.Error()
		}
		defer stmt.Close()
		row := stmt.QueryRow(request.FirstName, request.LastName, request.Email, request.Phone, request.Gender, request.Status, request.HasOrders, request.Active, util.CREATED_BY)
		errMsg := ""
		if row.Err() != nil {
			errMsg = row.Err().Error()
		}
		if errMsg != "" {
			log.Printf("Error while creating customer in DB: %v", errMsg)
			if strings.Contains(errMsg, "unique constraint \"customer_phone_key\"") {
				errMsg = "Customer with the given phone already exists"
			}
			return id, errMsg
		}
		fmt.Printf("rows returned :%v\n", row)
		//errMsg = row.Scan(&id).Error()
		if errMsg != "" {
			log.Printf("Error while creating customer in DB: %v", errMsg)
			return id, errMsg
		}
	}
	{
		err := tx.Commit()
		if err != nil {
			return id, err.Error()
		}
	}
	//id = *custId
	return id, ""
}

func (customerRepo *customerRepository) UpdateCustomer(request model.CustomerRequest) error {
	currentTime := time.Now()
	updateCustomerStat := `
		UPDATE "customer" 
		SET first_name = $1,last_name = $2,phone = $3,gender = $4,active = $5, updated_by = $6, updated_at = $7
		WHERE id = $6
		`
	err := customerRepo.db.QueryRow(updateCustomerStat, request.FirstName, request.LastName, request.Phone,
		request.Gender, request.Active, request.Id, util.UPDATED_BY, currentTime).Err()

	if err != nil {
		log.Printf("Error while updating customer in DB: %v", err.Error())
		return err
	}

	return nil
}

func (customerRepo *customerRepository) GetCustomers(active bool) ([]model.Customer, error) {
	fmt.Printf("GetCustomers active : %v\n", active)
	customerQuery := `
		SELECT id,first_name,last_name,gender,email,phone,source,status,active,account_id FROM "customer" 
		WHERE active = $1
		`
	rows, err := customerRepo.db.Query(customerQuery, active)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// An Customer slice to hold data from returned rows.
	var customers []model.Customer

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var cust model.Customer
		var source *string
		var acct *int64
		if err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Gender, &cust.Email, &cust.Phone, &source,
			&cust.Status, &cust.Active, &acct); err != nil {
			if source != nil {
				cust.Source = *source
			}
			if acct != nil {
				fmt.Printf("Customer Account id :%v\n", *acct)
			}
			return customers, err
		}
		customers = append(customers, cust)
	}
	if err = rows.Err(); err != nil {
		return customers, err
	}
	return customers, nil

}

func (customerRepo *customerRepository) GetCustomerById(id int64) (model.Customer, error) {
	customerQuery := `
		SELECT id,first_name,last_name,gender,email,phone,source,status,active,account_id FROM "customer" 
		WHERE id = $1
		`
	row := customerRepo.db.QueryRow(customerQuery, id)
	return customerRepo.scanResult(row)
}

func (customerRepo *customerRepository) GetCustomerByEmail(email string) (model.Customer, error) {
	customerQuery := `
		SELECT id,first_name,last_name,gender,email,phone,source,status,active,account_id FROM "customer" 
		WHERE email = $1
		`
	row := customerRepo.db.QueryRow(customerQuery, email)
	return customerRepo.scanResult(row)
}

func (customerRepo *customerRepository) UpdateCustomerStatus(id int64, status string) (model.Customer, error) {
	var cust model.Customer
	updateCustomerStat := `
		UPDATE "customer" 
		SET status = $1
		WHERE id = $2
		`
	err := customerRepo.db.QueryRow(updateCustomerStat, status, id).Err()

	if err != nil {
		log.Printf("Error while updating customer status in DB: %v", err.Error())
		return cust, err
	}
	customerQuery := `SELECT * FROM "customer" where id = $1`
	row := customerRepo.db.QueryRow(customerQuery, id)
	if row == nil {
		log.Printf("Error while retrieving updated customer record: ")
		return cust, nil
	}
	cust, err = customerRepo.scanResult(row)
	return cust, err
}

func (customerRepo *customerRepository) scanResult(row *sql.Row) (model.Customer, error) {
	// A Customer object to hold data from returned row.
	var cust model.Customer

	if row == nil {
		return cust, nil
	}

	// Loop through rows, using Scan to assign column data to struct fields.
	if err := row.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Gender, &cust.Email, &cust.Phone, &cust.Source,
		&cust.Status, &cust.Active, &cust.Account); err != nil {
		return cust, err
	}
	return cust, nil
}
