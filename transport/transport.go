package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/mntor/littlebill/account"
	"github.com/mntor/littlebill/services"
	"github.com/mntor/littlebill/transfer"
)

var (
	// ErrBadRouting error used to point to routing configuration problem
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// MakeCreateAccountEndpoint returns endpoint for create account method
func MakeCreateAccountEndpoint(svc services.AccountsService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		account := account.Account{
			Name:     req.Name,
			Balance:  req.Balance,
			Currency: req.Currency,
		}
		errExecute := svc.Create(account)
		if errExecute != nil {
			return StatusResponse{OK: false, Message: errExecute.Error()}, nil
		}
		return StatusResponse{OK: true, Message: "OK"}, nil
	}
}

// MakeListAccountEndpoint returns endpoint for list account method
func MakeListAccountEndpoint(svc services.AccountsService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		accountsList, errExecute := svc.List()
		if errExecute != nil {
			return listAccountResponse{Error: errExecute.Error()}, nil
		}
		return listAccountResponse{Accounts: accountsList}, nil
	}
}

// MakeListTransferEndpoint returns endpoint for list transfer method
func MakeListTransferEndpoint(svc services.TransfersService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(listTransferRequest)
		transfersList, errExecute := svc.List(req.Account)
		if errExecute != nil {
			return listTransferResponse{Error: errExecute.Error()}, nil
		}
		return listTransferResponse{Transfers: transfersList}, nil
	}
}

// MakeExecuteTransferEndpoint returns endpoint for execute transfer method
func MakeExecuteTransferEndpoint(svc services.TransfersService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(transfer.Transfer)
		errExecute := svc.Execute(req)
		if errExecute != nil {
			return StatusResponse{OK: false, Message: errExecute.Error()}, errExecute
		}
		return StatusResponse{OK: true, Message: "OK"}, nil
	}
}

// DecodeListAccountRequest decodes requst to list account method
func DecodeListAccountRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return
}

// DecodeListTransferRequest decodes requst to list transfer method
func DecodeListTransferRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	account, ok := vars["account"]
	if !ok {
		return nil, ErrBadRouting
	}
	return listTransferRequest{Account: account}, nil
}

// DecodeExecuteTransferRequest decodes execute request to transfers
func DecodeExecuteTransferRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request transfer.Transfer
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// DecodeCreateAccountRequest decodes create request to accounts
func DecodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// EncodeResponse encodes any struct to json and sends to client
func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// EncodeError encodes error into json
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest) // TODO: make different codes depending on actual error
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// StatusResponse contains response with status of method
type StatusResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// listTransferRequest contains request with data to list transfers method
type listTransferRequest struct {
	Account string `json:"account"`
}

// listTransferResponse contains response with data to list transfers method
type listTransferResponse struct {
	Transfers []transfer.Transfer `json:"transfers"`
	Error     string              `json:"error,omitempty"`
}

// listAccountResponse contains response with data to list accounts method
type listAccountResponse struct {
	Accounts []account.Account `json:"accounts"`
	Error    string            `json:"error,omitempty"`
}

// createAccountRequest contains request with data to create account method
type createAccountRequest struct {
	Name     string `json:"name"`
	Balance  int    `json:"balance"`
	Currency string `json:"currency"`
}
