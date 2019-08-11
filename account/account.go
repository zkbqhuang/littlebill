package account

// Account contains all data account entry
type Account struct {
	Name     string
	Balance  int
	Currency string
}

// Repository provides access to balance store
type Repository interface {
	List() ([]Account, error)
	Create(Account) error
}
