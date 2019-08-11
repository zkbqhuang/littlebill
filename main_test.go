package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExecuteTransferHandler(t *testing.T) {
	checker := &httpChecker{
		method:   "POST",
		path:     "/transfers/execute",
		body:     `{"account_from": "acc_from", "account_to": "acc_to", "amount": 100}`,
		expected: `{"ok":true,"message":"OK"}`,
		status:   http.StatusOK,
	}
	errInit := checker.init()
	if errInit != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errInit)
	}
	defer checker.db.Close()

	checker.mock.ExpectExec("INSERT INTO transfers \\(account_from, account_to, amount\\) VALUES \\(\\$1, \\$2, \\$3\\)").WithArgs("acc_from", "acc_to", 100).WillReturnResult(sqlmock.NewResult(1, 1))

	errCheck := checker.check()
	if errCheck != nil {
		t.Error(errCheck)
	}
}

func TestListTransferHandler(t *testing.T) {
	checker := &httpChecker{
		method:   "GET",
		path:     "/transfers/acc_1",
		body:     "",
		expected: `{"transfers":[{"account_from":"acc_1","account_to":"acc_2","amount":99,"direction":"outgoing"},{"account_from":"acc_2","account_to":"acc_1","amount":100,"direction":"incoming"}]}`,
		status:   http.StatusOK,
	}
	errInit := checker.init()
	if errInit != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errInit)
	}
	defer checker.db.Close()

	rows := sqlmock.NewRows([]string{"direction", "account_from", "account_to", "amount"}).
		AddRow("outgoing", "acc_1", "acc_2", 99).
		AddRow("incoming", "acc_2", "acc_1", 100)
	checker.mock.ExpectQuery("SELECT CASE WHEN account_from = \\$1 THEN 'outgoing' ELSE 'incoming' END AS direction, account_from, account_to, amount FROM transfers WHERE account_from = \\$1 or account_to = \\$1 ORDER BY dt DESC LIMIT \\$2").WithArgs("acc_1", 20).WillReturnRows(rows)

	errCheck := checker.check()
	if errCheck != nil {
		t.Error(errCheck)
	}
}

func TestListAccountsHandler(t *testing.T) {
	checker := &httpChecker{
		method:   "GET",
		path:     "/accounts",
		body:     "",
		expected: `{"accounts":[{"Name":"acc_1","Balance":1000,"Currency":"USD"}]}`,
		status:   http.StatusOK,
	}
	errInit := checker.init()
	if errInit != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errInit)
	}
	defer checker.db.Close()

	rows := sqlmock.NewRows([]string{"account", "balance", "curr"}).
		AddRow("acc_1", "1000", "USD")
	checker.mock.ExpectQuery("^SELECT (.+) FROM (.+)").WillReturnRows(rows)

	errCheck := checker.check()
	if errCheck != nil {
		t.Error(errCheck)
	}
}

func TestCreateAccountHandler(t *testing.T) {
	checker := &httpChecker{
		method:   "PUT",
		path:     "/accounts",
		body:     `{"name": "account_name", "balance": 1000, "currency": "USD"}`,
		expected: `{"ok":true,"message":"OK"}`,
		status:   http.StatusOK,
	}
	errInit := checker.init()
	if errInit != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", errInit)
	}
	defer checker.db.Close()

	checker.mock.ExpectBegin()
	checker.mock.ExpectExec("INSERT INTO accounts").WithArgs("account_name").WillReturnResult(sqlmock.NewResult(1, 1))
	checker.mock.ExpectExec("INSERT INTO balances").WithArgs("account_name", 1000).WillReturnResult(sqlmock.NewResult(1, 1))
	checker.mock.ExpectExec("INSERT INTO currs").WithArgs("account_name", "USD").WillReturnResult(sqlmock.NewResult(1, 1))
	checker.mock.ExpectCommit()

	errCheck := checker.check()
	if errCheck != nil {
		t.Error(errCheck)
	}
}

type httpChecker struct {
	db       *sql.DB
	mock     sqlmock.Sqlmock
	expected string
	body     string
	method   string
	path     string
	status   int
}

func (c *httpChecker) init() (errSQL error) {
	c.db, c.mock, errSQL = sqlmock.New()
	return errSQL
}

func (c *httpChecker) check() error {
	var body io.Reader
	if c.body != "" {
		body = strings.NewReader(c.body)
	}
	req, errReq := http.NewRequest(c.method, c.path, body)
	if errReq != nil {
		return errReq
	}

	recorder := httptest.NewRecorder()
	handler := makeRouter(c.db)

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != c.status {
		return fmt.Errorf("handler returned wrong status code: got %v want %v",
			status, c.status)
	}

	got := strings.TrimSpace(recorder.Body.String())
	if got != c.expected {
		return fmt.Errorf("handler returned unexpected body: got '%v' want '%v'",
			got, c.expected)
	}
	return nil
}
