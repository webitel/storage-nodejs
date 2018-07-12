package web

import (
	goi18n "github.com/nicksnyder/go-i18n/i18n"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/mlog"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/utils"
	"net/http"
)

type Context struct {
	App           *app.App
	Log           *mlog.Logger
	Session       model.Session
	Err           *model.AppError
	T             goi18n.TranslateFunc
	Params        *Params
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
			mlog.String("err_where", err.Where),
			mlog.Int("http_code", err.StatusCode),
			mlog.String("err_details", err.DetailedError),
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
			mlog.String("err_where", err.Where),
			mlog.Int("http_code", err.StatusCode),
			mlog.String("err_details", err.DetailedError),
		)
	}
}

func (c *Context) LogDebug(err *model.AppError) {
	c.Log.Debug(
		err.SystemMessage(utils.TDefault),
		mlog.String("err_where", err.Where),
		mlog.Int("http_code", err.StatusCode),
		mlog.String("err_details", err.DetailedError),
	)
}

func (c *Context) SessionRequired() {
	if len(c.Session.UserId) == 0 {
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
		c.SetInvalidUrlParam("domain")
	}
	return c
}