package store

import (
	"context"
)

type LayeredStoreDatabaseLayer interface {
	LayeredStoreSupplier
	StoreData
}

type LayeredStore struct {
	TmpContext    context.Context
	DatabaseLayer LayeredStoreDatabaseLayer
	ElasticLayer  *ElasticSupplier

	CdrSupplier    *CdrSupplier
	LayerChainHead LayeredStoreSupplier
}

func NewLayeredStore(db LayeredStoreDatabaseLayer) Store {
	store := &LayeredStore{
		TmpContext:    context.TODO(),
		DatabaseLayer: db,
		ElasticLayer:  NewElasticSupplier(),
	}

	store.CdrSupplier = &CdrSupplier{
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

func (s *LayeredStore) Session() SessionStore {
	return s.DatabaseLayer.Session()
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

func (s *LayeredStore) Cdr() CdrStoreData {
	return s.DatabaseLayer.Cdr()
}

func (s *LayeredStore) Schedule() ScheduleStore {
	return s.DatabaseLayer.Schedule()
}
