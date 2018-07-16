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
}

func (srv *JobServer) InitWorkers() *Workers {
	workers := &Workers{
		ConfigService: srv.ConfigService,
	}
	workers.Watcher = srv.MakeWatcher(workers, DEFAULT_WATCHER_POLLING_INTERVAL)

	return workers
}

func (workers *Workers) Start() *Workers {
	mlog.Info("Starting workers")

	workers.startOnce.Do(func() {
		go workers.Watcher.Start()
	})

	return workers
}

func (workers *Workers) handleConfigChange(oldConfig *model.Config, newConfig *model.Config) {
	mlog.Debug("Workers received config change.")
}

func (workers *Workers) Stop() *Workers {

	workers.Watcher.Stop()

	mlog.Info("Stopped workers")

	return workers
}
