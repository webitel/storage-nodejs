package synchronizer

import (
	"fmt"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
)

type job struct {
	file model.SyncJob
	app  *app.App
}

func (j *job) Execute() {
	store, err := j.app.GetFileBackendStore(j.file.ProfileId, j.file.ProfileUpdatedAt)
	if err != nil {
		wlog.Error(err.Error())
		//todo set db error

		return
	}

	err = store.Remove(&model.File{
		BaseFile:  j.file.BaseFile,
		Id:        j.file.Id,
		DomainId:  j.file.DomainId,
		Uuid:      "",
		ProfileId: j.file.ProfileId,
		CreatedAt: 0,
	})

	if err != nil {
		wlog.Error(fmt.Sprintf("file %d, error: %s", j.file.FileId, err.Error()))
	}

	err = j.app.Store.SyncFile().Clean(j.file.Id)
	if err != nil {
		wlog.Error(fmt.Sprintf("file %d, error: %s", j.file.FileId, err.Error()))
	}

	wlog.Debug(fmt.Sprintf("file %d removed \"%s\" from store \"%s\"", j.file.FileId, j.file.Name, store.Name()))
}
