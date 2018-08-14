package app

import (
	"github.com/webitel/storage/model"
	"io"
)

func (app *App) ListFiles(domain string, page, perPage int) ([]*model.File, *model.AppError) {
	if result := <-app.Store.File().GetAllPageByDomain(domain, page*perPage, perPage); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.File), nil
	}
}

func (app *App) GetFile(domain, uuid string) ([]*model.FileWithProfile, *model.AppError) {
	if result := <-app.Store.File().Get(domain, uuid); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.FileWithProfile), nil
	}
}

func (app *App) SaveMediaFile(src io.Reader, mediaFile *model.MediaFile) (*model.MediaFile, *model.AppError) {
	directory := app.MediaFileStore.GetStoreDirectory(mediaFile.Domain)
	size, err := app.MediaFileStore.WriteFile(src, directory, mediaFile.GetStoreName())
	if err != nil {
		return nil, err
	}
	mediaFile.Size = int64(size)
	mediaFile.Instance = app.GetInstanceId()

	if result := <-app.Store.MediaFile().Save(mediaFile); result.Err != nil {
		app.MediaFileStore.Remove(mediaFile)
		return nil, result.Err
	} else {
		return result.Data.(*model.MediaFile), nil
	}
}

func (app *App) ListMediaFiles(domain string, page, perPage int) ([]*model.MediaFile, *model.AppError) {
	if result := <-app.Store.MediaFile().GetAllByDomain(domain, page*perPage, perPage); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.([]*model.MediaFile), nil
	}
}

func (app *App) CountMediaFilesByDomain(domain string) (int64, *model.AppError) {
	if result := <-app.Store.MediaFile().GetCountByDomain(domain); result.Err != nil {
		return 0, result.Err
	} else {
		return result.Data.(int64), nil
	}
}
