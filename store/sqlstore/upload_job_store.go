package sqlstore

import (
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
)

type SqlUploadJobStore struct {
	SqlStore
}

const uploadJobTableName = "upload_file_jobs"

func NewSqlUploadJobStore(sqlStore SqlStore) store.UploadJobStore {
	us := &SqlUploadJobStore{sqlStore}
	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.JobUploadFile{}, uploadJobTableName).SetKeys(true, "id")
		table.ColMap("Name").SetNotNull(true).SetMaxSize(100)
		table.ColMap("Uuid").SetNotNull(true).SetMaxSize(36)
		table.ColMap("Domain").SetNotNull(true).SetMaxSize(100)
		table.ColMap("MimeType").SetNotNull(false).SetMaxSize(36)
		table.ColMap("Size").SetNotNull(true)
		table.ColMap("EmailMsg").SetNotNull(false).SetMaxSize(500)
		table.ColMap("EmailSub").SetNotNull(false).SetMaxSize(150)
		table.ColMap("Instance").SetNotNull(false).SetMaxSize(10)
		table.ColMap("FailedCount").SetNotNull(true)
	}
	return us
}

func (self *SqlUploadJobStore) CreateIndexesIfNotExists() {

}

func (self *SqlUploadJobStore) Save(job *model.JobUploadFile) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if err := self.GetMaster().Insert(job); err != nil {
			result.Err = model.NewAppError("SqlUploadJobStore.Save", "store.sql_upload_job.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (self *SqlUploadJobStore) List(limit int, instance string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var jobs []*model.JobUploadFile

		res, err := self.GetMaster().Query("SELECT id, name, uuid, domain, mime_type, size, email_msg, email_sub, instance, failed_count FROM upload_file_jobs LIMIT $1", limit)
		if err != nil {
			result.Err = model.NewAppError("SqlUploadJobStore.List", "store.sql_upload_job.list.app_error", nil, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Close()

		for res.Next() {
			job := new(model.JobUploadFile)
			err = res.Scan(&job.Id, &job.Name, &job.Uuid, &job.Domain, &job.MimeType, &job.Size, &job.EmailMsg, &job.EmailSub, &job.Instance, &job.FailedCount)
			if err != nil {
				result.Err = model.NewAppError("SqlUploadJobStore.List", "store.sql_upload_job.list.app_error", nil, err.Error(), http.StatusInternalServerError)
				return
			}
			jobs = append(jobs, job)
		}
		result.Data = jobs
	})
}
