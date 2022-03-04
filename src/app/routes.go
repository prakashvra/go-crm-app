package app

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) Routes() *gin.Engine {
	router := s.router

	endpoints := router.Group("/v1/api")
	{
		endpoints.GET("/status", s.ApiStatus())
		endpoints.POST("/lead", s.CreateCustomer())
		endpoints.PUT("/lead", s.UpdateCustomer())
		endpoints.PUT("/customer", s.UpdateCustomer())
		endpoints.GET("/customers/:id", s.GetCustomer())
		endpoints.GET("/customers", s.GetAllCustomers())
		endpoints.GET("/lead/:id/convert", s.ConvertLead())
	}
	return router
}
