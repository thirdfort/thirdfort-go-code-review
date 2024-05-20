package models

import (
	"reflect"
	"time"

	"github.com/jinzhu/copier"
	"github.com/rs/xid"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"gorm.io/gorm"
)

type PersonalInformationRequest struct {
	PathTxID
	Name NameRequest `json:"name,omitempty"`
	Dob  DobRequest  `json:"dob,omitempty"`
}

func (p *PersonalInformationRequest) ToModel(subtype string) *PersonalInformation {
	base := Base{TransactionID: p.TxID, Type: subtype}
	return &PersonalInformation{
		Base: base,
		Name: &Name{Base: base, NameRequest: p.Name},
		Dob:  &Dob{Base: base, DobRequest: p.Dob},
	}
}

type PersonalInformationPatchRequest struct {
	PathTxID
	Status string           `json:"status" enum:"completed,cancelled" example:"completed" description:"Status of the task"`
	Name   NamePatchRequest `json:"name,omitempty"`
	Dob    DobPatchRequest  `json:"dob,omitempty"`
}

func (p *PersonalInformationPatchRequest) ToModel(subtype string) *PersonalInformation {
	base := Base{TransactionID: p.TxID, Type: subtype, Status: p.Status}
	baseName := base
	baseName.Type = internal.TypePersonalInformationName
	baseDob := base
	baseDob.Type = internal.TypePersonalInformationDob
	return &PersonalInformation{
		Base: base,
		Name: &Name{Base: baseName, NameRequest: p.Name.Data},
		Dob:  &Dob{Base: baseDob, DobRequest: p.Dob.Data},
	}
}

type PersonalInformationResponse struct {
	TaskFields
	Name NameResponse `json:"name" description:"Name data"`
	Dob  DobResponse  `json:"dob" description:"Date of Birth data"`
}

type NameRequest struct {
	First   *string `json:"first" validate:"required" description:"First name" example:"Mary"`
	Last    *string `json:"last" validate:"required" description:"Surname" example:"Smith"`
	Other   *string `json:"other" description:"Other/middle name" example:"Rodriguez"`
	Changed bool    `json:"changed"`
}

type NamePatchRequest struct {
	Data NameRequest `json:"data" description:"Task data"`
}

type NameResponse struct {
	ExpectationFields
	Data NameRequest `json:"data" description:"Task data"`
}

type DobRequest struct {
	Dob     NullTime `json:"dob" gorm:"column:dob" description:"Date of Birth" example:"1980-06-31"`
	Changed bool     `json:"changed"`
}

type DobPatchRequest struct {
	Data DobRequest `json:"data" description:"Task data"`
}

type DobResponse struct {
	ExpectationFields
	Data DobRequest `json:"data" description:"Task data"`
}

type PersonalInformation struct {
	Base `json:"-" gorm:"-"`
	Name *Name `json:"name"`
	Dob  *Dob  `json:"dob"`
}

func (PersonalInformation) TableName() string {
	return "personal_information"
}

func (p *PersonalInformation) ToResponse() *PersonalInformationResponse {
	var (
		name     NameResponse
		dob      DobResponse
		nameData NameRequest
		dobData  DobRequest
	)

	copier.Copy(&nameData, p.Name)
	copier.Copy(&name, p.Name)
	name.Data = nameData
	name.Status = MapPaStatus(p.Name.Base.Status)

	copier.Copy(&dobData, p.Dob)
	copier.Copy(&dob, p.Dob)
	dob.Data = dobData
	dob.Status = p.Dob.Base.Status

	return &PersonalInformationResponse{
		TaskFields: TaskFields{
			Status:     getCombinedStatus(p.Name.Base.Status, p.Dob.Base.Status),
			ReasonCode: GetReasonCodeForTask(p.Name.Base.ReasonCode, p.Dob.Base.ReasonCode),
			Type:       internal.TypePersonalInformation,
			CreatedAt:  p.Name.Base.CreatedAt,
		},
		Name: name,
		Dob:  dob,
	}
}

type Name struct {
	Base
	NameRequest
}

type Dob struct {
	Base
	DobRequest
}

func (Name) TableName() string {
	return "name"
}

func (t Name) CopyIDs() *Name {
	var tt Name
	tt.TransactionID = t.TransactionID
	return &tt
}

func (t Name) NewID() *Name {
	id := xid.New().String()
	t.ID = &id
	t.NewCreated()
	return &t
}

func (t Name) NewCreated() {
	t.CreatedAt = time.Now()
	t.NewUpdated()
}

func (t Name) NewUpdated() {
	t.UpdatedAt = time.Now()
}

func (t Name) Delete() {
	t.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

func (t Name) GetIDs() (*string, *string, string) {
	return t.TransactionID, t.ID, t.Type
}

func (t Name) GetBase() Base {
	return t.Base
}

func (t Name) GetBaseItem(base Base) *Name {
	return &Name{Base: base}
}

func (t Name) SetBase(base Base) *Name {
	t.Base = base
	return &t
}

func (t Name) Equals(n Name) bool {
	t.Base = Base{}
	n.Base = Base{}
	t.Changed = false
	n.Changed = false
	return reflect.DeepEqual(t, n)
}

func (t Name) IsEmpty() bool {
	return t.First == nil && t.Last == nil
}

func (t Name) IsValidToSubmit() bool {
	return t.NameRequest.First != nil && t.NameRequest.Last != nil
}

func (t Name) FillEmpty() Name {
	if t.Other == nil {
		t.Other = internal.StrPtr("")
	}

	return t
}

func (Dob) TableName() string {
	return "dob"
}

func (t Dob) CopyIDs() *Dob {
	var tt Dob
	tt.TransactionID = t.TransactionID
	return &tt
}

func (t Dob) NewID() *Dob {
	id := xid.New().String()
	t.ID = &id
	t.NewCreated()
	return &t
}

func (t Dob) NewCreated() {
	t.CreatedAt = time.Now()
	t.NewUpdated()
}

func (t Dob) NewUpdated() {
	t.UpdatedAt = time.Now()
}

func (t Dob) Delete() {
	t.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

func (t Dob) GetIDs() (*string, *string, string) {
	return t.TransactionID, t.ID, t.Type
}

func (t Dob) GetBase() Base {
	return t.Base
}

func (t Dob) GetBaseItem(base Base) *Dob {
	return &Dob{Base: base}
}

func (t Dob) SetBase(base Base) *Dob {
	t.Base = base
	return &t
}

func (t Dob) IsEmpty() bool {
	return t.Dob.Time.IsZero()
}

func (t Dob) IsValidToSubmit() bool {
	return t.Dob.Valid
}

func (t Dob) Equals(d Dob) bool {
	t.Base = Base{}
	d.Base = Base{}
	t.Changed = false
	d.Changed = false
	return t.Dob.ToString() == d.Dob.ToString()
}
