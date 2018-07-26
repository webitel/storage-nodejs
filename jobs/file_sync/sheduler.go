package file_sync

import (
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
	"time"
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

func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {
	t, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	t = t.AddDate(0, 0, 1)
	return &t
}

func (scheduler *Scheduler) ScheduleJob(cfg *model.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model.AppError) {
	mlog.Debug("Scheduling Job", mlog.String("scheduler", scheduler.Name()))

	data := map[string]string{}

	if job, err := scheduler.App.Jobs.CreateJob(model.JOB_TYPE_SYNC_FILES, data); err != nil {
		return nil, err
	} else {
		return job, nil
	}
}
