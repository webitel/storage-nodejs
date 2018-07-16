package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
)

type Schedulers struct {
	stop       chan bool
	stopped    chan bool
	listenerId string
	startOnce  sync.Once
	jobs       *JobServer

	schedulers   []model.Scheduler
	nextRunTimes []*time.Time
}

func (srv *JobServer) InitSchedulers() *Schedulers {
	mlog.Debug("Initialising schedulers.")

	schedulers := &Schedulers{
		stop:    make(chan bool),
		stopped: make(chan bool),
		jobs:    srv,
	}

	if srv.UploadRecordingsJob != nil {
		schedulers.schedulers = append(schedulers.schedulers, srv.UploadRecordingsJob.MakeScheduler())
	}

	schedulers.nextRunTimes = make([]*time.Time, len(schedulers.schedulers))
	return schedulers
}

func (schedulers *Schedulers) Start() *Schedulers {
	go func() {
		schedulers.startOnce.Do(func() {
			mlog.Info("Starting schedulers.")

			defer func() {
				mlog.Info("Schedulers stopped.")
				close(schedulers.stopped)
			}()

			now := time.Now()
			for idx, scheduler := range schedulers.schedulers {
				if !scheduler.Enabled(schedulers.jobs.Config()) {
					schedulers.nextRunTimes[idx] = nil
				} else {
					schedulers.setNextRunTime(schedulers.jobs.Config(), idx, now, false)
				}
			}

			for {
				select {
				case <-schedulers.stop:
					mlog.Debug("Schedulers received stop signal.")
					return
				case now = <-time.After(1 * time.Minute):
					cfg := schedulers.jobs.Config()

					for idx, nextTime := range schedulers.nextRunTimes {
						if nextTime == nil {
							continue
						}

						if time.Now().After(*nextTime) {
							scheduler := schedulers.schedulers[idx]
							if scheduler != nil {
								if scheduler.Enabled(cfg) {
									if _, err := schedulers.scheduleJob(cfg, scheduler); err != nil {
										mlog.Warn(fmt.Sprintf("Failed to schedule job with scheduler: %v", scheduler.Name()))
										mlog.Error(fmt.Sprint(err))
									} else {
										schedulers.setNextRunTime(cfg, idx, now, true)
									}
								}
							}
						}
					}
				}
			}
		})
	}()

	return schedulers
}

func (schedulers *Schedulers) Stop() *Schedulers {
	mlog.Info("Stopping schedulers.")
	close(schedulers.stop)
	<-schedulers.stopped
	return schedulers
}

func (schedulers *Schedulers) setNextRunTime(cfg *model.Config, idx int, now time.Time, pendingJobs bool) {
	scheduler := schedulers.schedulers[idx]
	schedulers.nextRunTimes[idx] = scheduler.NextScheduleTime(cfg, now, pendingJobs, nil)
	mlog.Debug(fmt.Sprintf("Next run time for scheduler %v: %v", scheduler.Name(), schedulers.nextRunTimes[idx]))
}

func (schedulers *Schedulers) scheduleJob(cfg *model.Config, scheduler model.Scheduler) (model.Job, *model.AppError) {
	return scheduler.ScheduleJob(cfg, false, nil)
}
