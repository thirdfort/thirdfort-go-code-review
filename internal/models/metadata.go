package models

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/datatypes"
)

type ExpectationMetadata struct {
	Shared
	ID              *string        `json:"id" description:"ID" example:"cq1qh5c23amg0302nqv0" gorm:"column:id;primaryKey"`
	TransactionID   *string        `json:"transaction_id" gorm:"column:transaction_id"`
	ExpectationID   *string        `json:"expectation_id" gorm:"column:expectation_id;primaryKey"`
	ExpectationType string         `json:"type" gorm:"column:type"`
	Data            datatypes.JSON `json:"data" gorm:"type:jsonb"`
}

func (ExpectationMetadata) TableName() string {
	return "expectation_metadata"
}

func (r ExpectationMetadata) NewID() *ExpectationMetadata {
	id := xid.New().String()
	r.ID = &id
	r.NewCreated()

	return &r
}

func (r ExpectationMetadata) NewCreated() {
	r.CreatedAt = time.Now()
	r.NewUpdated()
}

func (r ExpectationMetadata) NewUpdated() {
	r.UpdatedAt = time.Now()
}

func (r ExpectationMetadata) GetIDs() (*string, *string, *string) {
	return r.ID, r.TransactionID, r.ExpectationID
}
