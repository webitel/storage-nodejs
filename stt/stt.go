package stt

import (
	"context"

	"github.com/webitel/storage/model"
)

type Stt interface {
	Transcript(ctx context.Context, fileUri, locale string) (model.FileTranscript, error)
}
