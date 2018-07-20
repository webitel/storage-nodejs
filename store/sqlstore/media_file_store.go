package sqlstore

import (
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
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

		table.ColMap("ProfileId").SetTransient(true)
		table.ColMap("Uuid").SetTransient(true)
	}

	return us
}
