package file_sync

import (
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/interfaces"
)

type SyncFilesJobInterfaceImpl struct {
	App *app.App
}

func init() {
	app.RegisterSyncFilesJobInterface(func(a *app.App) interfaces.SyncFilesJobInterface {
		return &SyncFilesJobInterfaceImpl{a}
	})
}
