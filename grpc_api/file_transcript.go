package grpc_api

import (
	"context"

	"github.com/webitel/storage/model"

	"github.com/webitel/protos/storage"
	"github.com/webitel/storage/controller"
)

type fileTranscript struct {
	ctrl *controller.Controller
}

func NewFileTranscriptApi(api *controller.Controller) *fileTranscript {
	return &fileTranscript{
		ctrl: api,
	}
}

func (api *fileTranscript) CreateFileTranscript(ctx context.Context, in *storage.StartFileTranscriptRequest) (*storage.StartFileTranscriptResponse, error) {
	session, err := api.ctrl.GetSessionFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	ops := &model.TranscriptOptions{
		Locale: in.GetLocale(),
	}

	if in.GetProfile().GetId() > 0 {
		ops.ProfileId = model.NewInt(int(in.GetProfile().GetId()))
	}

	list, err := api.ctrl.TranscriptFiles(session, in.GetFileId(), ops)
	if err != nil {
		return nil, err
	}

	res := &storage.StartFileTranscriptResponse{
		Items: make([]*storage.StartFileTranscriptResponse_TranscriptJob, 0, len(list)),
	}

	for _, v := range list {
		res.Items = append(res.Items, &storage.StartFileTranscriptResponse_TranscriptJob{
			Id:        v.Id,
			FileId:    v.FileId,
			CreatedAt: v.CreatedAt,
		})
	}

	return res, nil
}
