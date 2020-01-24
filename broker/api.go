package broker

import "github.com/webitel/storage/model"

type API interface {
	ScrollData(request *model.SearchEngineScroll) (*model.SearchEngineResponse, *model.AppError)
}
