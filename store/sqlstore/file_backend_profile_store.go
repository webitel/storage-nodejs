package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
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
		table.ColMap("TypeId").SetNotNull(true)
		table.ColMap("CreatedAt").SetNotNull(true)
		table.ColMap("UpdatedAt").SetNotNull(true)

		table = db.AddTableWithName(model.FileBackendProfileType{}, "file_backend_profile_type").SetKeys(true, "id")
		table.ColMap("Name").SetMaxSize(50)
		table.ColMap("Code").SetMaxSize(10)
	}

	return us
}

func (self SqlFileBackendProfileStore) CreateIndexesIfNotExists() {

}

func (s SqlFileBackendProfileStore) Save(profile *model.FileBackendProfile) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		profile.PreSave()
		if err := s.GetMaster().Insert(profile); err != nil {
			result.Err = model.NewAppError("SqlBackendProfileStore.Save", "store.sql_file_backend_profile.save.saving.app_error", nil,
				fmt.Sprintf("name=%s, %s", profile.Name, err.Error()), http.StatusInternalServerError)
		} else {
			result.Data = profile
		}
	})
}

func (s SqlFileBackendProfileStore) Get(id int, domain string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {

		query := `SELECT * FROM file_backend_profiles WHERE id = :Id AND domain = :Domain`

		profile := &model.FileBackendProfile{}

		if err := s.GetReplica().SelectOne(profile, query, map[string]interface{}{"Id": id, "Domain": domain}); err != nil {
			result.Err = model.NewAppError("SqlBackendProfileStore.Get", "store.sql_file_backend_profile.get.finding.app_error", nil,
				fmt.Sprintf("id=%d, domain=%s, %s", id, domain, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = profile
		}
	})
}

func (s SqlFileBackendProfileStore) GetById(id int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {

		query := `SELECT * FROM file_backend_profiles WHERE id = :Id`

		profile := &model.FileBackendProfile{}

		if err := s.GetReplica().SelectOne(profile, query, map[string]interface{}{"Id": id}); err != nil {
			result.Err = model.NewAppError("SqlBackendProfileStore.Get", "store.sql_file_backend_profile.get.finding.app_error", nil,
				fmt.Sprintf("id=%d, %s", id, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = profile
		}
	})
}

func (s SqlFileBackendProfileStore) List(domain string, offset, limit int) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		var profiles []*model.FileBackendProfile

		query := `SELECT * FROM file_backend_profiles 
			WHERE domain = :Domain  
			LIMIT :Limit OFFSET :Offset`

		if _, err := s.GetReplica().Select(&profiles, query, map[string]interface{}{"Domain": domain, "Offset": offset, "Limit": limit}); err != nil {
			result.Err = model.NewAppError("SqlBackendProfileStore.List", "store.sql_file_backend_profile.get_all.finding.app_error", nil, err.Error(), http.StatusInternalServerError)
		} else {
			result.Data = profiles
		}
	})
}

func (s SqlFileBackendProfileStore) Delete(id int, domain string) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		res, err := s.GetMaster().Exec("DELETE FROM file_backend_profiles WHERE id = :Id AND domain = :Domain", map[string]interface{}{"Id": id, "Domain": domain})
		if err != nil {
			result.Err = model.NewAppError("SqlBackendProfileStore.Delete", "store.sql_file_backend_profile.delete.app_error", nil,
				fmt.Sprintf("id=%d, err: %s", id, err.Error()), http.StatusInternalServerError)
			return
		}
		count, _ := res.RowsAffected()
		if count == 0 {
			result.Err = model.NewAppError("SqlBackendProfileStore.Delete", "store.sql_file_backend_profile.delete.not_found.app_error", map[string]interface{}{"Id": id},
				"", http.StatusNotFound)
		}
	})
}

func (s SqlFileBackendProfileStore) Update(profile *model.FileBackendProfile) store.StoreChannel {
	return store.Do(func(result *store.StoreResult) {
		if sqlResult, err := s.GetReplica().Exec(`
		UPDATE 
			file_backend_profiles
		SET name = :Name
			,type_id = :TypeId
			,"default" = :Default
			,expire_day = :ExpireDay
			,disabled = :Disabled
			,max_size_mb = :MaxSizeMb
			,properties = :Properties
		WHERE id = :Id AND domain = :Domain`,
			map[string]interface{}{
				"Id":     profile.Id,
				"Domain": profile.Domain,

				"Name":       profile.Name,
				"TypeId":     profile.TypeId,
				"Default":    profile.Default,
				"ExpireDay":  profile.ExpireDay,
				"Disabled":   profile.Disabled,
				"MaxSizeMb":  profile.MaxSizeMb,
				"Properties": model.StringInterfaceToJson(profile.Properties),
			}); err != nil {
			result.Err = model.NewAppError("SqlFileBackendProfileStore.Update", "store.sql_file_backend_profile.update.app_error",
				nil, err.Error(), http.StatusInternalServerError)
		} else {
			rows, err := sqlResult.RowsAffected()

			if err != nil {
				result.Err = model.NewAppError("SqlFileBackendProfileStore.Update", "store.sql_file_backend_profile.update.app_error",
					nil, err.Error(), http.StatusInternalServerError)
			} else {
				if rows == 1 {
					result.Data = true
				} else {
					result.Data = false
				}
			}
		}
	})
}
