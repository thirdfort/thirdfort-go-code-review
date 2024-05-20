package models

type Actor struct {
	Shared
	ID           *string       `json:"id" gorm:"primaryKey;column:id"`
	Name         string        `json:"name"`
	Mobile       string        `json:"mobile"`
	Email        string        `json:"email"`
	Fingerprints []Fingerprint `json:"fingerprint"`
}

func (*Actor) TableName() string {
	return "actor"
}

type Fingerprint struct {
	Shared
	Fingerprint *string `json:"id" gorm:"primaryKey;column:fingerprint"`
	ActorID     *string `json:"actor_id" gorm:"column:actor_id"`
	PaActorID   *string `json:"pa_actor_id" gorm:"column:pa_actor_id"`
}

func (*Fingerprint) TableName() string {
	return "fingerprint"
}
