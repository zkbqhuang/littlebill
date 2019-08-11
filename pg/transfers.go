package pg

import (
	"database/sql"

	"github.com/mntor/littlebill/transfer"
)

type transferRepository struct {
	db *sql.DB
}

// NewTransferRepository returns a new instance transfer repository.
func NewTransferRepository(db *sql.DB) transfer.Repository {
	return &transferRepository{
		db: db,
	}
}

// Store stores transfer in DB
func (tr *transferRepository) Store(transfer transfer.Transfer) error {
	_, errQuery := tr.db.Exec(
		"INSERT INTO transfers (account_from, account_to, amount) VALUES ($1, $2, $3)",
		transfer.AccountFrom,
		transfer.AccountTo,
		transfer.Amount,
	)
	return errQuery
}

// List returns list of transfers by account name
func (tr *transferRepository) List(account string) (ret []transfer.Transfer, err error) {
	rows, errQuery := tr.db.Query("SELECT CASE WHEN account_from = $1 THEN 'outgoing' ELSE 'incoming' END AS direction, account_from, account_to, amount FROM transfers WHERE account_from = $1 or account_to = $1 ORDER BY dt DESC LIMIT $2",
		account,
		20, // TODO: let pass this value from user
	)
	if errQuery != nil {
		return nil, errQuery
	}
	for rows.Next() {
		var entry transfer.Transfer
		errScan := rows.Scan(&entry.Direction, &entry.AccountFrom, &entry.AccountTo, &entry.Amount)
		if errScan != nil {
			return nil, errScan
		}
		ret = append(ret, entry)
	}
	return ret, nil
}
