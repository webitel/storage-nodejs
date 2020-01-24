package app

import (
	"fmt"
	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
	"io"
)

func (app *App) AddUploadJobFile(src io.ReadCloser, file *model.JobUploadFile) *model.AppError {

	size, err := app.FileCache.Write(src, file)
	if err != nil {
		return err
	}

	file.Size = size
	file.Instance = app.GetInstanceId()

	res := <-app.Store.UploadJob().Save(file)
	if res.Err != nil {
		if err = app.FileCache.Remove(file); err != nil {
			wlog.Error(fmt.Sprintf("Failed to remove cache file %v", err))
		}
	}

	return res.Err
}

func (app *App) RemoveUploadJob(id int) *model.AppError {
	return nil
}