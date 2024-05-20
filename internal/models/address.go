package models

import (
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"github.com/rs/xid"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"gorm.io/gorm"
)

type AddressRequest struct {
	PathTxID
	Status  string      `json:"status" enum:"not_started,in_progress,completed,cancelled,in_review"`
	Address AddressData `json:"address" description:"Address Task data"`
}

type AddressData struct {
	Data AddressFields `json:"data" description:"Address Task data"`
}

type AddressFields struct {
	Address1       *string `json:"address_1" gorm:"column:address_1" description:"First line of address" example:"10 Downing Street"`
	Address2       *string `json:"address_2" gorm:"column:address_2" description:"Second line of address" example:"Flat 42"`
	BuildingName   *string `json:"building_name" description:"Name of building" example:"Buckingham Palace"`
	BuildingNumber *string `json:"building_number" description:"Number of building" example:"12"`
	Country        *string `json:"country" description:"Three letter countrycode" example:"GBR"`
	FlatNumber     *string `json:"flat_number" description:"Flat number" example:"42"`
	Postcode       *string `json:"postcode" description:"Postcode" example:"SW1A 1AA"`
	State          *string `json:"state" description:"Name of state" example:"Wyoming"`
	Street         *string `json:"street" description:"Name of street" example:"Main street"`
	SubStreet      *string `json:"sub_street" description:"Name of sub street" example:"Lane 1"`
	Town           *string `json:"town" description:"Name of city or town" example:"London"`
}

func (p AddressRequest) ToModel(subtype string) *Address {
	base := Base{TransactionID: p.TxID, Type: subtype, Status: p.Status}
	return &Address{
		Base:           base,
		Address1:       p.Address.Data.Address1,
		Address2:       p.Address.Data.Address2,
		BuildingName:   p.Address.Data.BuildingName,
		BuildingNumber: p.Address.Data.BuildingNumber,
		Country:        p.Address.Data.Country,
		FlatNumber:     p.Address.Data.FlatNumber,
		Postcode:       p.Address.Data.Postcode,
		State:          p.Address.Data.State,
		Street:         p.Address.Data.Street,
		SubStreet:      p.Address.Data.SubStreet,
		Town:           p.Address.Data.Town,
	}
}

type AddressResponse struct {
	TaskFields
	AddressExpectation `json:"address" description:"Expectation metadata"`
}

type AddressExpectation struct {
	ExpectationFields
	Data AddressFields `json:"data" description:"Task data"`
}

type Address struct {
	Base
	Address1       *string `json:"address_1,omitempty" gorm:"column:address_1" description:"First line of address" example:"10 Downing Street"`
	Address2       *string `json:"address_2,omitempty" gorm:"column:address_2" description:"Second line of address" example:"Flat 42"`
	BuildingName   *string `json:"building_name,omitempty" description:"Name of building" example:"Buckingham Palace"`
	BuildingNumber *string `json:"building_number,omitempty" description:"Number of building" example:"12"`
	Country        *string `json:"country" description:"Three letter countrycode" example:"GBR"`
	FlatNumber     *string `json:"flat_number,omitempty" description:"Flat number" example:"42"`
	Postcode       *string `json:"postcode,omitempty" description:"Postcode" example:"SW1A 1AA"`
	State          *string `json:"state,omitempty" description:"Name of state" example:"Wyoming"`
	Street         *string `json:"street,omitempty" description:"Name of street" example:"Main street"`
	SubStreet      *string `json:"sub_street,omitempty" description:"Name of sub street" example:"Lane 1"`
	Town           *string `json:"town,omitempty" description:"Name of city or town" example:"London"`
	Changed        bool    `json:"changed"`
}

func (a *Address) ToResponse() *AddressResponse {
	var (
		data AddressRequest
		exp  AddressExpectation
	)
	copier.Copy(&data.Address.Data, a)
	copier.Copy(&exp, a.Base)
	exp.Data = data.Address.Data
	exp.Status = MapPaStatus(a.Base.Status)
	return &AddressResponse{
		TaskFields: TaskFields{
			Status:     getCombinedStatus(a.Base.Status),
			ReasonCode: GetReasonCodeForTask(a.Base.ReasonCode),
			Type:       a.Base.Type,
			CreatedAt:  a.Base.CreatedAt,
		},
		AddressExpectation: exp,
	}
}

func (Address) TableName() string {
	return "address"
}

func (t Address) CopyIDs() *Address {
	var tt Address
	tt.TransactionID = t.TransactionID
	return &tt
}

func (t Address) NewID() *Address {
	id := xid.New().String()
	t.ID = &id
	t.NewCreated()
	return &t
}

func (t Address) NewCreated() {
	t.CreatedAt = time.Now()
	t.NewUpdated()
}

func (t Address) NewUpdated() {
	t.UpdatedAt = time.Now()
}

func (t Address) Delete() {
	t.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

func (t Address) GetIDs() (*string, *string, string) {
	return t.TransactionID, t.ID, t.Type
}

func (t Address) GetBase() Base {
	return t.Base
}

func (t Address) GetBaseItem(base Base) *Address {
	return &Address{Base: base}
}

func (t Address) SetBase(base Base) *Address {
	t.Base = base
	return &t
}

func (t Address) Equals(addr Address) bool {
	t.Base = Base{}
	addr.Base = Base{}
	t.Changed = false
	addr.Changed = false
	return reflect.DeepEqual(t, addr)
}

func (t Address) IsEmpty() bool {
	// Check if this is an empty struct
	base := Base{}
	t.Base = base
	return Address{Base: base} == t
}

func (t Address) IsValidToSubmit() bool {
	// Only country is required to submit
	// Other data is optional for now. Require confirmation
	if t.Country == nil {
		return false
	}

	return true

}

func (t Address) FillEmpty() Address {
	if t.Address1 == nil || (t.Address1 != nil && *t.Address1 == "") {
		t.Address1 = internal.StrPtr("")
	}
	if t.Address2 == nil || (t.Address2 != nil && *t.Address2 == "") {
		t.Address2 = internal.StrPtr("")
	}
	if t.BuildingName == nil || (t.BuildingName != nil && *t.BuildingName == "") {
		t.BuildingName = internal.StrPtr("")
	}
	if t.BuildingNumber == nil || (t.BuildingNumber != nil && *t.BuildingNumber == "") {
		t.BuildingNumber = internal.StrPtr("")
	}
	if t.FlatNumber == nil || (t.FlatNumber != nil && *t.FlatNumber == "") {
		t.FlatNumber = internal.StrPtr("")
	}
	if t.Postcode == nil || (t.Postcode != nil && *t.Postcode == "") {
		t.Postcode = internal.StrPtr("")
	}
	if t.State == nil || (t.State != nil && *t.State == "") {
		t.State = internal.StrPtr("")
	}
	if t.Street == nil || (t.Street != nil && *t.Street == "") {
		t.Street = internal.StrPtr("")
	}
	if t.SubStreet == nil || (t.SubStreet != nil && *t.SubStreet == "") {
		t.SubStreet = internal.StrPtr("")
	}
	if t.Town == nil || (t.Town != nil && *t.Town == "") {
		t.Town = internal.StrPtr("")
	}

	return t
}
