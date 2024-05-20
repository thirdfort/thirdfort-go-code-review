package web

import (
	"github.com/loopfz/gadgeto/tonic"
	"github.com/wI2L/fizz"
)

func (s *WebService) setV1Routes() {
	v1Router := s.Fizz.Group("/v1", "v1", "Main endpoint (v1) to the API")

	// Transaction routes
	txRouter := v1Router.Group("/transaction", "transaction", "Transaction endpoints")
	s.transactionRoutes(txRouter)
}

func (s *WebService) transactionRoutes(router *fizz.RouterGroup) {
	router.Handle("", "GET", []fizz.OperationOption{
		fizz.ID("GetAllTransactions"),
		fizz.Summary("GET all transactions"),
	}, tonic.Handler(s.GetTransactions, 200))

	router.Handle(":txID", "GET", []fizz.OperationOption{
		fizz.ID("GetTransactionByID"),
		fizz.Summary("GET Transaction by ID"),
	}, tonic.Handler(s.GetTransaction, 200))

	router.Handle(":txID", "PATCH", []fizz.OperationOption{
		fizz.ID("PatchTransactionByID"),
		fizz.Summary("GET Transaction by ID"),
	}, tonic.Handler(s.PatchTransaction, 200))
}
