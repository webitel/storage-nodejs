package file_sync

import (
	"fmt"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/jobs"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"github.com/webitel/storage/utils"
	"github.com/webitel/wlog"
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
	wlog.Debug("Worker started", wlog.String("worker", worker.name))

	defer func() {
		wlog.Debug("Worker finished", wlog.String("worker", worker.name))
		worker.stopped <- true
	}()

	for {
		select {
		case <-worker.stop:
			wlog.Debug("Worker received stop signal", wlog.String("worker", worker.name))
			return
		case job := <-worker.jobs:
			wlog.Debug("Worker received a new candidate job.", wlog.String("worker", worker.name))
			worker.DoJob(&job)
		}
	}
}

func (worker *Worker) Stop() {
	wlog.Debug("Worker stopping", wlog.String("worker", worker.name))
	worker.stop <- true
	<-worker.stopped
}

func (worker *Worker) JobChannel() chan<- model.Job {
	return worker.jobs
}

func (worker *Worker) DoJob(job *model.Job) {
	if claimed, err := worker.jobServer.ClaimJob(job); err != nil {
		wlog.Info("Worker experienced an error while trying to claim job",
			wlog.String("worker", worker.name),
			wlog.String("job_id", job.Id),
			wlog.String("error", err.Error()))
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
				wlog.Error("Worker: Failed to run remove files", wlog.String("worker", worker.name), wlog.String("job_id", job.Id), wlog.String("error", err.Error()))
				worker.setJobError(job, err)
				return
			} else if done {
				wlog.Info("Worker: Job is complete", wlog.String("worker", worker.name), wlog.String("job_id", job.Id))
				worker.setJobSuccess(job)
			} else {
				if err := worker.app.Jobs.UpdateInProgressJobData(job); err != nil {
					wlog.Error("Worker: Failed to update remove files status data for job", wlog.String("worker", worker.name), wlog.String("job_id", job.Id), wlog.String("error", err.Error()))
					worker.setJobError(job, err)
					return
				}
				goto loop
			}
			return

		case <-worker.stop:
			wlog.Debug("Worker: Job has been canceled via Worker Stop", wlog.String("worker", worker.name), wlog.String("job_id", job.Id))
			worker.setJobCanceled(job)
			return
		}
	}
}

func (worker *Worker) setJobSuccess(job *model.Job) {
	if err := worker.app.Jobs.SetJobSuccess(job); err != nil {
		wlog.Error("Worker: Failed to set success for job", wlog.String("worker", worker.name), wlog.String("job_id", job.Id), wlog.String("error", err.Error()))
		worker.setJobError(job, err)
	}
}

func (worker *Worker) setJobError(job *model.Job, appError *model.AppError) {
	if err := worker.app.Jobs.SetJobError(job, appError); err != nil {
		wlog.Error("Worker: Failed to set job error", wlog.String("worker", worker.name), wlog.String("job_id", job.Id), wlog.String("error", err.Error()))
	}
}

func (worker *Worker) setJobCanceled(job *model.Job) {
	if err := worker.app.Jobs.SetJobCanceled(job); err != nil {
		wlog.Error("Worker: Failed to mark job as canceled", wlog.String("worker", worker.name), wlog.String("job_id", job.Id), wlog.String("error", err.Error()))
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
			//TODO
			panic(err)
		}

		err = backend.Remove(f)
		if err != nil {
			wlog.Warn(fmt.Sprintf("Failed to remove file %s from store, err: %s", f.GetStoreName(), err.Error()))
			///TODO check error NOT_EXISTS and set true, and set not return error;
			//<-worker.app.Store.File().SetNoExistsById(f.Id, true)
			return err, false
		}

		resultDel = <-worker.app.Store.File().DeleteById(f.Id)
		if resultDel.Err != nil {
			wlog.Warn(fmt.Sprintf("Failed to remove store file %s from DB, err: %s", f.GetStoreName(), resultDel.Err.Error()))
			return err, false
		}
		wlog.Debug(fmt.Sprintf("File %s[%d] removed", f.GetStoreName(), f.Id))
		removed++
	}
	return nil, len(data) == 0
}
