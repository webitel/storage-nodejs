package grpc_api

import (
	"context"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/grpc_api/storage"
	"github.com/webitel/storage/model"
)

type backendProfiles struct {
	ctrl *controller.Controller
}

func NewBackendProfileApi(api *controller.Controller) *backendProfiles {
	return &backendProfiles{api}
}

func (api *backendProfiles) CreateBackendProfile(ctx context.Context, in *storage.CreateBackendProfileRequest) (*storage.BackendProfile, error) {
	session, err := api.ctrl.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	profile := &model.FileBackendProfile{
		Name:        in.GetName(),
		Description: in.GetDescription(),
		ExpireDay:   int(in.GetExpireDays()),
		Priority:    int(in.GetPriority()),
		Disabled:    in.GetDisabled(),
		MaxSizeMb:   int(in.GetMaxSize()),
		Properties:  nil,
		Type: model.Lookup{
			Id: int(in.GetType().GetId()),
		},
	}

	profile, err = api.ctrl.CreateBackendProfile(session, profile)
	if err != nil {
		return nil, err
	}

	return toGrpcProfile(profile), nil
}

func toGrpcProfile(src *model.FileBackendProfile) *storage.BackendProfile {
	return &storage.BackendProfile{
		Id:        src.Id,
		CreatedAt: src.CreatedAt,
		CreatedBy: &storage.Lookup{
			Id:   int64(src.CreatedBy.Id),
			Name: src.CreatedBy.Name,
		},
		UpdatedAt: src.UpdatedAt,
		UpdatedBy: &storage.Lookup{
			Id:   int64(src.UpdatedBy.Id),
			Name: src.UpdatedBy.Name,
		},
		DataSize:   int64(src.DataSize),
		DataCount:  src.DataCount,
		Name:       src.Name,
		ExpireDays: int32(src.ExpireDay),
		MaxSize:    int64(src.MaxSizeMb),
		Priority:   int32(src.Priority),
		Type: &storage.Lookup{
			Id:   int64(src.Type.Id),
			Name: src.Type.Name,
		},
		Properties:  nil,
		Description: src.Description,
		Disabled:    src.Disabled,
	}
}
