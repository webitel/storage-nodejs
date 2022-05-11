package main

import (
	"fmt"

	"github.com/webitel/storage/apis"
	"github.com/webitel/storage/app"
	"github.com/webitel/storage/grpc_api"
	"github.com/webitel/wlog"

	_ "github.com/webitel/storage/stt"
	_ "github.com/webitel/storage/synchronizer"
	_ "github.com/webitel/storage/uploader"

	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/webitel/storage/apis/private"
)

func main() {
	interruptChan := make(chan os.Signal, 1)
	a, err := app.New()
	if err != nil {
		panic(err.Error())
	}
	wlog.Info(fmt.Sprintf("server version: %s", a.Version()))

	serverErr := a.StartServer()
	if serverErr != nil {
		wlog.Critical(serverErr.Error())
		return
	}
	apis.Init(a, a.Srv.Router)

	serverErr = a.StartInternalServer()
	if serverErr != nil {
		wlog.Critical(serverErr.Error())
		return
	}
	private.Init(a, a.InternalSrv.Router)

	a.Uploader.Start()
	a.Synchronizer.Start()

	grpc_api.Init(a, a.GrpcServer.Server())

	if err = a.StartGrpcServer(); err != nil {
		panic(err.Error())
	}

	setDebug()
	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	a.Shutdown()

	//a.Broker.Close()

	wlog.Info("Stopping synchronizer server")
	a.Synchronizer.Stop()

	wlog.Info("Stopping uploader server")
	a.Uploader.Stop()

}

func setDebug() {
	//debug.SetGCPercent(-1)

	go func() {
		wlog.Info("Start debug server on :8090")
		err := http.ListenAndServe(":8090", nil)
		if err != nil {
			wlog.Critical(err.Error())
		}
	}()

}
