package apis

import (
	"fmt"
	"github.com/webitel/storage/model"
	"net/http"
)

func (api *API) initCdr() {
	api.PublicRoutes.Cdr.Handle("/text", api.ApiHandler(searchCdr)).Methods("POST")
	api.PublicRoutes.Cdr.Handle("/text/scroll", api.ApiHandler(scrollCdr)).Methods("POST")

	api.PublicRoutes.Cdr.Handle("/{id}", api.ApiSessionRequired(getLegCdr)).Methods("GET")
	api.PublicRoutes.Cdr.Handle("/{id}/b", api.ApiSessionRequired(getCdrCall)).Methods("GET")
}

func scrollCdr(c *Context, w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	scrollReq := model.SearchEngineScrollFromJson(r.Body)

	var res *model.SearchEngineResponse
	res, c.Err = c.App.ScrollData(scrollReq)
	if c.Err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.ToJson()))
}

func searchCdr(c *Context, w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	searchReq := model.SearchEngineRequestFromJson(r.Body)
	searchReq.Index = "cdr*"
	searchReq.Type = "cdr"

	searchReq.Filter.AddFilter(map[string]interface{}{
		"term": map[string]interface{}{
			"presence_id": "10.10.10.144",
		},
	})

	tst, _ := searchReq.Filter.Source()
	fmt.Println(tst)

	if searchReq.Sort.IsEmpty() {
		searchReq.Sort.Data = map[string]interface{}{
			"created_time": map[string]interface{}{
				"order":         "desc",
				"unmapped_type": "boolean",
			},
		}
	}

	var res *model.SearchEngineResponse
	res, c.Err = c.App.SearchData(searchReq)
	if c.Err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res.ToJson()))
}

func getLegCdr(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireId()
	if c.Err != nil {
		return
	}

	var leg *model.CdrData

	if r.URL.Query().Get("leg") == "b" {
		leg, c.Err = c.App.GetCdrLegBDataByUuid(c.Params.Id)
	} else {
		leg, c.Err = c.App.GetCdrLegADataByUuid(c.Params.Id)
	}

	if c.Err != nil {
		return
	}
	w.Write([]byte(leg.ToJSON()))
}

func getCdrCall(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireId()
	if c.Err != nil {
		return
	}

	var call *model.CdrCall

	call, c.Err = c.App.GetCdrByUuidCall(c.Params.Id)

	if c.Err != nil {
		return
	}
	w.Write([]byte(call.ToJSON()))
}
