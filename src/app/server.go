package app

import (
	"crm-app/src/api"
	"log"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router          *gin.Engine
	customerService api.CustomerService
	accountService  api.AccountService
}

func NewServer(router *gin.Engine, customerService api.CustomerService, accountService api.AccountService) *Server {
	return &Server{
		router:          router,
		customerService: customerService,
		accountService:  accountService,
	}
}

func (s *Server) Run() error {
	//initialise the routes
	routes := s.Routes()

	//start the server
	err := routes.Run()

	if err != nil {
		log.Printf("Error running the server: %v", err)
		return err
	}

	return nil
}
