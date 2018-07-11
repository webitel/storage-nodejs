package sqlstore

import (
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
)

type SqlFileBackendProfileStore struct {
	SqlStore
}

func NewSqlFileBackendProfileStore(sqlStore SqlStore) store.FileBackendProfileStore {
	us := &SqlFileBackendProfileStore{sqlStore}
	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.FileBackendProfile{}, "file_backend_profiles").SetKeys(true, "id")
		table.ColMap("Name").SetNotNull(true).SetMaxSize(100)
		table.ColMap("Domain").SetNotNull(true).SetMaxSize(100)
		table.ColMap("Default").SetNotNull(true)
		table.ColMap("ExpireDay").SetNotNull(true)
		table.ColMap("Disabled")
		table.ColMap("MaxSizeMb").SetNotNull(true)
		table.ColMap("Properties").SetNotNull(true)
	}

	return us
}

func (self *SqlFileBackendProfileStore) CreateIndexesIfNotExists() {

}
