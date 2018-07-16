package jobs

import (
	tjobs "github.com/webitel/storage/jobs/interfaces"
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
	Workers       []*Workers
	Schedulers    *Schedulers

	UploadRecordingsJob tjobs.UploadRecordingsFilesJobInterface
}

func NewJobServer(configService ConfigService, store store.Store) *JobServer {
	return &JobServer{
		ConfigService: configService,
		Store:         store,
		Workers:       make([]*Workers, 5, 5),
	}
}

func (srv *JobServer) Config() *model.Config {
	return srv.ConfigService.Config()
}

func (srv *JobServer) StopWorkers() {
	if srv.Workers != nil {
		for _, worker := range srv.Workers {
			worker.Stop()
		}
	}
}

func (srv *JobServer) StartSchedulers() {
	srv.Schedulers = srv.InitSchedulers().Start()
}

func (srv *JobServer) StopSchedulers() {
	if srv.Schedulers != nil {
		srv.Schedulers.Stop()
	}
}
