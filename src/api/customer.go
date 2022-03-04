package api

import (
	"crm-app/src/model"
	"crm-app/src/repository"
	"errors"
	"log"
)

// CustomerService contains the business logic
type CustomerService interface {
	AddNewCustomer(customer model.CustomerRequest) (int, error)
	UpdateCustomer(customer model.CustomerRequest) error
	GetCustomers(active bool) ([]model.Customer, error)
	GetCustomerDetails(id int64, email string) (model.Customer, error)
	ConvertLead(id int64) (model.Customer, error)
}

type customerService struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerService(customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{
		customerRepo: customerRepo,
	}
}

func (custService *customerService) AddNewCustomer(customer model.CustomerRequest) (int, error) {
	var id int
	if customer.FirstName == "" {
		return id, errors.New("First name is required")
	}
	if customer.LastName == "" {
		return id, errors.New("Last name is required")
	}
	if customer.Email == "" {
		return id, errors.New("Email is required")
	}
	if customer.Phone == "" {
		return id, errors.New("Phone is required")
	}
	customer.Active = true
	customer.Status = "Lead"
	var svcErr error
	var errMsg string
	cust, err := custService.customerRepo.GetCustomerByEmail(customer.Email)
	if err != nil {
		log.Printf("Error retrieving customer details with email %v\n", customer.Email)
		log.Println(err)
	}
	if (model.Customer{} == cust) {
		id, errMsg = custService.customerRepo.CreateNewCustomer(customer)
		if errMsg != "" {
			svcErr = errors.New(errMsg)
		}
	} else {
		svcErr = errors.New("Customer with the email already exists")
	}
	if svcErr != nil {
		return id, svcErr
	}
	return id, nil
}

func (custService *customerService) UpdateCustomer(customer model.CustomerRequest) error {
	if customer.FirstName == "" {
		return errors.New("First name is required")
	}
	if customer.LastName == "" {
		return errors.New("Last name is required")
	}
	if customer.Email == "" {
		return errors.New("Email is required")
	}
	if customer.Phone == "" {
		return errors.New("Phone is required")
	}
	err := custService.customerRepo.UpdateCustomer(customer)
	if err != nil {
		return err
	}
	return nil
}

func (custService *customerService) GetCustomers(active bool) ([]model.Customer, error) {
	return custService.customerRepo.GetCustomers(active)
}

func (custService *customerService) GetCustomerDetails(id int64, email string) (model.Customer, error) {
	var cust model.Customer
	var err error
	if id != 0 {
		cust, err = custService.customerRepo.GetCustomerById(id)
	} else if email != "" && len(email) != 0 {
		cust, err = custService.customerRepo.GetCustomerByEmail(email)
	}
	return cust, err
}

func (custService *customerService) ConvertLead(id int64) (model.Customer, error) {
	//update the customer status from 'Lead' to 'Customer'
	cust, err := custService.customerRepo.UpdateCustomerStatus(id, "Customer")
	if err != nil {
		log.Printf("Error while converting lead to customer, could not update status \n", err.Error())
	}
	return cust, err
}
