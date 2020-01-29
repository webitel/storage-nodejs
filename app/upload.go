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

	err = app.Store.UploadJob().Create(file)
	if err != nil {
		wlog.Error(fmt.Sprintf("Failed to store file %s, %v", file.Uuid, err))
		if errRem := app.FileCache.Remove(file); errRem != nil {
			wlog.Error(fmt.Sprintf("Failed to remove cache file %v", err))
		}
	}

	return err
}

func (app *App) RemoveUploadJob(id int) *model.AppError {
	return nil
}
