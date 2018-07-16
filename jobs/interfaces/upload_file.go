package interfaces

import (
	"github.com/webitel/storage/model"
)

type UploadRecordingsFilesJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
