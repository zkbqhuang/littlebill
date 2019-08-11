package services

import (
	"github.com/mntor/littlebill/transfer"
)

// TransfersService provides transfers operations
type TransfersService interface {
	Execute(transfer transfer.Transfer) error
	List(account string) ([]transfer.Transfer, error)
}

// TransferService implements
type TransferService struct {
	TransfersRepo transfer.Repository
}

// Execute makes transition from one account to another
func (ts *TransferService) Execute(transfer transfer.Transfer) error {
	return ts.TransfersRepo.Store(transfer)
}

// List returns list of last 20 transfers by account name
func (ts *TransferService) List(account string) ([]transfer.Transfer, error) {
	return ts.TransfersRepo.List(account)
}
