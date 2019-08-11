package pg

import (
	"database/sql"

	"github.com/mntor/littlebill/account"
)

type accountRepository struct {
	db *sql.DB
}

// NewAccountRepository returns a new instance account repository.
func NewAccountRepository(db *sql.DB) account.Repository {
	return &accountRepository{
		db: db,
	}
}

// List returns list of transfers by account name
func (tr *accountRepository) List() (ret []account.Account, err error) {
	// TODO: add LIMIT to query
	rows, errQuery := tr.db.Query("SELECT b.account, b.balance, c.curr FROM balances AS b JOIN currs AS c ON c.account = b.account")
	if errQuery != nil {
		return nil, errQuery
	}
	for rows.Next() {
		var entry account.Account
		errScan := rows.Scan(&entry.Name, &entry.Balance, &entry.Currency)
		if errScan != nil {
			return nil, errScan
		}
		ret = append(ret, entry)
	}
	return ret, nil
}

// Create creates new account with init balance
func (tr *accountRepository) Create(account account.Account) error {
	trx, trxErr := tr.db.Begin()
	if trxErr != nil {
		return trxErr
	}
	_, errAcc := trx.Exec("INSERT INTO accounts VALUES ($1)", account.Name)
	if errAcc != nil {
		trx.Rollback()
		return errAcc
	}
	_, errBal := trx.Exec("INSERT INTO balances VALUES ($1, $2)", account.Name, account.Balance)
	if errBal != nil {
		trx.Rollback()
		return errBal
	}
	_, errCurr := trx.Exec("INSERT INTO currs VALUES ($1, $2)", account.Name, account.Currency)
	if errCurr != nil {
		trx.Rollback()
		return errCurr
	}

	return trx.Commit()
}
