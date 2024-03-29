package apis

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/webitel/storage/model"
)

type fileResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	MimeType  string `json:"mime"`
	SharedUrl string `json:"shared"`
}

func (api *API) InitFile() {
	api.PublicRoutes.Files.Handle("/{id}/stream", api.ApiSessionRequired(streamRecordFile)).Methods("GET")
	api.PublicRoutes.Files.Handle("/{id}/download", api.ApiSessionRequired(downloadRecordFile)).Methods("GET")
	api.PublicRoutes.Files.Handle("/{id}/upload", api.ApiSessionRequired(uploadAnyFile)).Methods("POST")
	api.PublicRoutes.Files.Handle("/{id}/transcript", api.ApiSessionRequired(transcriptFile)).Methods("GET")
}

func transcriptFile(c *Context, w http.ResponseWriter, r *http.Request) {
	var id int
	var err error

	c.RequireId()

	if c.Err != nil {
		return
	}

	if id, err = strconv.Atoi(c.Params.Id); err != nil {
		c.SetInvalidUrlParam("id")
		return
	}

	//TODO
	var tr *model.FileTranscript
	if tr, c.Err = c.App.TranscriptFile(int64(id), model.TranscriptOptions{}); c.Err != nil {

		return
	}

	data, _ := json.Marshal(tr)
	w.Write(data)
}

func uploadAnyFile(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireId()

	if c.Err != nil {
		return
	}

	defer r.Body.Close()

	files := make([]*fileResponse, 0)

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
			if err := c.App.SyncUpload(part, file); err != nil {
				c.Err = err
				return
			}
			sig, _ := c.App.GeneratePreSignetResourceSignature(model.AnyFileRouteName, "download", file.Id, file.DomainId)

			files = append(files, &fileResponse{
				Id:        file.Id,
				Name:      file.Name,
				Size:      file.Size,
				MimeType:  file.MimeType,
				SharedUrl: sig,
			})
			part.Close()
		}
	} else {
		file := &model.JobUploadFile{}
		file.Name = model.NewId() + "_" + r.URL.Query().Get("name")
		file.MimeType = r.Header.Get("Content-Type")
		file.DomainId = c.Session.DomainId
		file.Uuid = c.Params.Id

		// TODO PERMISSION
		if err := c.App.SyncUpload(r.Body, file); err != nil {
			c.Err = err
			return
		}

		sig, _ := c.App.GeneratePreSignetResourceSignature(model.AnyFileRouteName, "download", file.Id, file.DomainId)
		files = append(files, &fileResponse{
			Id:        file.Id,
			Name:      file.Name,
			Size:      file.Size,
			MimeType:  file.MimeType,
			SharedUrl: sig,
		})
	}

	if c.Err != nil {
		return
	}

	data, _ := json.Marshal(files)
	w.Write(data)

	// todo app generate public download
}
