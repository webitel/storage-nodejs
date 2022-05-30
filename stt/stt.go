package stt

import (
	"context"

	"github.com/webitel/storage/model"
)

type Stt interface {
	Transcript(ctx context.Context, id int64, fileUri, locale string) (model.FileTranscript, error)
	Callback(req map[string]interface{}) error
}
