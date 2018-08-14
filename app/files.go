package app

import (
	"github.com/webitel/storage/model"
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
