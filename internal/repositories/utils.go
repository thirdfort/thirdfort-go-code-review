package repositories

import (
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func MapExpectationData(data json.RawMessage) (*datatypes.JSON, error) {
	var out datatypes.JSON

	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func DeletedAt() gorm.DeletedAt {
	return gorm.DeletedAt{}
}
