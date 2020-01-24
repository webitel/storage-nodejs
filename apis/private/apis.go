package private

import (
	"github.com/gorilla/mux"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/web"
	"net/http"
)

type RoutesInternal struct {
	Root    *mux.Router // ''
	ApiRoot *mux.Router // 'sys'
	Files   *mux.Router // '/records'
	Test    *mux.Router // '/test'
}

type API struct {
	App    *app.App
	Routes *RoutesInternal
}

func Init(a *app.App, root *mux.Router) *API {
	api := &API{
		App:    a,
		Routes: &RoutesInternal{},
	}

	api.Routes.Root = root
	api.Routes.ApiRoot = root.PathPrefix(model.API_INTERNAL_URL_SUFFIX_V1).Subrouter()
	api.Routes.Files = api.Routes.ApiRoot.PathPrefix("/formLoadFile").Subrouter()
	api.Routes.Test = api.Routes.ApiRoot.PathPrefix("/test").Subrouter()

	api.InitTest()

	api.InitFile()
	return api
}

func (api *API) InitTest() {
	//api.Routes.Test.Handle("", api.ApiHandler(getTest)).Methods("GET")
}

func (api *API) Handle404(w http.ResponseWriter, r *http.Request) {
	web.Handle404(api.App, w, r)
}

//
//func getTest(c *Context, w http.ResponseWriter, r *http.Request)  {
//	res := <-c.App.Store.FileBackendProfile().Get(1, "10.10.10.144")
//	if res.Err != nil {
//		c.Err = res.Err
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	w.Write([]byte(res.Data.(*model.FileBackendProfile).ToJson()))
//}
var ReturnStatusOK = web.ReturnStatusOK
