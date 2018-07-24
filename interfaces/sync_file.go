package interfaces

import (
	"github.com/webitel/storage/model"
)

type SyncFilesJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
