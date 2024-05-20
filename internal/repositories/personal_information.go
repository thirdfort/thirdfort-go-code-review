package repositories

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"gorm.io/gorm"
)

func (s *Store) GetPersonalInformation(ctx context.Context, task *models.PersonalInformation) (*models.PersonalInformation, error) {
	base := models.Base{
		TransactionID: task.TransactionID,
		Shared: models.Shared{
			DeletedAt: gorm.DeletedAt{},
		},
	}

	name := models.Name{Base: base}
	name.Type = internal.TypePersonalInformationName
	dob := models.Dob{Base: base}
	dob.Type = internal.TypePersonalInformationDob

	var storedName models.Name
	var storedDob models.Dob
	q := s.db.WithContext(ctx).First(&name).Scan(&storedName)
	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrNotFound
		}
		s.logError(ctx, q)
		return nil, q.Error
	}

	q = s.db.WithContext(ctx).First(&dob).Scan(&storedDob)
	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrNotFound
		}
		s.logError(ctx, q)
		return nil, q.Error
	}

	return &models.PersonalInformation{
		Name: &storedName,
		Dob:  &storedDob,
	}, nil
}

func (s *Store) PutPersonalInformation(ctx context.Context, task *models.PersonalInformation) (*models.PersonalInformation, error) {
	name := *task.Name
	dob := *task.Dob
	oldTask, err := s.GetPersonalInformation(ctx, task)
	if err != nil {
		return nil, err
	}

	name.Base = oldTask.Name.Base
	dob.Base = oldTask.Dob.Base
	if task.Status != "" {
		name.Base.Status = task.Status
		dob.Base.Status = task.Status
	}

	q := s.db.WithContext(ctx).Save(&name)
	if q.Error != nil {
		return nil, errors.Wrap(q.Error, "PutPersonalInformation: Name")
	}

	q = s.db.WithContext(ctx).Save(&dob)
	if q.Error != nil {
		return nil, errors.Wrap(q.Error, "PutPersonalInformation: Dob")
	}

	return s.GetPersonalInformation(ctx, task)
}
