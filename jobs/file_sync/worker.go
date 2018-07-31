package file_sync

import (
	"fmt"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/jobs"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"github.com/webitel/storage/utils"
	"time"
)

const LIMIT_FATCH_COUNT = 100

type Worker struct {
	name      string
	stop      chan bool
	stopped   chan bool
	jobs      chan model.Job
	jobServer *jobs.JobServer
	app       *app.App
}

func (m *SyncFilesJobInterfaceImpl) MakeWorker() model.Worker {
	worker := Worker{
		name:      "Sync Files",
		stop:      make(chan bool, 1),
		stopped:   make(chan bool, 1),
		jobs:      make(chan model.Job),
		jobServer: m.App.Jobs,
		app:       m.App,
	}

	return &worker
}

func (worker *Worker) Run() {
	mlog.Debug("Worker started", mlog.String("worker", worker.name))

	defer func() {
		mlog.Debug("Worker finished", mlog.String("worker", worker.name))
		worker.stopped <- true
	}()

	for {
		select {
		case <-worker.stop:
			mlog.Debug("Worker received stop signal", mlog.String("worker", worker.name))
			return
		case job := <-worker.jobs:
			mlog.Debug("Worker received a new candidate job.", mlog.String("worker", worker.name))
			worker.DoJob(&job)
		}
	}
}

func (worker *Worker) Stop() {
	mlog.Debug("Worker stopping", mlog.String("worker", worker.name))
	worker.stop <- true
	<-worker.stopped
}

func (worker *Worker) JobChannel() chan<- model.Job {
	return worker.jobs
}

func (worker *Worker) DoJob(job *model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		mlog.Info("Worker experienced an error while trying to claim job",
			mlog.String("worker", worker.name),
			mlog.String("job_id", job.Id),
			mlog.String("error", err.Error()))
		return
	} else if !claimed {
		return
	}

	for {
		select {
		case <-time.After(100 * time.Millisecond):
		loop:
			start := model.GetMillis()
			err, done := worker.removeFiles(job)
			job.Progress = model.GetMillis() - start
			if err != nil {
				mlog.Error("Worker: Failed to run remove files", mlog.String("worker", worker.name), mlog.String("job_id", job.Id), mlog.String("error", err.Error()))
				worker.setJobError(job, err)
				return
			} else if done {
				mlog.Info("Worker: Job is complete", mlog.String("worker", worker.name), mlog.String("job_id", job.Id))
				worker.setJobSuccess(job)
			} else {
				if err := worker.app.Jobs.UpdateInProgressJobData(job); err != nil {
					mlog.Error("Worker: Failed to update remove files status data for job", mlog.String("worker", worker.name), mlog.String("job_id", job.Id), mlog.String("error", err.Error()))
					worker.setJobError(job, err)
					return
				}
				goto loop
			}
			return

		case <-worker.stop:
			mlog.Debug("Worker: Job has been canceled via Worker Stop", mlog.String("worker", worker.name), mlog.String("job_id", job.Id))
			worker.setJobCanceled(job)
			return
		}
	}
}

func (worker *Worker) setJobSuccess(job *model.Job) {
	if err := worker.app.Jobs.SetJobSuccess(job); err != nil {
		mlog.Error("Worker: Failed to set success for job", mlog.String("worker", worker.name), mlog.String("job_id", job.Id), mlog.String("error", err.Error()))
		worker.setJobError(job, err)
	}
}

func (worker *Worker) setJobError(job *model.Job, appError *model.AppError) {
	if err := worker.app.Jobs.SetJobError(job, appError); err != nil {
		mlog.Error("Worker: Failed to set job error", mlog.String("worker", worker.name), mlog.String("job_id", job.Id), mlog.String("error", err.Error()))
	}
}

func (worker *Worker) setJobCanceled(job *model.Job) {
	if err := worker.app.Jobs.SetJobCanceled(job); err != nil {
		mlog.Error("Worker: Failed to mark job as canceled", mlog.String("worker", worker.name), mlog.String("job_id", job.Id), mlog.String("error", err.Error()))
	}
}

func (worker *Worker) removeFiles(job *model.Job) (*model.AppError, bool) {
	res := <-worker.app.Store.File().FetchDeleted(LIMIT_FATCH_COUNT)
	if res.Err != nil {
		return res.Err, false
	}
	var backend utils.FileBackend
	var err *model.AppError
	var resultDel store.StoreResult
	var removed = 0

	defer func() {
		job.Data["removed"] = fmt.Sprintf("%d", removed)
	}()

	data := res.Data.([]*model.FileWithProfile)

	for _, f := range data {
		backend, err = worker.app.GetFileBackendStore(f.ProfileId, f.ProfileUpdatedAt)
		if err != nil {
			panic(err)
		}
		err = backend.Remove(f)
		if err != nil {
			panic(err.Error())
		}

		resultDel = <-worker.app.Store.File().DeleteById(f.Id)
		if res.Err != nil {
			panic(resultDel.Err)
		}
		mlog.Debug(fmt.Sprintf("File %s[%d] removed", f.GetStoreName(), f.Id))
		removed++
	}
	return nil, len(data) == 0
}
