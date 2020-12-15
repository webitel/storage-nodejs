package private

import (
	"github.com/webitel/storage/model"
	tts2 "github.com/webitel/storage/tts"
	"io"
	"net/http"
)

func (api *API) InitTTS() {
	api.Routes.TTS.Handle("/polly", api.ApiHandler(ttsPolly)).Methods("GET")
	api.Routes.TTS.Handle("/microsoft", api.ApiHandler(ttsMicrosoft)).Methods("GET")
}

func ttsPolly(c *Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := tts2.TTSParams{
		Key:      query.Get("key"),
		Token:    query.Get("token"),
		Format:   query.Get("format"),
		Voice:    query.Get("voice"),
		Region:   query.Get("region"),
		Text:     query.Get("text"),
		TextType: query.Get("text_type"),
	}

	out, t, err := tts2.Poly(params)
	if err != nil {
		c.Err = model.NewAppError("TTS", "tts.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	defer out.Close()

	if t != nil {
		w.Header().Set("Content-Type", *t)
	}
	io.Copy(w, out)
}

func ttsMicrosoft(c *Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := tts2.TTSParams{
		Key:      query.Get("key"),
		Token:    query.Get("token"),
		Format:   query.Get("format"),
		Voice:    query.Get("voice"),
		Region:   query.Get("region"),
		Text:     query.Get("text"),
		TextType: query.Get("text_type"),
	}

	out, t, err := tts2.Microsoft(params)
	if err != nil {
		c.Err = model.NewAppError("TTS", "tts.app_error", nil, err.Error(), http.StatusInternalServerError)
		return
	}

	defer out.Close()

	if t != nil {
		w.Header().Set("Content-Type", *t)
	}
	io.Copy(w, out)
}