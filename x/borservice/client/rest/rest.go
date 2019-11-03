package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	restName = "name"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/prizes", storeName), namesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/prizes", storeName), initPrizeHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/prizes", storeName), setPrizeHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/prizes/{%s}", storeName, restName), resolveNameHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/prizes/{%s}/whois", storeName, restName), whoIsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/prizes", storeName), deletePrizeHandler(cliCtx)).Methods("DELETE")
}
