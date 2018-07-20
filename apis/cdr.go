package apis

import (
	"github.com/webitel/storage/model"
	"net/http"
)

func (api *API) initCdr() {
	api.PublicRoutes.Cdr.Handle("/{id}", api.ApiSessionRequired(getLegCdr)).Methods("GET")
	api.PublicRoutes.Cdr.Handle("/{id}/b", api.ApiSessionRequired(getCdrCall)).Methods("GET")
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
