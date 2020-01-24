package apis

import (
	"github.com/webitel/storage/model"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

func (api *API) InitMediaFile() {
	api.PublicRoutes.MediaFiles.Handle("", api.ApiHandler(listMediaFiles)).Methods("GET")
	api.PublicRoutes.MediaFiles.Handle("/{id}", api.ApiHandler(getMediaFile)).Methods("GET")
	api.PublicRoutes.MediaFiles.Handle("/{id}/stream", api.ApiHandler(streamMediaFile)).Methods("GET")

	// Old version
	api.PublicRoutes.MediaFiles.Handle("/{id}", api.ApiHandler(saveMediaFile)).Methods("POST")
	api.PublicRoutes.MediaFiles.Handle("/", api.ApiHandler(saveMediaFile)).Methods("POST")

	api.PublicRoutes.MediaFiles.Handle("/{id}", api.ApiHandler(removeMediaFile)).Methods("DELETE")
}

func getMediaFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	c.RequireId()

	if c.Err != nil {
		return
	}

	var file *model.MediaFile

	if file, c.Err = c.App.GetMediaFileByName(c.Params.Id, c.Params.Domain); c.Err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(file.ToJson()))
}

func streamMediaFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	c.RequireId()

	if c.Err != nil {
		return
	}

	var file *model.MediaFile

	if file, c.Err = c.App.GetMediaFileByName(c.Params.Id, c.Params.Domain); c.Err != nil {
		return
	}

	ranges, err := parseRange(r.Header.Get("Range"), int64(file.Size))
	if err != nil {
		c.Err = err
		return
	}

	var offset int64 = 0
	sendSize := file.Size
	code := http.StatusOK

	switch {
	case len(ranges) == 1:
		code = http.StatusPartialContent
		offset = ranges[0].Start
		sendSize = ranges[0].Length
		w.Header().Set("Content-Range", ranges[0].ContentRange(file.Size))
	default:
		//TODO
	}

	reader, err := c.App.MediaFileStore.Reader(file, offset)
	if err != nil {
		c.Err = err
		return
	}
	defer reader.Close()

	w.Header().Set("Accept-Ranges", "bytes")
	//w.Header().Set("Cache-Control", "no-cache, must-revalidate, max-age=0")
	w.Header().Set("Content-Type", file.MimeType)

	if w.Header().Get("Content-Encoding") == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
	}

	w.WriteHeader(code)
	io.CopyN(w, reader, sendSize)
}

func saveMediaFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()

	if c.Err != nil {
		return
	}

	defer r.Body.Close()

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(mediaType, "multipart/form-data") {
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
	} else {
		file := &model.MediaFile{}
		file.Properties = model.StringInterface{}
		file.Name = r.URL.Query().Get("name")
		file.Domain = c.Params.Domain
		file.MimeType = r.Header.Get("Content-Type")

		if _, err := c.App.SaveMediaFile(r.Body, file); err != nil {
			c.Err = err
			return
		}
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

func removeMediaFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	c.RequireId()

	if c.Err != nil {
		return
	}

	_, c.Err = c.App.RemoveMediaFileByName(c.Params.Id, c.Params.Domain)

	if c.Err != nil {
		return
	}

	ReturnStatusOK(w)
}
