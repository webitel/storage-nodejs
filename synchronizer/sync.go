package synchronizer

import (
	"fmt"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/interfaces"
	"github.com/webitel/storage/pool"
	"github.com/webitel/wlog"
	"sync"
	"time"
)

type synchronizer struct {
	App               *app.App
	betweenAttemptSec int64
	limit             int
	schedule          chan struct{}
	pollingInterval   time.Duration
	stopSignal        chan struct{}
	pool              interfaces.PoolInterface
	mx                sync.RWMutex
	stopped           bool
}

func init() {
	app.RegisterSynchronizer(func(a *app.App) interfaces.SynchronizerFilesInterface {
		wlog.Debug("Initialize synchronizer")
		return &synchronizer{
			App:               a,
			limit:             100,
			betweenAttemptSec: 60,
			schedule:          make(chan struct{}, 1),
			stopSignal:        make(chan struct{}),
			pollingInterval:   time.Second * 2,
			pool:              pool.NewPool(5, 10), //FIXME added config
		}
	})
}

func (s *synchronizer) Start() {
	wlog.Debug("Run synchronizer")
	go s.run()
}

func (s *synchronizer) run() {
	var i int
	for {
		select {
		case <-s.schedule:
		case <-time.After(s.pollingInterval):
		start:

			jobs, err := s.App.Store.SyncFile().FetchRemoveJobs(s.limit)
			if err != nil {
				wlog.Error(err.Error())
				continue
			}
			count := len(jobs)

			if count > 0 {
				wlog.Debug(fmt.Sprintf("fetch %d remove file jobs", count))
				for i = 0; i < count; i++ {
					s.pool.Exec(&job{
						app:  s.App,
						file: *jobs[i],
					})
				}

				if count == s.limit && !s.isStopped() {
					goto start
				}
			}

		case <-s.stopSignal:
			wlog.Debug("Synchronizer received stop signal.")
			return
		}
	}
}

func (s *synchronizer) isStopped() bool {
	s.mx.RLock()
	defer s.mx.RUnlock()
	return s.stopped
}

func (s *synchronizer) Stop() {
	s.mx.Lock()
	s.stopped = true
	s.mx.Unlock()

	s.stopSignal <- struct{}{}
	s.pool.Close()
	s.pool.Wait()
	wlog.Debug("Synchronizer stopped.")
}
