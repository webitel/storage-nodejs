package jobs

import (
	"sync"

	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
)

type Workers struct {
	startOnce     sync.Once
	ConfigService ConfigService
	Watcher       *Watcher

	DataRetention model.Worker
	SyncFile      model.Worker
}

func (srv *JobServer) InitWorkers() *Workers {
	workers := &Workers{
		ConfigService: srv.ConfigService,
	}
	workers.Watcher = srv.MakeWatcher(workers, DEFAULT_WATCHER_POLLING_INTERVAL)

	if syncFilesJobInterface := srv.SyncFilesJob; syncFilesJobInterface != nil {
		workers.SyncFile = syncFilesJobInterface.MakeWorker()
	}

	return workers
}

func (workers *Workers) Start() *Workers {
	mlog.Info("Starting workers")

	workers.startOnce.Do(func() {

		if workers.SyncFile != nil {
			go workers.SyncFile.Run()
		}

		go workers.Watcher.Start()
	})

	return workers
}

func (workers *Workers) Stop() *Workers {

	workers.Watcher.Stop()

	if workers.SyncFile != nil {
		workers.SyncFile.Stop()
	}

	mlog.Info("Stopped workers")

	return workers
}
