package store

import (
	"context"
)

type LayeredStoreDatabaseLayer interface {
	LayeredStoreSupplier
	Store
}

type LayeredStore struct {
	TmpContext     context.Context
	DatabaseLayer  LayeredStoreDatabaseLayer
	LayerChainHead LayeredStoreSupplier
}

func NewLayeredStore(db LayeredStoreDatabaseLayer) Store {
	store := &LayeredStore{
		TmpContext:    context.TODO(),
		DatabaseLayer: db,
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

func (s *LayeredStore) Recording() RecordingStore {
	return s.DatabaseLayer.Recording()
}

func (s *LayeredStore) Job() JobStore {
	return s.DatabaseLayer.Job()
}
