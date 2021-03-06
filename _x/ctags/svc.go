package ctags

import (
	"context"

	"github.com/tusharsoni/copper/cerror"
)

// Svc provides methods to manage tags on entities
type Svc interface {
	AddTag(ctx context.Context, tag, entityID string) error
	RemoveTag(ctx context.Context, t, entityID string) error
	ListTags(ctx context.Context, entityID string) ([]string, error)
	HasTag(ctx context.Context, entityID, tag string) (bool, error)
	ListEntityIDs(ctx context.Context, tag string) ([]string, error)
}

type svcImpl struct {
	tags repo
}

func newSvcImpl(tags repo) Svc {
	return &svcImpl{
		tags: tags,
	}
}

func (s *svcImpl) AddTag(ctx context.Context, t, entityID string) error {
	return s.tags.Add(ctx, &tag{
		Tag:      t,
		EntityID: entityID,
	})
}

func (s *svcImpl) RemoveTag(ctx context.Context, t, entityID string) error {
	return s.tags.Delete(ctx, t, entityID)
}

func (s *svcImpl) ListTags(ctx context.Context, entityID string) ([]string, error) {
	tags, err := s.tags.FindByEntityID(ctx, entityID)
	if err != nil {
		return nil, cerror.New(err, "failed to find tags by entity id", map[string]interface{}{
			"entityID": entityID,
		})
	}

	t := make([]string, len(tags))
	for i, tag := range tags {
		t[i] = tag.Tag
	}

	return t, nil
}

func (s *svcImpl) HasTag(ctx context.Context, entityID, tag string) (bool, error) {
	tags, err := s.ListTags(ctx, entityID)
	if err != nil {
		return false, cerror.New(err, "failed to list tags", map[string]interface{}{
			"entityID": entityID,
		})
	}

	for _, t := range tags {
		if t == tag {
			return true, nil
		}
	}

	return false, nil
}

func (s *svcImpl) ListEntityIDs(ctx context.Context, tag string) ([]string, error) {
	tags, err := s.tags.FindByTag(ctx, tag)
	if err != nil {
		return nil, cerror.New(err, "failed to find tags by tag", map[string]interface{}{
			"tag": tag,
		})
	}

	entityIDs := make([]string, len(tags))
	for i, tag := range tags {
		entityIDs[i] = tag.EntityID
	}

	return entityIDs, nil
}
