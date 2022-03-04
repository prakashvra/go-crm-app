package model

type CustomerRequest struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Source    string `json:"source"`
	Status    string `json:"status"`
	HasOrders bool   `json:"has_orders"`
	Active    bool   `json:"active"`
}

type Customer struct {
	Account
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Source    string `json:"source"`
	Status    string `json:"status"`
	Active    bool   `json:"active"`
	CreatedBy string `json:"createdBy"`
	CreatedDt string `json:"createdDt"`
	UpdatedBy string `json:"updatedBy"`
	UpdatedDt string `json:"updatedDt"`
}

type AccountRequest struct {
	Id              int    `json:"id"`
	Type            string `json:"type"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Gender          string `json:"gender"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Source          string `json:"source"`
	HasOrders       bool   `json:"has_orders"`
	BillingAddress1 string `json:"billingAddress1"`
	BillingAddress2 string `json:"billingAddress2"`
	City            string `json:"city"`
	Zipcode         string `json:"zipcode"`
	State           string `json:"state"`
	Country         string `json:"country"`
	Status          string `json:"status"`
	Active          bool   `json:"active"`
}

type Account struct {
	Id              int    `json:"id"`
	Type            string `json:"type"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Gender          string `json:"gender"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Source          string `json:"source"`
	HasOrders       bool   `json:"has_orders"`
	BillingAddress1 string `json:"billingAddress1"`
	BillingAddress2 string `json:"billingAddress2"`
	City            string `json:"city"`
	Zipcode         string `json:"zipcode"`
	State           string `json:"state"`
	Country         string `json:"country"`
	Status          string `json:"status"`
	Active          bool   `json:"active"`
	CreatedBy       string `json:"createdBy"`
	CreatedDt       string `json:"createdDt"`
	UpdatedBy       string `json:"updatedBy"`
	UpdatedDt       string `json:"updatedDt"`
}
