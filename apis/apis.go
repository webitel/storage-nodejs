package apis

import (
	"github.com/gorilla/mux"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/web"
	"net/http"
)

type RoutesPublic struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'api/v2'

	BackendProfile *mux.Router // '/backend_profiles'
}

type API struct {
	App          *app.App
	PublicRoutes *RoutesPublic
}

func Init(a *app.App, root *mux.Router) *API {
	api := &API{
		App:          a,
		PublicRoutes: &RoutesPublic{},
	}
	api.PublicRoutes.Root = root
	api.PublicRoutes.ApiRoot = root.PathPrefix(model.API_URL_SUFFIX).Subrouter()
	api.PublicRoutes.BackendProfile = api.PublicRoutes.ApiRoot.PathPrefix("/backend_profiles").Subrouter()

	api.InitBackendProfile()
	return api
}

func (api *API) Handle404(w http.ResponseWriter, r *http.Request) {
	web.Handle404(api.App, w, r)
}

var ReturnStatusOK = web.ReturnStatusOK