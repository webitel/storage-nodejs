package web

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	PAGE_DEFAULT     = 0
	PER_PAGE_DEFAULT = 60
	PER_PAGE_MAXIMUM = 1000
)

type Params struct {
	Domain    string
	Id        string
	Name      string
	Page      int
	PerPage   int
	Expires   int64
	Signature string
}

func ParamsFromRequest(r *http.Request) *Params {
	params := &Params{}

	props := mux.Vars(r)
	query := r.URL.Query()

	params.Domain = query.Get("domain_id")
	params.Name = query.Get("name")

	if val, ok := props["id"]; ok {
		params.Id = val
	}

	if val, err := strconv.Atoi(query.Get("page")); err != nil || val < 0 {
		params.Page = PAGE_DEFAULT
	} else {
		params.Page = val
	}

	if val, err := strconv.Atoi(query.Get("per_page")); err != nil || val < 0 {
		params.PerPage = PER_PAGE_DEFAULT
	} else if val > PER_PAGE_MAXIMUM {
		params.PerPage = PER_PAGE_MAXIMUM
	} else {
		params.PerPage = val
	}

	if val, err := strconv.Atoi(query.Get("expires")); err == nil || val > 0 {
		params.Expires = int64(val)
	}
	params.Signature = query.Get("signature")

	return params
}
