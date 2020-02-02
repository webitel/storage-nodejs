package apis

import (
	"github.com/gorilla/mux"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/web"
	"net/http"
)

type RoutesPublic struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'api/v2'

	CallRecordingsFiles *mux.Router // '/files'
	MediaFiles          *mux.Router // '/media
}

type API struct {
	App          *app.App
	PublicRoutes *RoutesPublic
	ctrl         *controller.Controller
}

func Init(a *app.App, root *mux.Router) *API {
	api := &API{
		App:          a,
		PublicRoutes: &RoutesPublic{},
		ctrl:         controller.NewController(a),
	}
	api.PublicRoutes.Root = root
	api.PublicRoutes.ApiRoot = root.PathPrefix(model.API_URL_SUFFIX).Subrouter()

	api.PublicRoutes.MediaFiles = api.PublicRoutes.ApiRoot.PathPrefix("/media").Subrouter()
	api.PublicRoutes.CallRecordingsFiles = api.PublicRoutes.ApiRoot.PathPrefix("/recordings").Subrouter()

	api.InitMediaFile()
	api.InitCallRecordingsFiles()

	return api
}

func (api *API) Handle404(w http.ResponseWriter, r *http.Request) {
	web.Handle404(api.App, w, r)
}

var ReturnStatusOK = web.ReturnStatusOK
