package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thirdfort/thirdfort-go-code-review/internal"

	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func (s *Service) GetActor(ctx context.Context) (*models.Actor, error) {
	id := s.getIDFromCtx(ctx)
	actor, err := s.DataStore.GetActor(ctx, &models.Actor{ID: &id})
	if err != nil {
		return nil, errors.Wrap(err, "GetActor: DataStore.GetActor")
	}

	s.cache.SetActor(id, actor)

	return actor, nil
}

func (s *Service) getIDFromCtx(ctx context.Context) string {
	valueMap := internal.GetContextHeaders(ctx)
	return internal.SafeGetValueFromMap(valueMap, "ID")
}
