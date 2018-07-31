package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
)

type Schedulers struct {
	stop          chan bool
	stopped       chan bool
	configChanged chan *model.Config
	listenerId    string
	startOnce     sync.Once
	jobs          *JobServer

	schedulers []model.Scheduler
}

func (srv *JobServer) InitSchedulers() *Schedulers {
	mlog.Debug("Initialising schedulers.")

	schedulers := &Schedulers{
		stop:          make(chan bool),
		stopped:       make(chan bool),
		configChanged: make(chan *model.Config),
		jobs:          srv,
	}

	if syncFilesJobInterface := srv.SyncFilesJob; syncFilesJobInterface != nil {
		schedulers.schedulers = append(schedulers.schedulers, syncFilesJobInterface.MakeScheduler())
	}

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

			schedulers.scheduleJobs(&now)

			for {
				select {
				case <-schedulers.stop:
					mlog.Debug("Schedulers received stop signal.")
					return

				case now = <-time.After(time.Minute):
					schedulers.scheduleJobs(&now)
				}
			}
		})
	}()

	return schedulers
}

func (schedulers *Schedulers) scheduleJobs(now *time.Time) {
	var nextTime int64
	var appErr *model.AppError

	res := <-schedulers.jobs.Store.Schedule().GetAllWithNoJobs(1000, 0)

	if res.Err != nil {
		mlog.Critical(res.Err.Error())
		return
	}

	for _, item := range res.Data.([]*model.Schedule) {
		data := make(map[string]string)
		nextTime = item.NextTime(*now)

		_, appErr = schedulers.jobs.CreateJob(item.Type, &item.Id, nextTime, data)
		if appErr != nil {
			mlog.Warn(fmt.Sprintf("Failed to schedule job with scheduler: %s", item.Name))
			mlog.Error(fmt.Sprint(appErr))
		} else {
			mlog.Debug(fmt.Sprintf("Next run time for scheduler %v: %v", item.Name, time.Unix(nextTime, 0).String()))
		}
	}
}

func (schedulers *Schedulers) Stop() *Schedulers {
	mlog.Info("Stopping schedulers.")
	close(schedulers.stop)
	<-schedulers.stopped
	return schedulers
}
