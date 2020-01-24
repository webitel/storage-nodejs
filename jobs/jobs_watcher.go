package jobs

import (
	"fmt"
	"time"

	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
)

// Default polling interval for jobs termination.
// (Defining as `var` rather than `const` allows tests to lower the interval.)
var DEFAULT_WATCHER_POLLING_INTERVAL = 15000

type Watcher struct {
	srv     *JobServer
	workers *Workers

	stop            chan bool
	stopped         chan bool
	pollingInterval int
}

func (srv *JobServer) MakeWatcher(workers *Workers, pollingInterval int) *Watcher {
	return &Watcher{
		stop:            make(chan bool, 1),
		stopped:         make(chan bool, 1),
		pollingInterval: pollingInterval,
		workers:         workers,
		srv:             srv,
	}
}

func (watcher *Watcher) Start() {
	wlog.Debug("Watcher Started")

	// Delay for some random number of milliseconds before starting to ensure that multiple
	// instances of the jobserver  don't poll at a time too close to each other.
	//rand.Seed(time.Now().UTC().UnixNano())
	//<-time.After(time.Duration(rand.Intn(watcher.pollingInterval)) * time.Millisecond)
	watcher.PollAndNotify()
	defer func() {
		wlog.Debug("Watcher Finished")
		watcher.stopped <- true
	}()

	for {
		select {
		case <-watcher.stop:
			wlog.Debug("Watcher: Received stop signal")
			return
		case <-time.After(time.Duration(watcher.pollingInterval) * time.Millisecond):
			watcher.PollAndNotify()
		}
	}
}

func (watcher *Watcher) Stop() {
	wlog.Debug("Watcher Stopping")
	watcher.stop <- true
	<-watcher.stopped
}

func (watcher *Watcher) PollAndNotify() {
	if result := <-watcher.srv.Store.Job().GetAllByStatusAndLessScheduleTime(model.JOB_STATUS_PENDING, model.GetMillis()); result.Err != nil {
		wlog.Error(fmt.Sprintf("Error occurred getting all pending statuses: %v", result.Err.Error()))
	} else {
		jobs := result.Data.([]*model.Job)

		for _, job := range jobs {
			if w, ok := watcher.workers.middleware[job.Type]; ok {
				w.JobChannel() <- *job
			} else {
				wlog.Warn(fmt.Sprintf("Not found middleware %s", job.Type))
			}
		}
	}
}
