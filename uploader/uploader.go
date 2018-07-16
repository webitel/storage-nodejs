package uploader

import (
	"fmt"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/einterfaces"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/pool"
	"github.com/webitel/storage/store"
	"sync"
	"time"
)

type UploaderInterfaceImpl struct {
	App               *app.App
	betweenAttemptSec int64
	limit             int
	schedule          chan struct{}
	pollingInterval   time.Duration
	stopSignal        chan struct{}
	pool              einterfaces.PoolInterface
	mx                sync.RWMutex
	stopped           bool
}

func init() {
	app.RegisterUploader(func(a *app.App) einterfaces.UploadRecordingsFilesInterface {
		mlog.Debug("Initialize uploader")
		return &UploaderInterfaceImpl{
			App:               a,
			limit:             10,
			betweenAttemptSec: 60,
			schedule:          make(chan struct{}, 1),
			stopSignal:        make(chan struct{}),
			pollingInterval:   time.Second * 2,
			pool:              pool.NewPool(10, 1),
		}
	})
}

func (u *UploaderInterfaceImpl) Start() {
	mlog.Debug("Run uploader")
	go u.run()
}

func (u *UploaderInterfaceImpl) run() {
	var result store.StoreResult
	var jobs = []*model.JobUploadFileWithProfile{}
	var count int
	var i int
	for {
		select {
		case <-u.schedule:
		case <-time.After(u.pollingInterval):
		start:
			if result = <-u.App.Store.UploadJob().UpdateWithProfile(u.limit, u.App.GetInstanceId(), u.betweenAttemptSec); result.Err != nil {
				mlog.Critical(fmt.Sprint(result.Err))
				continue
			}
			jobs = result.Data.([]*model.JobUploadFileWithProfile)

			count = len(jobs)
			if count > 0 {
				mlog.Debug(fmt.Sprintf("Found uploading files %d", count))
				for i = 0; i < count; i++ {
					u.pool.Exec(&UploadTask{
						app: u.App,
						job: jobs[i],
					})
				}

				if count == u.limit && !u.isStopped() {
					goto start
				}
			}
		case <-u.stopSignal:
			mlog.Debug("Uploader received stop signal.")
			return
		}
	}
}

func (u *UploaderInterfaceImpl) isStopped() bool {
	u.mx.RLock()
	defer u.mx.RUnlock()
	return u.stopped
}

func (u *UploaderInterfaceImpl) Stop() {
	u.mx.Lock()
	u.stopped = true
	u.mx.Unlock()

	u.stopSignal <- struct{}{}
	u.pool.Close()
	u.pool.Wait()
	mlog.Debug("Uploader stopped.")
}
