package app

import (
	"crm-app/src/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) ApiStatus() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		response := map[string]string{
			"status":  "success",
			"message": "Customer service",
		}

		c.JSON(http.StatusOK, response)
	})
}

func (s *Server) CreateCustomer() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		fmt.Println(c.Request.Body)
		var customer model.CustomerRequest
		//map request JSON to CustomerRequest
		err := c.BindJSON(&customer)
		if err != nil {
			log.Printf("Error converting request JSON: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		fmt.Printf("Customer data : %v\n", customer.FirstName)
		//call create customer service
		var id int
		id, err = s.customerService.AddNewCustomer(customer)
		if err != nil {
			log.Printf("Error creating new customer using CustomerService : %v", err)
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		response := map[string]string{
			"status":     "Success",
			"message":    "Customer created successfuly",
			"customerId": string(id),
		}
		c.JSON(http.StatusOK, response)
	})
}

func (s *Server) UpdateCustomer() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		fmt.Println(c.Request.Body)
		var customer model.CustomerRequest
		//map request JSON to CustomerRequest
		err := c.BindJSON(&customer)
		if err != nil {
			log.Printf("Error converting request JSON: %v", err)
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		fmt.Printf("Customer data : %v\n", customer.Email)
		//call create customer service
		err = s.customerService.UpdateCustomer(customer)
		if err != nil {
			log.Printf("Error updating customer using CustomerService : %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := map[string]string{
			"status":  "Success",
			"message": "Customer details updated successfuly",
		}
		c.JSON(http.StatusOK, response)
	})
}

func (s *Server) GetCustomer() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		fmt.Printf("Customer id param : %v\n", id)
		//call get customer service
		var custId int64
		var converr error
		if id != "" {
			custId, converr = strconv.ParseInt(id, 0, 32)
			if converr != nil {
				fmt.Println("Error converting id param to integer", converr)
			}
		}
		customer, err := s.customerService.GetCustomerDetails(custId, "")
		if err != nil {
			log.Printf("Error getting customer details using CustomerService : %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := customer
		c.JSON(http.StatusOK, response)
	})
}

func (s *Server) GetAllCustomers() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		active := c.Param("active")
		fmt.Printf("Status param : %v\n", active)
		//call get customer service
		var isActive bool
		if active == "" && len(active) == 0 {
			isActive = true
		}
		customers, err := s.customerService.GetCustomers(isActive)
		if err != nil {
			log.Printf("Error getting customers using CustomerService : %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		response := customers
		c.JSON(http.StatusOK, response)
	})
}

func (s *Server) ConvertLead() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		id := c.Param("id")
		fmt.Printf("Customer id param : %v\n", id)
		var custId int64
		var converr error
		if id != "" {
			custId, converr = strconv.ParseInt(id, 0, 32)
			if converr != nil {
				fmt.Println("Error converting id param to integer", converr)
			}
		}
		//call customer convert lead service
		customer, err := s.customerService.ConvertLead(custId)
		if err != nil {
			log.Printf("Error converting lead to customer using CustomerService : %v", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		//create new account for customer
		var acctReq model.AccountRequest
		acctReq.FirstName = customer.FirstName
		acctReq.LastName = customer.LastName
		acctReq.Gender = customer.Gender
		acctReq.Email = customer.Email
		acctReq.Phone = customer.Phone
		acctReq.Source = customer.Source
		acctReq.HasOrders = customer.HasOrders
		acctReq.Status = "Created"
		acctReq.Active = true
		acctId, acctErr := s.accountService.AddNewAccount(acctReq)
		if acctErr != nil {
			log.Printf("Error creating account using AccountService : %v", acctErr)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		result := `{
			"status":  "Success",
			"message": "Lead converted successfuly",
			"data" : {"customerId" :` + string(custId) + `, "accountId" :` + string(acctId) + `}
		}`
		// Declared an empty interface
		var response map[string]interface{}
		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(result), &response)
		c.JSON(http.StatusOK, response)
	})
}
