package jobs

import (
	"sync"

	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
)

type Workers struct {
	startOnce sync.Once
	Watcher   *Watcher

	listenerId string
}

func (srv *JobServer) NewWorker(interval int) *Workers {
	worker := &Workers{}
	worker.Watcher = srv.MakeWatcher(worker, interval)

	go func() {
		worker.Start()
	}()

	return worker
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
	mlog.Info("Stopped workers")

	workers.Watcher.Stop()

	return workers
}
