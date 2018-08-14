package sqlstore

import (
	"fmt"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
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
		table.ColMap("MimeType").SetNotNull(false).SetMaxSize(20)
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
			result.Err = model.NewAppError("SqlMediaFileStore.Save", "store.sql_media_file.save.saving.app_error", nil,
				fmt.Sprintf("name=%s, %s", file.Name, err.Error()), http.StatusInternalServerError)
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
