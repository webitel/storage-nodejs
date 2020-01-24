package jobs

import (
	"fmt"
	"sync"
	"time"

	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
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
	wlog.Debug("Initialising schedulers.")

	schedulers := &Schedulers{
		stop:          make(chan bool),
		stopped:       make(chan bool),
		configChanged: make(chan *model.Config),
		jobs:          srv,
	}

	for _, w := range srv.middlewareJobs {
		schedulers.schedulers = append(schedulers.schedulers, w.MakeScheduler())
	}

	return schedulers
}

func (schedulers *Schedulers) Start() *Schedulers {

	go func() {
		schedulers.startOnce.Do(func() {
			wlog.Info("Starting schedulers.")

			defer func() {
				wlog.Info("Schedulers stopped.")
				close(schedulers.stopped)
			}()

			now := time.Now()

			schedulers.scheduleJobs(&now)

			for {
				select {
				case <-schedulers.stop:
					wlog.Debug("Schedulers received stop signal.")
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
		wlog.Critical(res.Err.Error())
		return
	}

	for _, item := range res.Data.([]*model.Schedule) {
		data := make(map[string]string)
		nextTime = item.NextTime(*now)

		_, appErr = schedulers.jobs.CreateJob(item.Type, &item.Id, nextTime, data)
		if appErr != nil {
			wlog.Warn(fmt.Sprintf("Failed to schedule job with scheduler: %s", item.Name))
			wlog.Error(fmt.Sprint(appErr))
		} else {
			wlog.Debug(fmt.Sprintf("Next run time for scheduler %v: %v", item.Name, time.Unix(nextTime, 0).String()))
		}
	}
}

func (schedulers *Schedulers) Stop() *Schedulers {
	wlog.Info("Stopping schedulers.")
	close(schedulers.stop)
	<-schedulers.stopped
	return schedulers
}
