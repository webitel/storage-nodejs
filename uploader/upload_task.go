package uploader

import (
	"fmt"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
)

type UploadTask struct {
	app *app.App
	job *model.JobUploadFileWithProfile
}

func (u *UploadTask) Name() string {
	return u.job.Uuid
}

//TODO added max count attempts ?

func (u *UploadTask) Execute() {
	store, err := u.app.GetFileBackendStore(u.job.ProfileId, u.job.ProfileUpdatedAt)

	if err != nil {
		mlog.Critical(err.Error())
		u.app.Store.UploadJob().SetStateError(int(u.job.Id), err.Error())
		return
	}

	mlog.Debug(fmt.Sprintf("Execute upload task %s to store %s", u.Name(), store.Name()))

	r, err := u.app.FileBackendLocal.Reader(u.job, 0)
	if err != nil {
		mlog.Critical(err.Error())
		u.app.Store.UploadJob().SetStateError(int(u.job.Id), err.Error())
		return
	}
	defer r.Close()

	storeName := u.job.GetStoreName()
	directory := store.GetStoreDirectory(u.job.Domain)

	if _, err = store.WriteFile(r, directory, storeName); err != nil {
		mlog.Critical(err.Error())
		u.app.Store.UploadJob().SetStateError(int(u.job.Id), err.Error())
		return
	}

	mlog.Debug(fmt.Sprintf("Store to %s/%s %d bytes", directory, storeName, u.job.Size))

	result := <-u.app.Store.File().MoveFromJob(int(u.job.Id), u.job.ProfileId, model.StringInterface{"directory": directory})
	if result.Err != nil {
		mlog.Critical(result.Err.Error())
		u.app.Store.UploadJob().SetStateError(int(u.job.Id), result.Err.Error())
		return
	}

	err = u.app.FileBackendLocal.RemoveFile(u.app.FileBackendLocal.GetStoreDirectory(""), u.job.GetStoreName())
	if err != nil {
		mlog.Critical(err.Error())
	}

	mlog.Debug(fmt.Sprintf("End execute upload task %s", u.Name()))
}
