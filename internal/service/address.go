package service

import (
	"context"
	"errors"

	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func FindAddress(tx_id string) (models.Address, error) {
	address, err := MyService.DataStore.FindAddress(context.TODO(), tx_id)
	return *address, err
}

func UpdateAddress(transactionId string, newAddr models.Address) (*models.Address, error) {

	address, _ := FindAddress(transactionId)

	if address.Status == "completed" {
		return nil, errors.New("Adddress is already complete")
	} else if address.Status == "pending" || address.Status == "accepted" {
		// update address
		MyService.DataStore.GetStore().UpdateAddress(context.TODO(), transactionId, newAddr)
		return &newAddr, nil
	}
	return nil, nil
}
