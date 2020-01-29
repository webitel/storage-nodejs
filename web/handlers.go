package web

import (
	"fmt"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/utils"
	"github.com/webitel/wlog"
	"net/http"
)

type Handler struct {
	App            *app.App
	Ctrl           *controller.Controller
	HandleFunc     func(*Context, http.ResponseWriter, *http.Request)
	RequireSession bool
	TrustRequester bool
	RequireMfa     bool
	IsStatic       bool
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wlog.Debug(fmt.Sprintf("%v - %v", r.Method, r.URL.Path))

	c := &Context{}
	c.App = h.App
	c.Ctrl = h.Ctrl
	c.T, _ = utils.GetTranslationsAndLocale(w, r)
	c.Params = ParamsFromRequest(r)
	c.RequestId = model.NewId()
	c.IpAddress = utils.GetIpAddress(r)
	c.Path = r.URL.Path
	c.Log = c.App.Log

	w.Header().Set(model.HEADER_REQUEST_ID, c.RequestId)
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		w.Header().Set("Expires", "0")
	}

	//TODO
	token, _ := app.ParseAuthTokenFromRequest(r)
	if len(token) != 0 && h.RequireSession {
		session, err := c.App.GetSession(token)
		if err != nil {
			c.Log.Info("Invalid session", wlog.Err(err))
			if err.StatusCode == http.StatusInternalServerError {
				c.Err = err
			} else {
				c.Err = model.NewAppError("ServeHTTP", "api.context.session_expired.app_error", nil, "token="+token, http.StatusUnauthorized)
			}
		} else {
			c.Session = *session
		}
	}

	c.Log = c.App.Log.With(
		wlog.String("path", c.Path),
		wlog.String("request_id", c.RequestId),
		wlog.String("ip_addr", c.IpAddress),
		wlog.Int64("user_id", c.Session.UserId),
		wlog.String("method", r.Method),
	)

	if c.Err == nil && h.RequireSession {
		c.SessionRequired()
	}

	if c.Err == nil {
		h.HandleFunc(c, w, r)
	}

	// Handle errors that have occurred
	if c.Err != nil {
		c.Err.Translate(c.T)
		c.Err.RequestId = c.RequestId

		if c.Err.Id == "api.context.session_expired.app_error" {
			c.LogInfo(c.Err)
		} else {
			c.LogError(c.Err)
		}

		c.Err.Where = r.URL.Path

		w.WriteHeader(c.Err.StatusCode)
		w.Write([]byte(c.Err.ToJson()))
	}
}
