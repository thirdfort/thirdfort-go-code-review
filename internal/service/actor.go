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

// Get actor or create if not found
func (s *Service) getCreateActor(ctx context.Context, id *string, actorData map[string]string) (*models.Actor, error) {

	actor, err := s.DataStore.GetActor(ctx, &models.Actor{ID: id})
	if err != nil {
		return nil, errors.Wrap(err, "getCreateActor:GetActor - could not match actor with fingerprint")
	}

	// We only have partial actor data from events, update
	if actor.Mobile == "" {
		actor, err = s.DataStore.UpdateActor(ctx, actor.ID, actorData)
		if err != nil {
			return nil, errors.Wrap(err, "getCreateActor: CreateActor")
		}
	}

	return actor, nil
}
