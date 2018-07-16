package sqlstore

import (
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
)

type SqlRecordingStore struct {
	SqlStore
}

func NewSqlRecordingStore(sqlStore SqlStore) store.RecordingStore {
	us := &SqlRecordingStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Recoding{}, "recording_files").SetKeys(true, "id")
		table.ColMap("Name").SetNotNull(true).SetMaxSize(100)
		table.ColMap("Uuid").SetNotNull(true).SetMaxSize(36)
		table.ColMap("ProfileId").SetNotNull(true)
		table.ColMap("Size").SetNotNull(true)
		table.ColMap("Domain").SetNotNull(true).SetMaxSize(100)
		table.ColMap("MimeType").SetNotNull(false).SetMaxSize(20)
		table.ColMap("Properties").SetNotNull(true)
	}

	return us
}

func (self SqlRecordingStore) CreateIndexesIfNotExists() {

}

//TODO reference tables ?
func (self *SqlRecordingStore) MoveFromJob(jobId, profileId int, properties model.StringInterface) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		_, err := self.GetReplica().Exec(`with del as (
  delete from upload_file_jobs
  where id = $1
  returning name, uuid, size, domain, mime_type, created_at
)
insert into recording_files (name, uuid, profile_id, size, domain, mime_type, properties, created_at)
select del.name, del.uuid, $2, del.size, del.domain, del.mime_type, $3, del.created_at
from del`, jobId, profileId, properties.ToJson())

		if err != nil {
			result.Err = model.NewAppError("SqlRecordingStore.MoveFromJob", "store.sql_recording.move_from_job.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	})
}
