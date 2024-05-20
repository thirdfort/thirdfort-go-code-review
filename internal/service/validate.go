package service

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/thirdfort/go-types/v2"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
)

const (
	DobMinAge = 18
	DobMaxAge = 150
)

func (s *Service) Validate(item any, typeName string) (any, error) {
	var (
		err    error
		paType any
	)
	switch typeName {
	case "address":
		var paAddr types.Address
		copier.Copy(&paAddr, item)
		err = paAddr.Validate()
		paType = paAddr
	default:
		err = s.validator.Struct(item)
		paType = item
	}

	return paType, err
}

func IsValidDateFormat(str string) error {
	_, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return errors.New("Date format incorrect, needs to be YYYY-MM-DD")
	}

	return nil
}

func IsValidAge(str string) error {
	dob, err := time.Parse(internal.DateLayout, str)
	if err != nil {
		return errors.New("Could not parse Dob, needs to be YYYY-MM-DD")
	}

	if time.Since(dob) < internal.TimeYear*DobMinAge || time.Since(dob) >= internal.TimeYear*DobMaxAge {
		return errors.New("Invalid age, needs to be >18 and <120")
	}

	return nil
}
