package grpc_api

import (
	"context"
	"errors"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/grpc_api/storage"
	"github.com/webitel/storage/model"
	"io"
	"net/http"
)

type file struct {
	ctrl *controller.Controller
}

func NewFileApi(api *controller.Controller) *file {
	return &file{api}
}

func (api *file) UploadFile(in storage.FileService_UploadFileServer) error {
	var chunk *storage.UploadFileRequest_Chunk

	res, gErr := in.Recv()
	if gErr != nil {

		return gErr
	}

	metadata, ok := res.Data.(*storage.UploadFileRequest_Metadata_)
	if !ok {
		// bad request
		return errors.New("bad metadata")
	}

	var fileRequest model.JobUploadFile
	fileRequest.DomainId = metadata.Metadata.DomainId
	fileRequest.Name = metadata.Metadata.Name
	fileRequest.MimeType = metadata.Metadata.MimeType
	fileRequest.Uuid = metadata.Metadata.Uuid

	pipeReader, pipeWriter := io.Pipe()

	go func(writer *io.PipeWriter) {
		for {
			res, gErr = in.Recv()
			if gErr != nil {
				//TODO
				break
			}

			if chunk, ok = res.Data.(*storage.UploadFileRequest_Chunk); !ok {
				//TODO
				break
			}

			if len(chunk.Chunk) == 0 {
				break
			}

			writer.Write(chunk.Chunk)
		}

		if gErr != nil {

		}

		writer.Close()

	}(pipeWriter)

	var err *model.AppError
	var publicUrl string
	if err = api.ctrl.UploadFileStream(pipeReader, &fileRequest); err != nil {
		return err
	}

	if publicUrl, err = api.ctrl.GeneratePreSignetResourceSignature(model.AnyFileRouteName, "download", fileRequest.Id, fileRequest.DomainId); err != nil {
		return err
	}

	return in.SendAndClose(&storage.UploadFileResponse{
		FileId:  fileRequest.Id,
		Code:    storage.UploadStatusCode_Ok,
		FileUrl: publicUrl,
	})
}

func (api *file) UploadFileUrl(ctx context.Context, in *storage.UploadFileUrlRequest) (*storage.UploadFileUrlResponse, error) {
	var err *model.AppError
	var publicUrl string

	if in.Url == "" || in.DomainId == 0 || in.Name == "" {
		return nil, errors.New("bad request")
	}

	res, httpErr := http.Get(in.GetUrl())
	if httpErr != nil {
		return nil, httpErr
	}

	var fileRequest model.JobUploadFile
	fileRequest.DomainId = in.GetDomainId()
	fileRequest.Name = in.GetName()
	fileRequest.MimeType = res.Header.Get("Content-Type")
	fileRequest.Uuid = in.GetUuid()
	if fileRequest.Uuid == "" {
		fileRequest.Uuid = model.NewId() // bad request ?
	}

	if err = api.ctrl.UploadFileStream(res.Body, &fileRequest); err != nil {
		return nil, err
	}

	if publicUrl, err = api.ctrl.GeneratePreSignetResourceSignature(model.AnyFileRouteName, "download", fileRequest.Id, fileRequest.DomainId); err != nil {
		return nil, err
	}

	return &storage.UploadFileUrlResponse{
		FileId:   fileRequest.Id,
		Code:     storage.UploadStatusCode_Ok,
		FileUrl:  publicUrl,
		MimeType: fileRequest.MimeType,
	}, nil
}
