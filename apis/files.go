package apis

import (
	"fmt"
	"github.com/webitel/storage/model"
	"io"
	"net/http"
	"strconv"
)

func (api *API) InitFiles() {
	api.PublicRoutes.Files.Handle("", api.ApiSessionRequired(listFiles)).Methods("GET")
	api.PublicRoutes.Files.Handle("/{id}", api.ApiSessionRequired(getFile)).Methods("GET")
}

func listFiles(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()

	if c.Err != nil {
		return
	}

	listFiles, err := c.App.ListFiles(c.Params.Domain, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	} else {
		w.Write([]byte(model.FileListToJson(listFiles)))
	}
}

func getFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireDomain()
	c.RequireId()

	if c.Err != nil {
		return
	}

	files, err := c.App.GetFile(c.Params.Domain, c.Params.Id)
	if err != nil {
		c.Err = err
		return
	}

	if len(files) > 0 {
		fmt.Println(files)
	}

	file := files[0]

	store, err := c.App.GetFileBackendStore(file.ProfileId, file.ProfileUpdatedAt)
	if err != nil {
		c.Err = err
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

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Cache-Control", "no-cache, must-revalidate, max-age=0")
	w.Header().Set("Content-Type", file.MimeType)

	if w.Header().Get("Content-Encoding") == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(sendSize, 10))
	}

	reader, err := store.Reader(file.File, offset)
	if err != nil {
		c.Err = err
		return
	}
	defer reader.Close()

	w.WriteHeader(code)
	io.CopyN(w, reader, sendSize)
}
