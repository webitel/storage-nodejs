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
	return model.JOB_STATUS_SUCCESS
}

func (scheduler *Scheduler) Enabled(cfg *model.Config) bool {
	return true
}

func (scheduler *Scheduler) NextScheduleTime(cfg *model.Config, now time.Time, pendingJobs bool, lastSuccessfulJob *model.Job) *time.Time {

	nextTime := time.Now().Add(1 * time.Second)
	return &nextTime
}

func (scheduler *Scheduler) ScheduleJob(cfg *model.Config, pendingJobs bool, lastSuccessfulJob *model.Job) (*model.Job, *model.AppError) {
	mlog.Debug("Scheduling Job", mlog.String("scheduler", scheduler.Name()))

	return &model.Job{}, nil
}
