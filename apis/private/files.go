package private

import (
	"fmt"
	"github.com/webitel/storage/model"
	"net/http"
)

func (api *API) InitFile() {
	api.Routes.Files.Handle("", api.ApiHandler(putFile)).Methods("PUT")
}

//  /sys/records?
// domain=10.10.10.144
// &id=65d252ab-3f9d-4293-b680-0728bb566acc
// &type=mp3
// &email=none
// &name=recordSession
// &email_sbj=none
// &email_msg=none
func putFile(c *Context, w http.ResponseWriter, r *http.Request) {
	var fileRequest model.JobUploadFile

	fileRequest.Domain = r.URL.Query().Get("domain")
	fileRequest.Uuid = r.URL.Query().Get("id")
	fileRequest.Name = fmt.Sprintf("%s.%s", r.URL.Query().Get("name"), r.URL.Query().Get("type"))
	fileRequest.MimeType = r.Header.Get("Content-Type")

	if r.URL.Query().Get("email_msg") != "" && r.URL.Query().Get("email_msg") != "none" {
		fileRequest.EmailMsg = r.URL.Query().Get("email_msg")
		fileRequest.EmailSub = r.URL.Query().Get("email_sbj")
	}

	defer r.Body.Close()

	err := c.App.AddUploadJobFile(r.Body, &fileRequest)
	if err != nil {
		c.Err = err
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"status\": \"+OK\"}"))
}
