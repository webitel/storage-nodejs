package grpc_api

import (
	"context"
	"github.com/webitel/engine/grpc_api/engine"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/grpc_api/storage"
	"github.com/webitel/storage/model"
)

type media struct {
	ctrl *controller.Controller
}

func NewMediaApi(api *controller.Controller) *media {
	return &media{api}
}

func (api *media) SearchMediaFile(ctx context.Context, in *storage.SearchMediaFileRequest) (*storage.ListMedia, error) {
	session, err := api.ctrl.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	var list []*model.MediaFile
	var endOfList bool

	req := &model.SearchMediaFile{
		ListRequest: model.ListRequest{
			Q:       in.GetQ(),
			Page:    int(in.GetPage()),
			PerPage: int(in.GetSize()),
		},
	}

	list, endOfList, err = api.ctrl.SearchMediaFile(session, in.GetDomainId(), req)

	if err != nil {
		return nil, err
	}

	items := make([]*storage.MediaFile, 0, len(list))
	for _, v := range list {
		items = append(items, toGrpcMediaFile(v))
	}
	return &storage.ListMedia{
		Next:  !endOfList,
		Items: items,
	}, nil
}

func (api *media) ReadMediaFile(ctx context.Context, in *storage.ReadMediaFileRequest) (*storage.MediaFile, error) {
	session, err := api.ctrl.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var file *model.MediaFile

	file, err = api.ctrl.GetMediaFile(session, in.GetDomainId(), int(in.GetId()))
	if err != nil {
		return nil, err
	}

	return toGrpcMediaFile(file), nil
}

func (api *media) DeleteMediaFile(ctx context.Context, in *storage.DeleteMediaFileRequest) (*storage.MediaFile, error) {
	session, err := api.ctrl.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	var file *model.MediaFile

	file, err = api.ctrl.DeleteMediaFile(session, in.GetDomainId(), int(in.GetId()))
	if err != nil {
		return nil, err
	}

	return toGrpcMediaFile(file), nil
}

func toGrpcMediaFile(src *model.MediaFile) *storage.MediaFile {
	return &storage.MediaFile{
		Id:        src.Id,
		CreatedAt: src.CreatedAt,
		CreatedBy: &engine.Lookup{
			Id:   int64(src.CreatedBy.Id),
			Name: src.CreatedBy.Name,
		},
		UpdatedAt: src.UpdatedAt,
		UpdatedBy: &engine.Lookup{
			Id:   int64(src.UpdatedBy.Id),
			Name: src.UpdatedBy.Name,
		},
		Name:     src.Name,
		Size:     src.Size,
		MimeType: src.MimeType,
	}
}
