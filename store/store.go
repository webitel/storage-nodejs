package store

import (
	"github.com/webitel/engine/auth_manager"
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
	StoreData
	StoreSearchEngine
}

type StoreData interface {
	UploadJob() UploadJobStore
	FileBackendProfile() FileBackendProfileStore
	File() FileStore
	Job() JobStore
	MediaFile() MediaFileStore
	Cdr() CdrStoreData
	Schedule() ScheduleStore
}

type StoreSearchEngine interface {
	Search(request *model.SearchEngineRequest) StoreChannel
	Scroll(scroll *model.SearchEngineScroll) StoreChannel
}

type UploadJobStore interface {
	Create(job *model.JobUploadFile) *model.AppError
	//Save(job *model.JobUploadFile) StoreChannel
	GetAllPageByInstance(limit int, instance string) StoreChannel
	UpdateWithProfile(limit int, instance string, betweenAttemptSec int64) StoreChannel
	SetStateError(id int, errMsg string) StoreChannel
}

type FileBackendProfileStore interface {
	CheckAccess(domainId, id int64, groups []int, access auth_manager.PermissionAccess) (bool, *model.AppError)
	Create(profile *model.FileBackendProfile) (*model.FileBackendProfile, *model.AppError)
	GetAllPage(domainId int64, offset, limit int) ([]*model.FileBackendProfile, *model.AppError)
	GetAllPageByGroups(domainId int64, groups []int, offset, limit int) ([]*model.FileBackendProfile, *model.AppError)
	Get(id, domainId int64) (*model.FileBackendProfile, *model.AppError)
	Update(profile *model.FileBackendProfile) (*model.FileBackendProfile, *model.AppError)
	Delete(domainId, id int64) *model.AppError

	GetById(id int) StoreChannel
	GetAllPageByDomain(domain string, limit, offset int) StoreChannel
}

type FileStore interface {
	Get(domain, uuid string) StoreChannel
	GetAllPageByDomain(domain string, offset, limit int) StoreChannel
	MoveFromJob(jobId, profileId int, properties model.StringInterface) StoreChannel
	FetchDeleted(limit int) StoreChannel
	DeleteById(id int64) StoreChannel
	SetNoExistsById(id int64, notExists bool) StoreChannel
}

type MediaFileStore interface {
	Create(file *model.MediaFile) (*model.MediaFile, *model.AppError)
	GetAllPage(domainId int64, req *model.ListRequest) ([]*model.MediaFile, *model.AppError)
	Get(domainId int64, id int) (*model.MediaFile, *model.AppError)
	Delete(domainId, id int64) *model.AppError

	Save(file *model.MediaFile) StoreChannel
	GetAllByDomain(domain string, offset, limit int) StoreChannel
	GetCountByDomain(domain string) StoreChannel
	GetByName(name, domain string) StoreChannel
	DeleteByName(name, domain string) StoreChannel
	DeleteById(id int64) StoreChannel
}

type ScheduleStore interface {
	GetAllEnablePage(limit, offset int) StoreChannel
	GetAllWithNoJobs(limit, offset int) StoreChannel
	GetAllPageByType(typeName string) StoreChannel
}

type JobStore interface {
	Save(job *model.Job) StoreChannel
	UpdateOptimistically(job *model.Job, currentStatus string) StoreChannel
	UpdateStatus(id string, status string) StoreChannel
	UpdateStatusOptimistically(id string, currentStatus string, newStatus string) StoreChannel
	Get(id string) StoreChannel
	GetAllPage(offset int, limit int) StoreChannel
	GetAllByType(jobType string) StoreChannel
	GetAllByTypePage(jobType string, offset int, limit int) StoreChannel
	GetAllByStatus(status string) StoreChannel
	GetAllByStatusAndLessScheduleTime(status string, t int64) StoreChannel
	GetNewestJobByStatusAndType(status string, jobType string) StoreChannel
	GetCountByStatusAndType(status string, jobType string) StoreChannel
	Delete(id string) StoreChannel
}

type CdrStoreData interface {
	GetLegADataByUuid(uuid string) StoreChannel
	GetLegBDataByUuid(uuid string) StoreChannel
	GetByUuidCall(uuid string) StoreChannel
}
