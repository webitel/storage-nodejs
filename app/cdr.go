package app

import "github.com/webitel/storage/model"

func (app *App) GetCdrLegADataByUuid(uuid string) (*model.CdrData, *model.AppError) {
	if result := <-app.Store.Cdr().GetLegADataByUuid(uuid); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.CdrData), nil
	}
}

func (app *App) GetCdrLegBDataByUuid(uuid string) (*model.CdrData, *model.AppError) {
	if result := <-app.Store.Cdr().GetLegBDataByUuid(uuid); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.CdrData), nil
	}
}

func (app *App) GetCdrByUuidCall(uuid string) (*model.CdrCall, *model.AppError) {
	if result := <-app.Store.Cdr().GetByUuidCall(uuid); result.Err != nil {
		return nil, result.Err
	} else {
		return result.Data.(*model.CdrCall), nil
	}
}
