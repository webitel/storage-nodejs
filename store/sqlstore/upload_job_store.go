package sqlstore

import (
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
	"strings"
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
		table.ColMap("CreatedAt").SetNotNull(true)
		table.ColMap("Attempts").SetNotNull(true)
	}
	return us
}

func (self *SqlUploadJobStore) CreateIndexesIfNotExists() {

}

func (self *SqlUploadJobStore) Save(job *model.JobUploadFile) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		job.PreSave()
		if err := self.GetMaster().Insert(job); err != nil {
			result.Err = model.NewAppError("SqlUploadJobStore.Save", "store.sql_upload_job.save.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (self *SqlUploadJobStore) GetAllPageByInstance(limit int, instance string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var jobs []*model.JobUploadFile

		res, err := self.GetReplica().Query("SELECT id, name, uuid, domain, mime_type, size, email_msg, email_sub, instance, attempts FROM upload_file_jobs LIMIT $1", limit)
		if err != nil {
			result.Err = model.NewAppError("SqlUploadJobStore.List", "store.sql_upload_job.list.app_error", nil, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Close()

		for res.Next() {
			job := new(model.JobUploadFile)
			err = res.Scan(&job.Id, &job.Name, &job.Uuid, &job.Domain, &job.MimeType, &job.Size, &job.EmailMsg, &job.EmailSub, &job.Instance, &job.Attempts)
			if err != nil {
				result.Err = model.NewAppError("SqlUploadJobStore.List", "store.sql_upload_job.list.app_error", nil, err.Error(), http.StatusInternalServerError)
				return
			}
			jobs = append(jobs, job)
		}
		result.Data = jobs
	})
}

//region  sql sqlUpdateWithProfile
const sqlUpdateWithProfile = `update upload_file_jobs uu
set attempts = attempts + 1
  ,state = 1
  ,updated_at = extract(EPOCH from now()) :: BIGINT
from (
       SELECT
         t.id,
         t.name,
         t.uuid,
         t.domain_id,
         t.mime_type,
         t.size,
         t.email_msg,
         t.email_sub,
         profile.id as profile_id,
         profile.updated_at
       FROM upload_file_jobs as t
         inner join lateral (              select
                                             tmp.domain_id,
                                             tmp.id,
                                             tmp.updated_at
                                           from (select
                                                   p1.domain_id,
                                                   p1.id,
                                                   p1.updated_at,
                                                   p1.priority
                                                 from file_backend_profiles p1
                                                 where p1.domain_id = t.domain_id and NOT p1.disabled is TRUE
                                                 union all select
                                                             t.domain_id,
                                                             p2.id,
                                                             p2.updated_at,
                                                             p2.priority
                                                           from file_backend_profiles p2
                                                           where p2.domain_id = $3 and
                                                                 NOT p2.disabled is TRUE) as tmp
                                           order by tmp.priority asc
                                           FETCH FIRST 1 ROW ONLY              ) profile ON profile.domain_id = t.domain_id
       WHERE state = 0 AND instance = $2 AND (t.updated_at < $4 OR attempts = 0)
       ORDER BY created_at ASC
       LIMIT $1) tmp
WHERE tmp.id = uu.id and state = 0
returning tmp.*`

var sqlUpdateWithProfileFmt = strings.Replace(sqlUpdateWithProfile, "\n", " ", -1)

//endregion

func (self *SqlUploadJobStore) UpdateWithProfile(limit int, instance string, betweenAttemptSec int64) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var jobs = make([]*model.JobUploadFileWithProfile, 0, limit)

		rows, err := self.GetMaster().Query(sqlUpdateWithProfileFmt, limit, instance, model.ROOT_FILE_BACKEND_DOMAIN, model.GetMillis()-betweenAttemptSec)
		if err != nil {
			result.Err = model.NewAppError("SqlUploadJobStore.UpdateWithProfile", "store.sql_upload_job.update_with_profile.app_error", nil, err.Error(), http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		for rows.Next() {
			job := &model.JobUploadFileWithProfile{}
			err = rows.Scan(&job.Id, &job.Name, &job.Uuid, &job.Domain, &job.MimeType, &job.Size, &job.EmailMsg, &job.EmailSub, &job.ProfileId, &job.ProfileUpdatedAt)
			if err != nil {
				result.Err = model.NewAppError("SqlUploadJobStore.UpdateWithProfile", "store.sql_upload_job.update_with_profile.scan.app_error", nil, err.Error(), http.StatusInternalServerError)
				return
			}
			jobs = append(jobs, job)
		}
		result.Data = jobs
	})
}

func (self *SqlUploadJobStore) SetStateError(id int, errMsg string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		self.GetMaster().Exec(`update upload_file_jobs
set state = 0,
  updated_at = $2
where id = $1`, id, model.GetMillis())
	})
}

func (self *SqlUploadJobStore) RemoveById(id int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		panic("TODO")
	})
}
