package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
)

type SqlFileStore struct {
	SqlStore
}

func NewSqlFileStore(sqlStore SqlStore) store.FileStore {
	us := &SqlFileStore{sqlStore}

	return us
}

func (self SqlFileStore) CreateIndexesIfNotExists() {

}

//TODO reference tables ?
func (self *SqlFileStore) MoveFromJob(jobId, profileId int, properties model.StringInterface) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		_, err := self.GetMaster().Exec(`with del as (
  delete from upload_file_jobs
  where id = $1
  returning name, uuid, size, domain, mime_type, created_at, instance
)
insert into files(name, uuid, profile_id, size, domain, mime_type, properties, created_at, instance)
select del.name, del.uuid, $2, del.size, del.domain, del.mime_type, $3, del.created_at, del.instance
from del`, jobId, profileId, properties.ToJson())

		if err != nil {
			result.Err = model.NewAppError("SqlFileStore.MoveFromJob", "store.sql_file.move_from_job.app_error", nil, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (self *SqlFileStore) GetAllPageByDomain(domain string, offset, limit int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var recordings []*model.File

		query := `SELECT * FROM files 
			WHERE domain = :Domain  
			LIMIT :Limit OFFSET :Offset`

		if _, err := self.GetReplica().Select(&recordings, query, map[string]interface{}{"Domain": domain, "Offset": offset, "Limit": limit}); err != nil {
			result.Err = model.NewAppError("SqlFileStore.List", "store.sql_file.get_all.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = recordings
		}
	})
}

func (self *SqlFileStore) Get(domain, uuid string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var files []*model.FileWithProfile
		_, err := self.GetReplica().Select(&files, `SELECT files.*, p.updated_at as profile_updated_at FROM files JOIN file_backend_profiles p on p.id = files.profile_id WHERE uuid = :Uuid AND files.domain = :Domain`,
			map[string]interface{}{"Domain": domain, "Uuid": uuid})

		if err != nil {
			result.Err = model.NewAppError("SqlFileStore.Get", "store.sql_file.get.app_error", nil, "uuid="+uuid+" "+err.Error(), http.StatusInternalServerError)
			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = files
		}
	})
}

func (self *SqlFileStore) Delete(domain string, id int64) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {

	})
}

func (self *SqlFileStore) DeleteById(id int64) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := self.GetMaster().Exec(
			`DELETE FROM
				files
			WHERE
				id = :Id`, map[string]interface{}{"Id": id}); err != nil {
			result.Err = model.NewAppError("SqlFileStore.DeleteById", "store.sql_file.delete.app_error", nil,
				fmt.Sprintf("id=%d, err: %s", err.Error()), http.StatusInternalServerError)
		} else {
			result.Data = id
		}
	})
}

func (self *SqlFileStore) FetchDeleted(limit int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var recordings []*model.FileWithProfile

		query := `SELECT files.*,  p.updated_at as profile_updated_at FROM files JOIN file_backend_profiles p on p.id = files.profile_id 
			WHERE removed is TRUE AND NOT not_exists is TRUE 
			LIMIT :Limit `

		if _, err := self.GetReplica().Select(&recordings, query, map[string]interface{}{"Limit": limit}); err != nil {
			result.Err = model.NewAppError("SqlFileStore.FetchDeleted", "store.sql_file.get_deleted.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = recordings
		}
	})
}

func (self *SqlFileStore) SetNoExistsById(id int64, notExists bool) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if _, err := self.GetMaster().Exec(
			`update files
				set not_exists = :NotExists 
				where where id = :Id`, map[string]interface{}{"Id": id, "NotExists": notExists}); err != nil {
			result.Err = model.NewAppError("SqlFileStore.SetNoExistsById", "store.sql_file.update_exists.app_error", nil,
				fmt.Sprintf("id=%d, err: %s", err.Error()), http.StatusInternalServerError)
		} else {
			result.Data = id
		}
	})
}
