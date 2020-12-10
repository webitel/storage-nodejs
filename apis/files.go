package apis

import (
	"github.com/webitel/storage/model"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

func (api *API) InitFile() {
	api.PublicRoutes.Files.Handle("/{id}/stream", api.ApiSessionRequired(streamRecordFile)).Methods("GET")
	api.PublicRoutes.Files.Handle("/{id}/download", api.ApiSessionRequired(downloadRecordFile)).Methods("GET")
	api.PublicRoutes.Files.Handle("/{id}/upload", api.ApiSessionRequired(uploadAnyFile)).Methods("POST")
}

func uploadAnyFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireId()

	if c.Err != nil {
		return
	}

	defer r.Body.Close()

	files := make([]*model.JobUploadFile, 0)

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(mediaType, "multipart/form-data") {
		writer := multipart.NewReader(r.Body, params["boundary"])

		for {
			part, err := writer.NextPart()
			if err == io.EOF {
				break
			}

			if err != nil {
				//panic(err)
				break //fixme
			}

			file := &model.JobUploadFile{}
			file.Name = model.NewId() + "_" + part.FileName()
			file.MimeType = part.Header.Get("Content-Type")
			file.DomainId = c.Session.DomainId
			file.Uuid = c.Params.Id

			// TODO PERMISSION
			if err := c.App.AddUploadJobFile(r.Body, file); err != nil {
				c.Err = err
				return
			}

			files = append(files, file)
			part.Close()
		}
	} else {
		file := &model.JobUploadFile{}
		file.Name = model.NewId() + "_" + r.URL.Query().Get("name")
		file.MimeType = r.Header.Get("Content-Type")
		file.DomainId = c.Session.DomainId
		file.Uuid = c.Params.Id

		// TODO PERMISSION
		if err := c.App.AddUploadJobFile(r.Body, file); err != nil {
			c.Err = err
			return
		}
	}

	if c.Err != nil {
		return
	}

	response := &ListResponse{
		Items: files,
	}

	w.Write([]byte(response.ToJson()))

	//c.App.GenerateSignature() // todo app generate public download
}
