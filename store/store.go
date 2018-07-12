package store

import (
	"time"

	"github.com/webitel/storage/model"
)

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

func Do(f func(result *StoreResult)) StoreChannel {
	storeChannel := make(StoreChannel, 1)
	go func() {
		result := StoreResult{}
		f(&result)
		storeChannel <- result
		close(storeChannel)
	}()
	return storeChannel
}

func Must(sc StoreChannel) interface{} {
	r := <-sc
	if r.Err != nil {

		time.Sleep(time.Second)
		panic(r.Err)
	}

	return r.Data
}

type Store interface {
	Session() SessionStore
	UploadJob() UploadJobStore
	FileBackendProfile() FileBackendProfileStore
}

type SessionStore interface {
	Get(sessionIdOrToken string) StoreChannel
}

type UploadJobStore interface {
	Save(job *model.JobUploadFile) StoreChannel
	List(limit int, instance string) StoreChannel
}

type FileBackendProfileStore interface {
	Save(profile *model.FileBackendProfile) StoreChannel
	Get(id int, domain string) StoreChannel
	GetList(domain string, limit, offset int) StoreChannel
	Delete(id int, domain string) StoreChannel
	Update(profile *model.FileBackendProfile) StoreChannel
}
