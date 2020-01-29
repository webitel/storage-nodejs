package grpc_api

import (
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/controller"
	"github.com/webitel/storage/grpc_api/storage"
	"google.golang.org/grpc"
)

type API struct {
	app             *app.App
	ctrl            *controller.Controller
	backendProfiles *backendProfiles
	media           *media
}

func Init(a *app.App, server *grpc.Server) {
	api := &API{
		app: a,
	}

	ctrl := controller.NewController(a)
	api.backendProfiles = NewBackendProfileApi(ctrl)
	api.media = NewMediaApi(ctrl)

	storage.RegisterBackendProfileServiceServer(server, api.backendProfiles)
	storage.RegisterMediaFileServiceServer(server, api.media)
}
