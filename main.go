package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/mntor/littlebill/pg"
	"github.com/mntor/littlebill/services"
	"github.com/mntor/littlebill/transport"

	_ "github.com/lib/pq"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	db, errDB := pg.NewPostgreConnection(pg.PostgreCredentials{
		Username: os.Getenv("PSQL_USERNAME"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Host:     os.Getenv("PSQL_HOSTNAME"),
		DBName:   os.Getenv("PSQL_DBNAME"),
	})
	if errDB != nil {
		logger.Log("msg", "DB", "error", errDB)
		os.Exit(1)
	}
	r := makeRouter(db)
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", r))
}

func makeRouter(db *sql.DB) (r *mux.Router) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(transport.EncodeError),
	}
	svcTransfer := &services.TransferService{
		TransfersRepo: pg.NewTransferRepository(db),
	}
	svcAccount := &services.AccountService{
		AccountRepo: pg.NewAccountRepository(db),
	}

	r = mux.NewRouter()
	r.Methods("POST").Path("/transfers/execute").Handler(httptransport.NewServer(
		transport.MakeExecuteTransferEndpoint(svcTransfer),
		transport.DecodeExecuteTransferRequest,
		transport.EncodeResponse,
		options...,
	))
	r.Methods("GET").Path("/transfers/{account}").Handler(httptransport.NewServer(
		transport.MakeListTransferEndpoint(svcTransfer),
		transport.DecodeListTransferRequest,
		transport.EncodeResponse,
		options...,
	))
	r.Methods("GET").Path("/accounts").Handler(httptransport.NewServer(
		transport.MakeListAccountEndpoint(svcAccount),
		transport.DecodeListAccountRequest,
		transport.EncodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/accounts").Handler(httptransport.NewServer(
		transport.MakeCreateAccountEndpoint(svcAccount),
		transport.DecodeCreateAccountRequest,
		transport.EncodeResponse,
		options...,
	))
	return
}
