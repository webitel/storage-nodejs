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

	middleware map[string]model.Worker
}

func (srv *JobServer) InitWorkers() *Workers {
	workers := &Workers{
		ConfigService: srv.ConfigService,
		middleware:    make(map[string]model.Worker),
	}
	workers.Watcher = srv.MakeWatcher(workers, DEFAULT_WATCHER_POLLING_INTERVAL)

	for key, j := range srv.middlewareJobs {
		workers.middleware[key] = j.MakeWorker()
	}

	return workers
}

func (workers *Workers) Start() *Workers {
	mlog.Info("Starting workers")

	workers.startOnce.Do(func() {

		for _, w := range workers.middleware {
			go w.Run()
		}

		go workers.Watcher.Start()
	})

	return workers
}

func (workers *Workers) Stop() *Workers {

	workers.Watcher.Stop()

	for _, w := range workers.middleware {
		w.Stop()
	}

	mlog.Info("Stopped workers")

	return workers
}
