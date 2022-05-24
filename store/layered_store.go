package store

import (
	"context"

	"github.com/webitel/storage/model"
)

type LayeredStoreDatabaseLayer interface {
	LayeredStoreSupplier
	StoreData
}

type LayeredStore struct {
	TmpContext    context.Context
	DatabaseLayer LayeredStoreDatabaseLayer
	ElasticLayer  *ElasticSupplier

	SearchEngine   *SearchEngine
	LayerChainHead LayeredStoreSupplier
}

func NewLayeredStore(db LayeredStoreDatabaseLayer, noSql *ElasticSupplier) Store {
	store := &LayeredStore{
		TmpContext:    context.TODO(),
		DatabaseLayer: db,
		ElasticLayer:  noSql,
	}

	store.SearchEngine = &SearchEngine{
		DatabaseLayer: store.DatabaseLayer,
		ElasticLayer:  store.ElasticLayer,
	}

	return store
}

type QueryFunction func(LayeredStoreSupplier) *LayeredStoreSupplierResult

func (s *LayeredStore) RunQuery(queryFunction QueryFunction) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		result := queryFunction(s.LayerChainHead)
		storeChannel <- result.StoreResult
	}()

	return storeChannel
}

func (s *LayeredStore) UploadJob() UploadJobStore {
	return s.DatabaseLayer.UploadJob()
}

func (s *LayeredStore) FileBackendProfile() FileBackendProfileStore {
	return s.DatabaseLayer.FileBackendProfile()
}

func (s *LayeredStore) File() FileStore {
	return s.DatabaseLayer.File()
}

func (s *LayeredStore) Job() JobStore {
	return s.DatabaseLayer.Job()
}

func (s *LayeredStore) MediaFile() MediaFileStore {
	return s.DatabaseLayer.MediaFile()
}

func (s *LayeredStore) Schedule() ScheduleStore {
	return s.DatabaseLayer.Schedule()
}

func (s *LayeredStore) SyncFile() SyncFileStore {
	return s.DatabaseLayer.SyncFile()
}

func (s *LayeredStore) CognitiveProfile() CognitiveProfileStore {
	return s.DatabaseLayer.CognitiveProfile()
}

func (s *LayeredStore) Scroll(scroll *model.SearchEngineScroll) StoreChannel {
	return s.SearchEngine.Scroll(scroll)
}

func (s *LayeredStore) Search(request *model.SearchEngineRequest) StoreChannel {
	return s.SearchEngine.Search(request)
}

func (s *LayeredStore) TranscriptFile() TranscriptFileStore {
	return s.DatabaseLayer.TranscriptFile()
}
