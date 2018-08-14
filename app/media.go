package app

import (
	"github.com/webitel/storage/model"
	"io"
	"net/http"
)

func (app *App) SaveMediaFile(src io.Reader, mediaFile *model.MediaFile) (*model.MediaFile, *model.AppError) {

	size, err := app.MediaFileStore.Write(src, mediaFile)
	if err != nil {
		return nil, err
	}
	mediaFile.Size = size
	mediaFile.Instance = app.GetInstanceId()

	if err = mediaFile.IsValid(); err != nil {
		return nil, err
	}

	if !model.StringInSlice(mediaFile.MimeType, app.Config().MediaFileStoreSettings.AllowMime) {
		return nil, model.NewAppError("app.SaveMediaFile", "model.media_file.mime_type.app_error", nil,
			"Not allowed mime type", http.StatusBadRequest)
	}

	if app.Config().MediaFileStoreSettings.MaxSizeByte != nil && *app.Config().MediaFileStoreSettings.MaxSizeByte < int(size) {
		return nil, model.NewAppError("app.SaveMediaFile", "model.media_file.size.app_error", nil,
			"", http.StatusBadRequest)
	}

	if result := <-app.Store.MediaFile().Save(mediaFile); result.Err != nil {
		if result.Err.Id != "store.sql_media_file.save.saving.duplicate" {
			app.MediaFileStore.Remove(mediaFile)
		}
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

func (app *App) GetMediaFile(id int64, domain string) (*model.MediaFile, *model.AppError) {
	if result := <-app.Store.MediaFile().Get(id, domain); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.MediaFile), nil
	}
}

func (app *App) GetMediaFileByName(name, domain string) (*model.MediaFile, *model.AppError) {
	if result := <-app.Store.MediaFile().GetByName(name, domain); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.MediaFile), nil
	}
}

func (app *App) RemoveMediaFileByName(name, domain string) (file *model.MediaFile, err *model.AppError) {

	file, err = app.GetMediaFileByName(name, domain)
	if err != nil {
		return
	}

	err = app.MediaFileStore.Remove(file)
	if err != nil {
		return
	}

	result := <-app.Store.MediaFile().DeleteById(file.Id)
	return nil, result.Err
}
