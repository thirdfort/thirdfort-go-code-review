package service

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/thirdfort/go-slogctx"
	"github.com/thirdfort/thirdfort-go-code-review/internal"
	"github.com/thirdfort/thirdfort-go-code-review/internal/models"
)

func (s *Service) GetActor(ctx context.Context) (*models.Actor, error) {
	fingerprint := s.getFingerprintFromCtx(ctx)

	actor := s.cache.GetActor(fingerprint)
	if actor != nil {
		return actor, nil
	}

	fp, err := s.DataStore.GetFingerprint(ctx, &models.Fingerprint{Fingerprint: &fingerprint})
	if err != nil {
		slogctx.Warn(ctx, "GetFingerprint", slog.String("fingerprint", fingerprint), slog.Any("err", err))
		return nil, errors.Wrap(err, "GetActor: DataStore.GetFingerprint")
	}

	actor, err = s.DataStore.GetActor(ctx, &models.Actor{ID: fp.ActorID})
	if err != nil {
		return nil, errors.Wrap(err, "GetActor: DataStore.GetActor")
	}

	s.cache.SetActor(fingerprint, actor)

	return actor, nil
}

func (s *Service) getFingerprintFromCtx(ctx context.Context) string {
	valueMap := internal.GetContextHeaders(ctx)
	return internal.SafeGetValueFromMap(valueMap, "Device-Fingerprint")
}

// Get actor or create if not found based on 'device-fingerprint'
func (s *Service) getCreateActor(ctx context.Context, fingerprint string, actorData map[string]string) (*models.Actor, error) {
	fp := &models.Fingerprint{Fingerprint: &fingerprint}
	fp, err := s.DataStore.GetFingerprint(ctx, fp)
	if err != nil {
		if errors.Is(err, internal.ErrNotFound) {
			actor, err := s.DataStore.CreateActor(ctx, actorData)
			if err != nil {
				return nil, errors.Wrap(err, "getCreateActor: CreateActor")
			}
			fp, err := s.DataStore.CreateFingerprint(ctx,
				&models.Fingerprint{
					Fingerprint: &fingerprint,
					ActorID:     actor.ID,
				})
			if err != nil {
				return nil, errors.Wrap(err, "getCreateActor: CreateFingerprint")
			}

			actor.Fingerprints = append(actor.Fingerprints, *fp)

			return actor, nil
		} else {
			return nil, errors.Wrap(err, "getCreateActor: GetActor")
		}
	}

	actor, err := s.DataStore.GetActor(ctx, &models.Actor{ID: fp.ActorID})
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
