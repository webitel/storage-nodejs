package interfaces

import (
	"github.com/webitel/storage/model"
)

type JobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}

type SyncFilesJobInterface interface {
	MakeWorker() model.Worker
	MakeScheduler() model.Scheduler
}
