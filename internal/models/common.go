package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"gorm.io/gorm"
)

type Shared struct {
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type Base struct {
	Shared        `json:"-"`
	ID            *string `json:"id" gorm:"uniqueIndex;primaryKey;<-:create;column:id"`
	Type          string  `json:"type" gorm:"primaryKey;<-:create;column:type"`
	PaType        string  `json:"pa_type" gorm:"primaryKey;<-:create;column:pa_type"`
	TransactionID *string `json:"transaction_id" path:"txID" gorm:"primaryKey;<-:create;column:transaction_id"`
	Status        string  `json:"status"`
	ExpectationID *string `json:"expectation_id" gorm:"column:expectation_id"`
	ReasonCode    string  `json:"reason_code"`
	ValidationErr string  `json:"validation_err" gorm:"-"`
}

type TaskFields struct {
	Status     string    `json:"status" description:"Status of the task" example:"in_progress" enum:"not_started,in_progress,completed,cancelled,in_review"`
	ReasonCode string    `json:"reason_code" description:"Reason code for failure"`
	Type       string    `json:"type,omitempty" description:"Type of task" example:"address"`
	CreatedAt  time.Time `json:"created_at" description:"Task creation time"`
}

type ExpectationFields struct {
	Status        string    `json:"status" description:"Status of the task" example:"in_progress" enum:"not_started,in_progress,completed,cancelled,in_review"`
	ExpectationID string    `json:"expectation_id" description:"Expectation ID"`
	ValidationErr *string   `json:"validation_err" description:"Validation errors"`
	CreatedAt     time.Time `json:"created_at" description:"Task creation time"`
}

type PathTxID struct {
	TxID *string `json:"-" path:"txID"`
}

type EmptyResponse struct {
	ExpectationFields
	Data interface{} `json:"data"`
}

type NullTime struct {
	Time  time.Time `json:"-"`
	Valid bool      `json:"-"`
}

func (n *NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		time := n.Time.Format(internal.DateLayout)
		return json.Marshal(time)
	}
	return json.Marshal(nil)
}

func (n *NullTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		n.Valid = false
		return nil
	}
	// 0001-01-01 00:00:00 +0000 UTC
	datetime, err := time.Parse(internal.DateLayout, s)
	if err != nil {
		return errors.Wrap(err, "parsing date")
	}

	n.Time = datetime
	n.Valid = true

	return nil
}

func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		n.Time, n.Valid = time.Time{}, false
		return nil
	}
	n.Valid = true
	n.Time = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

// 2000-02-22T00:00:00.000Z
func (n *NullTime) ToString() string {
	if n.Valid {
		return n.Time.UTC().Format(time.RFC3339)
	}
	return ""
}

// List of all internal models that are supported with simple CRUD
type DataType interface {
	Address | Dob | Document | Name
}

type Item[T DataType] interface {
	NewID() *T                                                // Returns type with new
	CopyIDs() *T                                              // Returns a copy with same Type, TransactionID and ID as the original
	NewCreated()                                              // Sets type with Created_At and Updated set to current time
	NewUpdated()                                              // Sets type with Updated_At set to current time
	Delete()                                                  // Sets type with Deleted_At set to current time
	GetIDs() (TransactionID *string, ID *string, Type string) // returns TransactionID, ID and Type from underlying DataType
	GetBase() Base
	GetBaseItem(base Base) *T
	SetBase(b Base) *T
	TableName() string
	IsEmpty() bool
	IsValidToSubmit() bool
}

type ResponseType[T DataType] interface {
	AddressResponse | DocumentResponse | PersonalInformationResponse
}

type ResponseItem[T DataType, R ResponseType[T]] interface {
	ToResponse() *R
}

type RequestType[T DataType] interface {
	AddressRequest | DocumentRequest | PersonalInformationRequest
	ToModel(subtype string) *T
}

func safe(val *string) *string {
	if val != nil {
		return val
	}
	return nil
}
