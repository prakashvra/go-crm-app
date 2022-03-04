package api

import (
	"crm-app/src/model"
	"crm-app/src/repository"
	"errors"
	"log"
)

type AccountService interface {
	AddNewAccount(account model.AccountRequest) (int64, error)
	//UpdateAccount(account model.AccountRequest) error
	//GetAccounts(active bool) ([]model.Account, error)
	//GetAccountDetails(id int64, email string) (model.Account, error)
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}

func (acctService *accountService) AddNewAccount(account model.AccountRequest) (int64, error) {
	var acctId int64
	var svcErr error
	var errMsg string
	acct, err := acctService.accountRepo.GetAccountByEmail(account.Email)
	if err != nil {
		log.Printf("Error retrieving account details with email %v\n", account.Email)
		log.Println(err)
	}
	if acct.Id != 0 {
		svcErr = errors.New("Account with the email already exists")
	} else {
		acctId, errMsg = acctService.accountRepo.CreateNewAccount(account)
		svcErr = errors.New(errMsg)
	}
	if svcErr != nil {
		return acctId, svcErr
	}
	return acctId, nil
}
