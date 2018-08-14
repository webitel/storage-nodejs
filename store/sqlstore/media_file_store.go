package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
	"strings"
)

type SqlMediaFileStore struct {
	SqlStore
}

func NewSqlMediaFileStore(sqlStore SqlStore) store.MediaFileStore {
	us := &SqlMediaFileStore{sqlStore}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.MediaFile{}, "media_files").SetKeys(true, "id")
		table.ColMap("Name").SetNotNull(true).SetMaxSize(100)
		table.ColMap("Size").SetNotNull(true)
		table.ColMap("Domain").SetNotNull(true).SetMaxSize(100)
		table.ColMap("MimeType").SetNotNull(false).SetMaxSize(40)
		table.ColMap("Instance").SetNotNull(false).SetMaxSize(20)
	}

	return us
}

func (self SqlMediaFileStore) CreateIndexesIfNotExists() {

}

func (self *SqlMediaFileStore) Save(file *model.MediaFile) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		file.PreSave()
		if err := self.GetMaster().Insert(file); err != nil {
			//TODO
			if strings.Index(err.Error(), "duplicate") > -1 {
				result.Err = model.NewAppError("SqlMediaFileStore.Save", "store.sql_media_file.save.saving.duplicate", nil,
					fmt.Sprintf("name=%s, %s", file.Name, err.Error()), http.StatusInternalServerError)
			} else {
				result.Err = model.NewAppError("SqlMediaFileStore.Save", "store.sql_media_file.save.saving.app_error", nil,
					fmt.Sprintf("name=%s, %s", file.Name, err.Error()), http.StatusInternalServerError)
			}
		} else {
			result.Data = file
		}
	})
}

func (self *SqlMediaFileStore) GetAllByDomain(domain string, offset, limit int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var files []*model.MediaFile

		query := `SELECT * FROM media_files 
			WHERE domain = :Domain  
			LIMIT :Limit OFFSET :Offset`

		if _, err := self.GetReplica().Select(&files, query, map[string]interface{}{"Domain": domain, "Offset": offset, "Limit": limit}); err != nil {
			result.Err = model.NewAppError("SqlMediaFileStore.List", "store.sql_media_file.get_all.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = files
		}
	})
}

func (self *SqlMediaFileStore) GetCountByDomain(domain string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		query := `SELECT count(*) FROM media_files 
			WHERE domain = :Domain`

		if count, err := self.GetReplica().SelectInt(query, map[string]interface{}{"Domain": domain}); err != nil {
			result.Err = model.NewAppError("SqlMediaFileStore.LiGetCountByDomainst", "store.sql_media_file.get_all.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = count
		}
	})
}

func (self *SqlMediaFileStore) Get(id int64, domain string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		query := `SELECT * FROM media_files WHERE id = :Id AND domain = :Domain`

		file := &model.MediaFile{}

		if err := self.GetReplica().SelectOne(file, query, map[string]interface{}{"Id": id, "Domain": domain}); err != nil {
			result.Err = model.NewAppError("SqlMediaFileStore.Get", "store.sql_media_file.get.finding.app_error", nil,
				fmt.Sprintf("id=%d, domain=%s, %s", id, domain, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = file
		}
	})
}

func (self *SqlMediaFileStore) GetByName(name, domain string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		query := `SELECT * FROM media_files WHERE name = :Name AND domain = :Domain`

		file := &model.MediaFile{}

		if err := self.GetReplica().SelectOne(file, query, map[string]interface{}{"Name": name, "Domain": domain}); err != nil {
			result.Err = model.NewAppError("SqlMediaFileStore.GetByName", "store.sql_media_file.get_by_name.finding.app_error", nil,
				fmt.Sprintf("name=%s, domain=%s, %s", name, domain, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = file
		}
	})
}

func (self *SqlMediaFileStore) DeleteByName(name, domain string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		res, err := self.GetMaster().Exec("DELETE FROM media_files WHERE name = :Name AND domain = :Domain", map[string]interface{}{"Name": name, "Domain": domain})
		if err != nil {
			result.Err = model.NewAppError("SqlMediaFileStore.DeleteByName", "store.sql_media_file.delete.app_error", nil,
				fmt.Sprintf("name=%s, err: %s", name, err.Error()), http.StatusInternalServerError)
			return
		}
		count, _ := res.RowsAffected()
		if count == 0 {
			result.Err = model.NewAppError("SqlMediaFileStore.DeleteByName", "store.sql_media_file.delete.not_found.app_error", map[string]interface{}{"Name": name},
				"", http.StatusNotFound)
		}
	})
}

func (self *SqlMediaFileStore) DeleteById(id int64) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		res, err := self.GetMaster().Exec("DELETE FROM media_files WHERE id = :Id", map[string]interface{}{"Id": id})
		if err != nil {
			result.Err = model.NewAppError("SqlMediaFileStore.DeleteById", "store.sql_media_file.delete.app_error", nil,
				fmt.Sprintf("id=%d, err: %s", id, err.Error()), http.StatusInternalServerError)
			return
		}
		count, _ := res.RowsAffected()
		if count == 0 {
			result.Err = model.NewAppError("SqlMediaFileStore.DeleteById", "store.sql_media_file.delete.not_found.app_error", map[string]interface{}{"Id": id},
				"", http.StatusNotFound)
		}
	})
}
