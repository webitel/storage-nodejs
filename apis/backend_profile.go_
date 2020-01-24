package apis

import (
	"github.com/webitel/storage/model"
	"net/http"
	"strconv"
)

func (api *API) InitBackendProfile() {
	api.PublicRoutes.BackendProfile.Handle("", api.ApiSessionRequired(listProfiles)).Methods("GET")
	api.PublicRoutes.BackendProfile.Handle("", api.ApiSessionRequired(createProfile)).Methods("POST")
	api.PublicRoutes.BackendProfile.Handle("/{id:[0-9]+}", api.ApiSessionRequired(getProfile)).Methods("GET")
	api.PublicRoutes.BackendProfile.Handle("/{id:[0-9]+}", api.ApiSessionRequired(deleteProfile)).Methods("DELETE")
	api.PublicRoutes.BackendProfile.Handle("/{id:[0-9]+}", api.ApiSessionRequired(updateProfile)).Methods("PUT")
}

func createProfile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	if c.Err != nil {
		return
	}

	profile := model.FileBackendProfileFromJson(r.Body)
	defer r.Body.Close()

	if profile == nil {
		c.SetInvalidParam("profile")
		return
	}
	profile.Domain = c.Params.Domain

	if profile, err := c.App.SaveFileBackendProfile(profile); err != nil {
		c.Err = err
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(profile.ToJson()))
	}
}

func getProfile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	c.RequireId()

	if c.Err != nil {
		return
	}

	var id int
	id, _ = strconv.Atoi(c.Params.Id)

	profile, err := c.App.GetFileBackendProfile(id, c.Params.Domain)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(profile.ToJson()))
}

func listProfiles(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()

	if c.Err != nil {
		return
	}

	listProfiles, err := c.App.ListFileBackendProfiles(c.Params.Domain, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	} else {
		w.Write([]byte(model.FileBackendProfileListToJson(listProfiles)))
	}
}

func deleteProfile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	c.RequireId()

	if c.Err != nil {
		return
	}

	var id int
	id, _ = strconv.Atoi(c.Params.Id)
	c.Err = c.App.RemoveFileBackendProfiles(id, c.Params.Domain)
	if c.Err != nil {
		return
	}

	ReturnStatusOK(w)
}

func updateProfile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	c.RequireId()
	if c.Err != nil {
		return
	}

	path := model.FileBackendProfilePathFromJson(r.Body)
	if path == nil {
		c.SetInvalidParam("profile")
		return
	}

	var id int
	var profile *model.FileBackendProfile
	id, _ = strconv.Atoi(c.Params.Id)

	profile, c.Err = c.App.GetFileBackendProfile(id, c.Params.Domain)
	if c.Err != nil {
		return
	}

	profile, c.Err = c.App.PathFileBackendProfile(profile, path)
	if c.Err != nil {
		return
	}
	w.Write([]byte(profile.ToJson()))
}
