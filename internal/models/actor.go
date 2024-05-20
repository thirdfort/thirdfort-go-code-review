package models

type Actor struct {
	Shared
	ID     *string `json:"id" gorm:"primaryKey;column:id"`
	Name   string  `json:"name"`
	Mobile string  `json:"mobile"`
	Email  string  `json:"email"`
}

func (*Actor) TableName() string {
	return "actor"
}
