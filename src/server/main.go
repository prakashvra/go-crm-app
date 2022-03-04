package main

import (
	"crm-app/src/api"
	"crm-app/src/app"
	"crm-app/src/repository"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during startup %s\\n", err)
		os.Exit(1)
	}
}

func run() error {
	// define connection string for postgresql crm db
	postgres_db_server := os.Getenv("POSTGRES_DB_SERVER")
	connString := "postgres://crmuser:crmuser123@" + postgres_db_server + "/crm_app?sslmode=disable"

	db, err := setupDBConn(connString)

	if err != nil {
		return err
	}

	// create customer repository
	customerRepo := repository.NewCustomerRepo(db)
	// create customer service
	customerService := api.NewCustomerService(customerRepo)

	// create account repo
	accountRepo := repository.NewAccountRepo(db)
	// create account service
	accountService := api.NewAccountService(accountRepo)

	// create router
	router := gin.Default()
	router.Use(cors.Default())

	// create server
	server := app.NewServer(router, customerService, accountService)

	//start the server
	err = server.Run()
	if err != nil {
		return err
	}
	return nil
}

func setupDBConn(connString string) (*sql.DB, error) {
	// connect to postgresql using the connection string
	db, err := sql.Open("postgres", connString)
	log.Println("Connected to DB :", db)
	if err != nil {
		log.Println("Error while connecting to DB :", err)
		return nil, err
	}

	// check the connection by pinging
	err = db.Ping()
	if err != nil {
		log.Println("Error while pinging DB :", err)
		return nil, err
	}
	return db, err
}
