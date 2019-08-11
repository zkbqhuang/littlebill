package transfer

// Transfer contains all data of transfer
type Transfer struct {
	AccountFrom string `json:"account_from"`
	AccountTo   string `json:"account_to"`
	Amount      int    `json:"amount"`
	Direction   string `json:"direction,omitempty"`
}

// Repository provides access to transfers store
type Repository interface {
	Store(transfer Transfer) error
	List(account string) ([]Transfer, error)
}
