package app

import (
	"fmt"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
	"io"
)

func (app *App) AddUploadJobFile(src io.ReadCloser, file *model.JobUploadFile) *model.AppError {
	path := file.GetPath("./cache")
	size, err := app.FileCacheBackend.WriteFile(src, path)
	if err != nil {
		return err
	}

	file.Size = int(size)
	file.Instance = app.GetInstanceId()

	res := <-app.Store.UploadJob().Save(file)
	if res.Err != nil {
		if err = app.FileCacheBackend.RemoveFile(path); err != nil {
			mlog.Error(fmt.Sprintf("Failed to remove cache file %v", err))
		}
	}

	return res.Err
}
