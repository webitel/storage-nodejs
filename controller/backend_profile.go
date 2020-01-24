package controller

import (
	"github.com/webitel/engine/auth_manager"
	"github.com/webitel/storage/model"
)

func (c *Controller) CreateBackendProfile(session *auth_manager.Session, profile *model.FileBackendProfile) (*model.FileBackendProfile, *model.AppError) {
	var err *model.AppError
	permission := session.GetPermission(model.PERMISSION_SCOPE_BACKEND_PROFILE_ROUTING)
	if !permission.CanCreate() {
		return nil, c.app.MakePermissionError(session, permission, auth_manager.PERMISSION_ACCESS_CREATE)
	}

	profile.DomainRecord = model.DomainRecord{
		Id:        0,
		DomainId:  session.Domain(profile.DomainId),
		CreatedAt: model.GetMillis(),
		CreatedBy: model.Lookup{
			Id: int(session.UserId),
		},
		UpdatedAt: model.GetMillis(),
		UpdatedBy: model.Lookup{
			Id: int(session.UserId),
		},
	}

	if err = profile.IsValid(); err != nil {
		return nil, err
	}

	return c.app.CreateFileBackendProfile(profile)
}
