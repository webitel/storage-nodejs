package jobs

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
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
	mlog.Debug("Watcher Started")

	// Delay for some random number of milliseconds before starting to ensure that multiple
	// instances of the jobserver  don't poll at a time too close to each other.
	rand.Seed(time.Now().UTC().UnixNano())
	<-time.After(time.Duration(rand.Intn(watcher.pollingInterval)) * time.Millisecond)

	defer func() {
		mlog.Debug("Watcher Finished")
		watcher.stopped <- true
	}()

	for {
		select {
		case <-watcher.stop:
			mlog.Debug("Watcher: Received stop signal")
			return
		case <-time.After(time.Duration(watcher.pollingInterval) * time.Millisecond):
			watcher.PollAndNotify()
		}
	}
}

func (watcher *Watcher) Stop() {
	mlog.Debug("Watcher Stopping")
	watcher.stop <- true
	<-watcher.stopped
}

func (watcher *Watcher) PollAndNotify() {
	if result := <-watcher.srv.Store.Job().GetAllByStatus(model.JOB_STATUS_PENDING); result.Err != nil {
		mlog.Error(fmt.Sprintf("Error occurred getting all pending statuses: %v", result.Err.Error()))
	} else {
		jobs := result.Data.([]*model.Job)

		for _, job := range jobs {
			if job.Type == model.JOB_TYPE_SYNC_FILES {
				if watcher.workers.SyncFile != nil {
					select {
					case watcher.workers.SyncFile.JobChannel() <- *job:
					default:
					}
				}
			}
		}
	}
}
