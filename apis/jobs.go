package apis

import (
	"fmt"
	"io"
	"net/http"
)

func (api *API) InitJobs() {
	api.PublicRoutes.Jobs.Handle("/callback", api.ApiHandler(callbackJob)).Methods("POST")
}

func callbackJob(c *Context, w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)
	fmt.Println(string(data))
	fmt.Println(r.Header)
	fmt.Println(r.URL.Query())

	switch r.Header.Get("X-Microsoftspeechservices-Event") {
	case "Challenge":
		w.Write([]byte(r.URL.Query().Get("validationToken")))
		return
	case "TranscriptionCompletion":

		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
}
