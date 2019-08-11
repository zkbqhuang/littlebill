package services

import (
	"github.com/mntor/littlebill/account"
)

// AccountsService provides accounts operations
type AccountsService interface {
	List() ([]account.Account, error)
	Create(account.Account) error
}

// AccountService implements account service
type AccountService struct {
	AccountRepo account.Repository
}

// List returns all available accounts
func (ts *AccountService) List() ([]account.Account, error) {
	return ts.AccountRepo.List()
}

// Create creates new account in DB
func (ts *AccountService) Create(account account.Account) error {
	return ts.AccountRepo.Create(account)
}
