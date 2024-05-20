package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/pkg/errors"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"github.com/thirdfort/thirdfort-go-code-review/internal/repositories"
	"gorm.io/gorm"
)

func (s *Service) GetPersonalInformation(ctx context.Context, item *models.PersonalInformation) (*models.PersonalInformation, error) {
	return s.DataStore.GetPersonalInformation(ctx, item)
}

func (s *Service) PatchPersonalInformation(ctx context.Context, item *models.PersonalInformation) (*models.PersonalInformation, error) {
	if item.Status == internal.StatusCancelled {
		return s.cancelPersonalInformation(ctx, item)
	}

	if item.Name.IsEmpty() && item.Dob.IsEmpty() && item.Status != internal.StatusCompleted {
		return nil, internal.ErrBadRequest
	}

	base := models.Base{
		TransactionID: item.TransactionID,
		Shared: models.Shared{
			DeletedAt: gorm.DeletedAt{},
		},
	}

	store := s.DataStore.GetStore()

	// We can only store the data in the database.
	if !item.Name.IsEmpty() {
		// get the already existing value
		base.Type = internal.TypePersonalInformationName
		originalName, err := repositories.New[models.Name](store).Get(ctx, *&models.Name{Base: base}, nil)
		if err != nil {
			return nil, err
		}

		if slices.Contains([]string{internal.StatusCompleted, internal.StatusCancelled}, originalName.Status) {
			return nil, internal.ErrBadRequest
		}

		// add a space at the last name if it's only a single letter
		if len(*item.Name.Last) == 1 {
			*item.Name.Last = fmt.Sprintf("%s ", *item.Name.Last)
		}

		originalName.NameRequest = item.Name.NameRequest
		originalName.Status = item.Status
		if item.Status != internal.StatusCompleted {
			originalName.Status = internal.StatusInProgress
			storedName, err := repositories.New[models.Name](store).UpdateChanges(ctx, *originalName, nil)
			if err != nil {
				return nil, err
			}

			item.Name = storedName
		}

		item.ExpectationID = originalName.ExpectationID
		item.ReasonCode = originalName.ReasonCode
	}

	if !item.Dob.IsEmpty() {
		// get the already existing value
		base.Type = internal.TypePersonalInformationDob

		originalDob, err := repositories.New[models.Dob](store).Get(ctx, *&models.Dob{Base: base}, nil)
		if err != nil {
			return nil, err
		}

		if slices.Contains([]string{internal.StatusCompleted, internal.StatusCancelled}, originalDob.Status) {
			return nil, internal.ErrBadRequest
		}

		originalDob.DobRequest = item.Dob.DobRequest
		originalDob.Status = item.Status
		if item.Status != internal.StatusCompleted {
			originalDob.Status = internal.StatusInProgress
			storedDob, err := repositories.New[models.Dob](store).UpdateChanges(ctx, *originalDob, nil)
			if err != nil {
				return nil, err
			}

			item.Dob = storedDob
		}

		item.ExpectationID = originalDob.ExpectationID
		item.ReasonCode = originalDob.ReasonCode
	}

	// If we aren't updating status we return the otherwise updated item
	if item.Status != internal.StatusCompleted {
		// Get the value from the database
		return s.DataStore.GetPersonalInformation(ctx, item)
	}

	// After this step it's more or less the same as the PutPersonalInformation
	return nil, nil
}

func (s *Service) cancelPersonalInformation(ctx context.Context, item *models.PersonalInformation) (*models.PersonalInformation, error) {
	base := models.Base{
		TransactionID: item.TransactionID,
		Shared: models.Shared{
			DeletedAt: gorm.DeletedAt{},
		},
	}

	store := s.DataStore.GetStore()

	cancelBase := models.Base{
		TransactionID: item.TransactionID,
		Status:        internal.StatusCancelled,
		Shared: models.Shared{
			DeletedAt: gorm.DeletedAt{},
		},
	}

	base.Type = internal.TypePersonalInformationName
	originalName, err := repositories.New[models.Name](store).Get(ctx, *&models.Name{Base: base}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Get: Name")
	}

	base.Type = internal.TypePersonalInformationDob
	originalDob, err := repositories.New[models.Dob](store).Get(ctx, *&models.Dob{Base: base}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Get: Dob")
	}

	// check if the original resource has reason code
	reasonCode := models.GetReasonCodeForTask(originalName.ReasonCode, originalDob.ReasonCode)
	if reasonCode == "" {
		return nil, internal.ErrBadRequest
	}

	// now update the originalName, originalDob as cancelled.
	cancelBase.Type = internal.TypePersonalInformationName
	originalName, err = repositories.New[models.Name](store).UpdateChanges(ctx, models.Name{Base: cancelBase}, nil)
	if err != nil {
		return nil, err
	}

	cancelBase.Type = internal.TypePersonalInformationDob
	originalDob, err = repositories.New[models.Dob](store).UpdateChanges(ctx, models.Dob{Base: cancelBase}, nil)
	if err != nil {
		return nil, err
	}

	return s.GetPersonalInformation(ctx, item)
}
