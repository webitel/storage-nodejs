package apis

import (
	"fmt"
	"github.com/webitel/storage/model"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
)

func (api *API) InitMediaFile() {
	api.PublicRoutes.MediaFiles.Handle("", api.ApiHandler(listMediaFiles)).Methods("GET")
	// Old version
	api.PublicRoutes.MediaFiles.Handle("/{id}", api.ApiHandler(saveMediaFile)).Methods("POST")
	api.PublicRoutes.MediaFiles.Handle("/", api.ApiHandler(saveMediaFile)).Methods("POST")
}

func saveMediaFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()

	if c.Err != nil {
		return
	}

	defer r.Body.Close()

	//TODO check mime multipart
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}
	fmt.Println(mediaType)
	writer := multipart.NewReader(r.Body, params["boundary"])

	for {
		part, err := writer.NextPart()
		if err == io.EOF {
			return
		}

		if err != nil {
			panic(err)
		}

		file := &model.MediaFile{}
		file.Properties = model.StringInterface{}
		file.Name = part.FileName()
		file.Domain = c.Params.Domain
		file.MimeType = part.Header.Get("Content-Type")
		if _, err := c.App.SaveMediaFile(part, file); err != nil {
			c.Err = err
			return
		}
		part.Close()
	}

}

func listMediaFiles(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()

	if c.Err != nil {
		return
	}

	var files []*model.MediaFile
	var count int64

	if files, c.Err = c.App.ListMediaFiles(c.Params.Domain, c.Params.Page, c.Params.PerPage); c.Err != nil {
		return
	}
	if count, c.Err = c.App.CountMediaFilesByDomain(c.Params.Domain); c.Err != nil {
		return
	}

	response := &ListResponse{
		Total: count,
		Items: files,
	}

	w.Write([]byte(response.ToJson()))

}
