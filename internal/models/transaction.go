package models

import (
	"time"

	"github.com/jinzhu/copier"
)

type Transaction struct {
	Shared
	ID           *string `json:"id" path:"txID" gorm:"primaryKey;column:id" description:"Transaction ID" example:"cq1qh5c23amg0302nqv0"`
	EventID      *string `json:"event_id" gorm:"column:event_id" description:"Event ID associated with transaction" example:"cpkz6tb23amg03015m70"`
	Name         *string `json:"name" description:"Transaction's name" example:"My Transaction"`
	TenantName   *string `json:"tenant_name" gorm:"column:tenant_name" description:"Name of tenant" example:"Thirdfort Limited"`
	ConsumerName *string `json:"consumer_name" gorm:"column:actor_name" description:"Consumer's name" example:"Bob Bobson"`
	Ref          *string `json:"ref" gorm:"column:ref" description:"Transaction reference" example:"MyTransaction01"`
	Status       string  `json:"status" description:"Status of the transaction" example:"not_started" enum:"not_started,in_progress,completed,cancelled"`
	ActorID      *string `json:"-" gorm:"primaryKey;column:actor_id"`
	Metadata     *string `json:"metadata" gorm:"type:jsonb" description:"Transaction Metadata from PA"`
}

type TransactionResponse struct {
	ID           string    `json:"id" path:"txID" description:"Transaction ID" example:"cq1qh5c23amg0302nqv0"`
	EventID      string    `json:"event_id" description:"Event ID associated with transaction" example:"cpkz6tb23amg03015m70"`
	Name         string    `json:"name" description:"Transaction's name" example:"My Transaction"`
	TenantName   string    `json:"tenant_name" description:"Name of tenant" example:"Thirdfort Limited"`
	ConsumerName string    `json:"consumer_name" description:"Consumer's name" example:"Bob Bobson"`
	Ref          string    `json:"ref" description:"Transaction reference" example:"MyTransaction01"`
	Status       string    `json:"status" description:"Status of the transaction" example:"not_started" enum:"not_started,in_progress,completed,cancelled"`
	CreatedAt    time.Time `json:"created_at" description:"Transaction creation time" example:"2021-07-01T12:00:00Z"`
}

func (a *Transaction) ToResponse() *TransactionResponse {
	var data TransactionResponse

	copier.Copy(&data, a)
	data.Status = MapPaStatus(a.Status)

	return &data
}

func (*Transaction) TableName() string {
	return "transaction"
}

type TransactionStatus struct {
	ID     string `json:"-" path:"txID" description:"Transaction ID" example:"cq1qh5c23amg0302nqv0"`
	Status string `json:"status" description:"Status of the transaction" example:"accepted" enum:"accepted,rejected"`
}
