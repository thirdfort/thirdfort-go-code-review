package repositories

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
	"gorm.io/gorm"
)

// Create new Actor with internal ID
func (s *Store) CreateActor(ctx context.Context, actorData map[string]string) (*models.Actor, error) {
	id := xid.New().String()
	actor := models.Actor{
		ID: &id,
		Shared: models.Shared{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if actorData != nil {
		name, ok := actorData["name"]
		if ok {
			actor.Name = name
		}
		mobile, ok := actorData["mobile"]
		if ok {
			actor.Mobile = mobile
		}
		email, ok := actorData["email"]
		if ok {
			actor.Email = email
		}
	}

	res := s.db.WithContext(ctx).Create(actor)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	return s.GetActor(ctx, &models.Actor{ID: &id})
}

func (s *Store) UpdateActor(ctx context.Context, actorID *string, actorData map[string]string) (*models.Actor, error) {
	actor := models.Actor{
		ID: actorID,
		Shared: models.Shared{
			UpdatedAt: time.Now(),
		},
	}

	s.fillActorWithData(&actor, actorData)

	res := s.db.WithContext(ctx).Save(actor)
	if res.Error != nil {
		s.logError(ctx, res)
		return nil, res.Error
	}

	return s.GetActor(ctx, &models.Actor{ID: actorID})
}

func (s *Store) fillActorWithData(actor *models.Actor, actorData map[string]string) {
	if actorData != nil {
		name, ok := actorData["name"]
		if ok {
			actor.Name = name
		}
		mobile, ok := actorData["mobile"]
		if ok {
			actor.Mobile = mobile
		}
		email, ok := actorData["email"]
		if ok {
			actor.Email = email
		}
	}
}

func (s *Store) GetActor(ctx context.Context, actor *models.Actor) (*models.Actor, error) {
	var ret models.Actor

	q := s.db.WithContext(ctx).Model(actor).Where(models.Actor{ID: actor.ID})

	res := q.First(&ret)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrNotFound
		}
		s.logError(ctx, res)
		return nil, res.Error
	}

	return &ret, nil
}
