package jobs

import (
	"github.com/webitel/storage/interfaces"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
)

type ConfigService interface {
	Config() *model.Config
}

type StaticConfigService struct {
	Cfg *model.Config
}

func (s StaticConfigService) Config() *model.Config { return s.Cfg }

type JobServer struct {
	ConfigService ConfigService
	Store         store.Store
	Workers       *Workers
	Schedulers    *Schedulers

	middlewareJobs map[string]interfaces.JobInterface
}

func NewJobServer(configService ConfigService, store store.Store) *JobServer {
	return &JobServer{
		ConfigService:  configService,
		Store:          store,
		middlewareJobs: make(map[string]interfaces.JobInterface),
	}
}

func (srv *JobServer) AddMiddleware(name string, w interfaces.JobInterface) {
	srv.middlewareJobs[name] = w
}

func (srv *JobServer) Config() *model.Config {
	return srv.ConfigService.Config()
}

func (srv *JobServer) StartWorkers() {
	srv.Workers = srv.InitWorkers().Start()
}

func (srv *JobServer) StartSchedulers() {
	srv.Schedulers = srv.InitSchedulers().Start()
}

func (srv *JobServer) StopWorkers() {
	if srv.Workers != nil {
		srv.Workers.Stop()
	}
}

func (srv *JobServer) StopSchedulers() {
	if srv.Schedulers != nil {
		srv.Schedulers.Stop()
	}
}
