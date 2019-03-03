package cacl

import (
	"context"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"
)

// Svc provides methods to manage permissions for a grantee (user, role etc.), resource, and action combination.
// For example:
// 	Alice (user grantee) can write (action) to files/bio.txt (resource).
// 	Admin (role grantee) can write (action) to files/bio.txt (resource).
type Svc interface {
	HasPermission(ctx context.Context, granteeID, resource, action string) (bool, error)
	GivePermission(ctx context.Context, granteeID, resource, action string) error
	RevokePermission(ctx context.Context, granteeID, resource, action string) error
}

type svcImpl struct {
	repo repo
}

func newSvcImpl(repo repo) Svc {
	return &svcImpl{
		repo: repo,
	}
}

func (s *svcImpl) HasPermission(ctx context.Context, granteeID, resource, action string) (bool, error) {
	_, err := s.repo.GetPermissionForGrantee(ctx, granteeID, resource, action)
	if err != nil && cerror.Cause(err) != gorm.ErrRecordNotFound {
		return false, cerror.New(err, "failed to get permission for grantte", map[string]string{
			"granteeID": granteeID,
			"resource":  resource,
			"action":    action,
		})
	} else if err != nil && cerror.Cause(err) == gorm.ErrRecordNotFound {
		return false, nil
	}

	return true, nil
}

func (s *svcImpl) GivePermission(ctx context.Context, granteeID, resource, action string) error {
	pUUID, err := uuid.NewRandom()
	if err != nil {
		return cerror.New(err, "failed to generate random uuid", nil)
	}

	p := permission{
		UUID:      pUUID.String(),
		GranteeID: granteeID,
		Resource:  resource,
		Action:    action,
	}

	err = s.repo.AddPermission(ctx, &p)
	if err != nil {
		return cerror.New(err, "failed to upsert permission", map[string]string{
			"uuid":      pUUID.String(),
			"granteeID": granteeID,
			"resource":  resource,
			"action":    action,
		})
	}

	return nil
}

func (s *svcImpl) RevokePermission(ctx context.Context, granteeID, resource, action string) error {
	p, err := s.repo.GetPermissionForGrantee(ctx, granteeID, resource, action)
	if err != nil && cerror.Cause(err) != gorm.ErrRecordNotFound {
		return cerror.New(err, "failed to get permission for grantte", map[string]string{
			"granteeID": granteeID,
			"resource":  resource,
			"action":    action,
		})
	}

	err = s.repo.DeletePermission(ctx, p.UUID)
	if err != nil {
		return cerror.New(err, "failed to delete permission", map[string]string{
			"uuid": p.UUID,
		})
	}

	return nil
}
