package file_sync

import (
	"context"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/jobs"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
	"time"
)

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

	cancelCtx, cancelCancelWatcher := context.WithCancel(context.Background())
	cancelWatcherChan := make(chan interface{}, 1)
	go worker.app.Jobs.CancellationWatcher(cancelCtx, job.Id, cancelWatcherChan)

	defer cancelCancelWatcher()

	for {
		select {
		case <-cancelWatcherChan:
			mlog.Debug("Worker: Job has been canceled via CancellationWatcher", mlog.String("worker", worker.name), mlog.String("job_id", job.Id))
			worker.setJobCanceled(job)
			return

		case <-worker.stop:
			mlog.Debug("Worker: Job has been canceled via Worker Stop", mlog.String("worker", worker.name), mlog.String("job_id", job.Id))
			worker.setJobCanceled(job)
			return

		case <-time.After(5 * time.Second):
			mlog.Info("Worker: Job is complete", mlog.String("worker", worker.name), mlog.String("job_id", job.Id))
			worker.setJobSuccess(job)
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
