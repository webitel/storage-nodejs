package app

import (
	"fmt"
	"github.com/webitel/storage/model"
)

func (app *App) Version() string {
	return fmt.Sprintf("%s [build:%s]", model.CurrentVersion, model.BuildNumber)
}
