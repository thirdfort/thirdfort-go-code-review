package models

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/rs/xid"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"gorm.io/gorm"
)

type Document struct {
	Base
	Documents []Documents `json:"documents" gorm:"foreignKey:DocumentID;references:ID"`
}

func (a *Document) ToResponse() *DocumentResponse {
	var (
		data DocumentRequest
		exp  DocumentExpectation
	)
	copier.Copy(&data.Documents.Data, a)
	copier.Copy(&exp, a.Base)
	exp.Data = data.Documents.Data
	exp.Status = MapPaStatus(a.Base.Status)
	return &DocumentResponse{
		TaskFields: TaskFields{
			Status:     getCombinedStatus(a.Base.Status),
			ReasonCode: GetReasonCodeForTask(a.Base.ReasonCode),
			Type:       a.Base.Type,
			CreatedAt:  a.Base.CreatedAt,
		},
		DocumentExpectation: exp,
	}
}

type DocumentRequest struct {
	PathTxID
	Status    string        `json:"status" enum:"not_started,in_progress,completed,cancelled,in_review"`
	Documents DocumentsData `json:"documents"`
}

type DocumentsData struct {
	Data DocumentFields `json:"data" description:"Documents data"`
}

type DocumentFields struct {
	Documents []Documents `json:"documents" description:"Documents fields"`
}

func (p DocumentRequest) ToModel(subtype string) *Document {
	base := Base{TransactionID: p.TxID, Type: subtype, Status: p.Status}
	return &Document{
		Base:      base,
		Documents: p.Documents.Data.Documents,
	}
}

type DocumentResponse struct {
	TaskFields
	DocumentExpectation `json:"document" description:"Expectation metadata"`
}

type DocumentExpectation struct {
	ExpectationFields
	Data DocumentFields `json:"data" description:"Task data"`
}

type Documents struct {
	Shared
	DocumentID *string `json:"-" gorm:"column:document_id;primaryKey"`
	ID         *string `json:"id" description:"Document's ID" example:"cq1qh5c23amg0302nqv0" gorm:"column:id;primaryKey"`
	Type       *string `json:"type" description:"Document's type" validate:"required" enum:"bank-statement,document,divorce,face-image,identity:driving-licence,identity:national-identity-card,identity:passport,identity:uk-biometric-residence-permit,identity:voter-id,inheritance,mortgage,poa,poo,sale-assets,savings,selfie" example:"bank-statement"`
}

func (Documents) TableName() string {
	return "documents"
}

func (Document) TableName() string {
	return "document"
}

func (t Document) CopyIDs() *Document {
	var tt Document
	tt.TransactionID = t.TransactionID
	tt.ID = t.ID
	return &tt
}

func (t Document) NewID() *Document {
	id := xid.New().String()
	t.ID = &id
	t.NewCreated()
	return &t
}

func (t Document) NewCreated() {
	t.CreatedAt = time.Now()
	t.NewUpdated()
}

func (t Document) NewUpdated() {
	t.UpdatedAt = time.Now()
}

func (t Document) Delete() {
	t.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
}

func (t Document) GetIDs() (*string, *string, string) {
	return t.TransactionID, t.ID, t.Type
}

func (t Document) GetBase() Base {
	return t.Base
}

func (t Document) SetBase(base Base) *Document {
	t.Base = base
	return &t
}

func (t Document) GetBaseItem(base Base) *Document {
	return &Document{Base: base}
}

func (t Document) IsEmpty() bool {
	if t.Status == internal.StatusCancelled {
		return false
	}

	return len(t.Documents) == 0
}

func (t Document) IsValidToSubmit() bool {
	if t.Status == internal.StatusCompleted && (t.Documents == nil || len(t.Documents) == 0) {
		return false
	}

	return true
}
