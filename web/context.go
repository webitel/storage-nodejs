package web

import (
	goi18n "github.com/nicksnyder/go-i18n/i18n"
	"github.com/webitel/engine/auth_manager"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/utils"
	"github.com/webitel/wlog"
	"net/http"
)

type Context struct {
	App           *app.App
	Log           *wlog.Logger
	Session       auth_manager.Session
	Err           *model.AppError
	T             goi18n.TranslateFunc
	Params        *Params
	Ctrl          *controller.Controller
	RequestId     string
	IpAddress     string
	Path          string
	siteURLHeader string
}

func (c *Context) LogError(err *model.AppError) {
	// Filter out 404s, endless reconnects and browser compatibility errors
	if err.StatusCode == http.StatusNotFound {
		c.LogDebug(err)
	} else {
		c.Log.Error(
			err.SystemMessage(utils.TDefault),
			wlog.String("err_where", err.Where),
			wlog.Int("http_code", err.StatusCode),
			wlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogInfo(err *model.AppError) {
	// Filter out 401s
	if err.StatusCode == http.StatusUnauthorized {
		c.LogDebug(err)
	} else {
		c.Log.Info(
			err.SystemMessage(utils.TDefault),
			wlog.String("err_where", err.Where),
			wlog.Int("http_code", err.StatusCode),
			wlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogDebug(err *model.AppError) {
	c.Log.Debug(
		err.SystemMessage(utils.TDefault),
		wlog.String("err_where", err.Where),
		wlog.Int("http_code", err.StatusCode),
		wlog.String("err_details", err.DetailedError),
	)
}

func (c *Context) SessionRequired() {
	if c.Session.UserId == 0 {
		c.Err = model.NewAppError("", "api.context.session_expired.app_error", nil, "UserRequired", http.StatusUnauthorized)
		return
	}
}

func (c *Context) SetInvalidParam(parameter string) {
	c.Err = NewInvalidParamError(parameter)
}

func (c *Context) SetInvalidUrlParam(parameter string) {
	c.Err = NewInvalidUrlParamError(parameter)
}

func NewInvalidParamError(parameter string) *model.AppError {
	err := model.NewAppError("Context", "api.context.invalid_body_param.app_error", map[string]interface{}{"Name": parameter}, "", http.StatusBadRequest)
	return err
}

func NewInvalidUrlParamError(parameter string) *model.AppError {
	err := model.NewAppError("Context", "api.context.invalid_url_param.app_error", map[string]interface{}{"Name": parameter}, "", http.StatusBadRequest)
	return err
}

func (c *Context) RequireId() *Context {
	if c.Err != nil {
		return c
	}

	if len(c.Params.Id) == 0 {
		c.SetInvalidUrlParam("id")
	}
	return c
}

func (c *Context) RequireDomain() *Context {
	if c.Err != nil {
		return c
	}

	if len(c.Params.Domain) == 0 {
		c.SetInvalidUrlParam("domain_id")
	}

	return c
}

func (c *Context) RequireExpire() *Context {
	if c.Err != nil {
		return c
	}

	if c.Params.Expires < 1 {
		c.SetInvalidUrlParam("expires")
	}
	return c
}

func (c *Context) SetSessionExpire() {
	c.Err = model.NewAppError("ServeHTTP", "api.context.session_expired.app_error", nil, "", http.StatusUnauthorized)
}

func (c *Context) SetSessionErrSignature() {
	c.Err = model.NewAppError("ServeHTTP", "api.context.session_signature.app_error", nil, "", http.StatusUnauthorized)
}

func (c *Context) RequireSignature() *Context {
	if c.Err != nil {
		return c
	}

	if len(c.Params.Signature) == 0 {
		c.SetInvalidUrlParam("signature")
	}
	return c
}
