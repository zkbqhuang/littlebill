package pg

import (
	"database/sql"
	"fmt"
)

// PostgreCredentials contains credentials needed for connection to server
type PostgreCredentials struct {
	Username string
	Password string
	DBName   string
	Host     string
}

// NewPostgreConnection returns new connection to DB or error
// TODO: use pool of connections
func NewPostgreConnection(cred PostgreCredentials) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cred.Username,
		cred.Password,
		cred.Host,
		cred.DBName,
	)
	db, errOpen := sql.Open("postgres", connStr)
	if errOpen != nil {
		return nil, errOpen
	}
	errPing := db.Ping()
	return db, errPing
}
