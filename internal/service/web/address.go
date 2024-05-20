package web

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"github.com/thirdfort/thirdfort-go-code-review/internal/service"
)

func FindAddress(c *gin.Context, req *models.Transaction) (*models.TransactionResponse, error) {
	LogAddressRequest(c, req.ID)
	addr, err := service.FindAddress(*req.ID)
	if err != nil {
		panic(err)
	}
	return &models.TransactionResponse{Address: &addr}, nil
}

func UpdateAddress(c *gin.Context, req *models.Transaction) (*models.TransactionResponse, error) {
	LogAddressRequest(c, req.ID)
	addr, err := service.UpdateAddress(*req.ID, *req.Address)
	if err != nil {
		panic(err)
	}
	return &models.TransactionResponse{Address: addr}, nil
}

func LogAddressRequest(c *gin.Context, txID *string) context.Context {
	ctx := getCtx(c)

	if txID != nil {
		fmt.Printf("{transactionID: %s}", *txID)
	}
	fmt.Printf("{Req: {method: %s, path: %s}}", c.Request.Method, c.Request.URL.Path)
	return ctx
}
