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

	// Address Routes
	addressRouter := v1Router.Group("/address", "address", "Address endpoints")
	s.AddressRoutes(addressRouter)
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

func (s *WebService) AddressRoutes(router *fizz.RouterGroup) {
	router.Handle("", "GET", []fizz.OperationOption{
		fizz.ID("GetAddressByID"),
		fizz.Summary("GET Address by ID"),
	}, tonic.Handler(FindAddress, 200))

	router.Handle("", "PATCH", []fizz.OperationOption{
		fizz.ID("PatchAddressByID"),
		fizz.Summary("PATCH Address by ID"),
	}, tonic.Handler(UpdateAddress, 200))
}
