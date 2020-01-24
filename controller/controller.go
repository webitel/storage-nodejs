package controller

import "github.com/webitel/storage/app"

type Controller struct {
	app *app.App
}

func NewController(a *app.App) *Controller {
	return &Controller{a}
}
