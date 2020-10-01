package grpc_api

import (
	"errors"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/grpc_api/storage"
	"github.com/webitel/storage/model"
	"io"
)

type file struct {
	ctrl *controller.Controller
}

func NewFileApi(api *controller.Controller) *file {
	return &file{api}
}

func (api *file) UploadFile(in storage.FileService_UploadFileServer) error {
	var chunk *storage.UploadFileRequest_Chunk

	res, err := in.Recv()
	if err != nil {

		return err
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
			res, err = in.Recv()
			if err != nil {
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

		if err != nil {

		}

		writer.Close()

	}(pipeWriter)

	if err := api.ctrl.UploadFileStream(pipeReader, &fileRequest); err != nil {
		return err
	}

	return in.SendAndClose(&storage.UploadFileResponse{
		FileId: fileRequest.Id,
		Code:   0,
	})
}
