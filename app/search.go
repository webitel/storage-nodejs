package app

import "github.com/webitel/storage/model"

func (app *App) SearchData(request *model.SearchEngineRequest) (*model.SearchEngineResponse, *model.AppError) {
	if res := <-app.Store.Search(request); res.Err != nil {
		return nil, res.Err
	} else {
		return res.Data.(*model.SearchEngineResponse), nil
	}
}

func (app *App) ScrollData(request *model.SearchEngineScroll) (*model.SearchEngineResponse, *model.AppError) {
	if res := <-app.Store.Scroll(request); res.Err != nil {
		return nil, res.Err
	} else {
		return res.Data.(*model.SearchEngineResponse), nil
	}
}
