package private

import (
	"github.com/webitel/storage/model"
	tts2 "github.com/webitel/storage/tts"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (api *API) InitTTS() {
	api.Routes.TTS.Handle("/polly", api.ApiHandler(ttsPolly)).Methods("GET")
	api.Routes.TTS.Handle("/microsoft", api.ApiHandler(ttsMicrosoft)).Methods("GET")
	api.Routes.TTS.Handle("/google", api.ApiHandler(ttsGoogle)).Methods("GET")
	api.Routes.TTS.Handle("/yandex", api.ApiHandler(ttsYandex)).Methods("GET")
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
		Language: query.Get("language"),
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

func ttsGoogle(c *Context, w http.ResponseWriter, r *http.Request) {
	var tmp string
	query := r.URL.Query()

	params := tts2.TTSParams{
		Key:      query.Get("key"),
		Token:    query.Get("token"),
		Format:   query.Get("format"),
		Voice:    query.Get("voice"),
		Region:   query.Get("region"),
		Text:     query.Get("text"),
		TextType: query.Get("text_type"),
		Language: query.Get("language"),
	}

	if tmp = query.Get("speakingRate"); tmp != "" {
		params.SpeakingRate, _ = strconv.ParseFloat(tmp, 32)
	}

	if tmp = query.Get("pitch"); tmp != "" {
		params.Pitch, _ = strconv.ParseFloat(tmp, 32)
	}

	if tmp = query.Get("volumeGainDb"); tmp != "" {
		params.VolumeGainDb, _ = strconv.ParseFloat(tmp, 32)
	}

	if tmp = query.Get("effectsProfileId"); tmp != "" {
		params.EffectsProfileId = strings.Split(tmp, ",")
	}

	params.KeyLocation = query.Get("keyLocation")

	out, t, err := tts2.Google(params)
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

func ttsYandex(c *Context, w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	params := tts2.TTSParams{
		Key:      query.Get("key"),
		Token:    query.Get("token"),
		Format:   query.Get("format"),
		Voice:    query.Get("voice"),
		Region:   query.Get("region"),
		Text:     query.Get("text"),
		TextType: query.Get("text_type"),
		Language: query.Get("language"),
	}

	out, t, err := tts2.Yandex(params)
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
