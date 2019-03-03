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
	UserHasPermission(ctx context.Context, userUUID string, resource, action string) (bool, error)
	GrantPermissions(ctx context.Context, granteeID, resource string, action []string) error
	RevokePermission(ctx context.Context, granteeID, resource, action string) error

	CreateRole(ctx context.Context, name string) error
	AddUserToRole(ctx context.Context, userUUID, roleUUID string) error
}

type svcImpl struct {
	repo repo
}

func newSvcImpl(repo repo) Svc {
	return &svcImpl{
		repo: repo,
	}
}

func (s *svcImpl) UserHasPermission(ctx context.Context, userUUID string, resource, action string) (bool, error) {
	has, err := s.HasPermission(ctx, userUUID, resource, action)
	if err != nil {
		return false, cerror.New(err, "failed to check permission", map[string]string{
			"userUUID": userUUID,
			"resource": resource,
			"action":   action,
		})
	}

	if has {
		return true, nil
	}

	roles, err := s.repo.FindRolesForUserUUID(ctx, userUUID)
	if err != nil {
		return false, cerror.New(err, "failed to find roles for user", map[string]string{
			"userUUID": userUUID,
		})
	}

	for _, r := range roles {
		has, err := s.HasPermission(ctx, r.UUID, resource, action)
		if err != nil {
			return false, cerror.New(err, "failed to check permission", map[string]string{
				"roleUUID": r.UUID,
				"resource": resource,
				"action":   action,
			})
		}

		if has {
			return true, nil
		}
	}

	return false, nil
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

func (s *svcImpl) GrantPermissions(ctx context.Context, granteeID, resource string, actions []string) error {
	pUUID, err := uuid.NewRandom()
	if err != nil {
		return cerror.New(err, "failed to generate random uuid", nil)
	}

	for _, action := range actions {
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

func (s *svcImpl) CreateRole(ctx context.Context, name string) error {
	rUUID, err := uuid.NewRandom()
	if err != nil {
		return cerror.New(err, "failed to generate random uuid", nil)
	}

	r := role{
		UUID: rUUID.String(),
		Name: name,
	}

	err = s.repo.AddRole(ctx, &r)
	if err != nil {
		return cerror.New(err, "failed to upsert role", map[string]string{
			"uuid": rUUID.String(),
			"name": name,
		})
	}

	return nil
}

func (s *svcImpl) AddUserToRole(ctx context.Context, userUUID, roleUUID string) error {
	return s.repo.AddUserRole(ctx, userUUID, roleUUID)
}
