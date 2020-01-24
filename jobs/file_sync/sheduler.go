package file_sync

import (
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
)

type Scheduler struct {
	App *app.App
}

func (s *SyncFilesJobInterfaceImpl) MakeScheduler() model.Scheduler {
	return &Scheduler{s.App}
}

func (scheduler *Scheduler) Name() string {
	return "Sync files"
}

func (scheduler *Scheduler) JobType() string {
	return model.JOB_TYPE_SYNC_FILES
}

func (scheduler *Scheduler) Enabled(cfg *model.Config) bool {
	return true
}

func (scheduler *Scheduler) ScheduleJob(cfg *model.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model.AppError) {
	wlog.Debug("Scheduling Job", wlog.String("scheduler", scheduler.Name()))

	data := map[string]string{}

	if job, err := scheduler.App.Jobs.CreateJob(model.JOB_TYPE_SYNC_FILES, nil, 0, data); err != nil {
		return nil, err
	} else {
		return job, nil
	}
}
