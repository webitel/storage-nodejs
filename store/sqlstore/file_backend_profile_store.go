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

	return us
}

func (self SqlFileBackendProfileStore) CreateIndexesIfNotExists() {

}

func (s SqlFileBackendProfileStore) Create(profile *model.FileBackendProfile) (*model.FileBackendProfile, *model.AppError) {
	err := s.GetMaster().SelectOne(&profile, `with p as (
    insert into file_backend_profiles (name, expire_day, priority, disabled, max_size_mb, properties, type_id,
                                           created_at, updated_at, created_by, updated_by,
                                           domain_id, description)
    values (:Name, :ExpireDay, :Priority, :Disabled, :MaxSize, :Properties, :TypeId, :CreatedAt, :UpdatedAt, :CreatedBy, :UpdatedBy,
            :DomainId, :Description)
    returning *
)
select p.id, call_center.cc_get_lookup(c.id, c.name) as created_by, p.created_at, call_center.cc_get_lookup(u.id, u.name) as updated_by,
       p.updated_at, p.name, p.description, p.expire_day, p.priority, p.disabled, p.max_size_mb, p.properties,
       call_center.cc_get_lookup(t.id, t.name) as type, p.data_size, p.data_count
from p
    inner join file_backend_profile_type t on t.id = p.type_id
    left join directory.wbt_user c on c.id = p.created_by
    left join directory.wbt_user u on u.id = p.updated_by`, map[string]interface{}{
		"Name":        profile.Name,
		"ExpireDay":   profile.ExpireDay,
		"Priority":    profile.Priority,
		"Disabled":    profile.Disabled,
		"MaxSize":     profile.MaxSizeMb,
		"Properties":  model.StringInterfaceToJson(profile.Properties),
		"TypeId":      profile.Type.Id,
		"CreatedAt":   profile.CreatedAt,
		"UpdatedAt":   profile.UpdatedAt,
		"CreatedBy":   profile.CreatedBy.Id,
		"UpdatedBy":   profile.UpdatedBy.Id,
		"DomainId":    profile.DomainId,
		"Description": profile.Description,
	})

	if err != nil {
		return nil, model.NewAppError("SqlFileBackendProfileStore.Create", "store.sql_file_backend_profile.get.create.app_error", nil,
			fmt.Sprintf("name=%v, %v", profile.Name, err.Error()), extractCodeFromErr(err))
	}

	return profile, nil
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
			result.Err = model.NewAppError("SqlBackendProfileStore.GetById", "store.sql_file_backend_profile.get.finding.app_error", nil,
				fmt.Sprintf("id=%d, %s", id, err.Error()), http.StatusInternalServerError)

			if err == sql.ErrNoRows {
				result.Err.StatusCode = http.StatusNotFound
			}
		} else {
			result.Data = profile
		}
	})
}

func (s SqlFileBackendProfileStore) GetAllPageByDomain(domain string, offset, limit int) store.StoreChannel {
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
		if sqlResult, err := s.GetMaster().Exec(`
		UPDATE 
			file_backend_profiles
		SET name = :Name
			,type_id = :TypeId
			,expire_day = :ExpireDay
			,disabled = :Disabled
			,max_size_mb = :MaxSizeMb
			,properties = :Properties
		WHERE id = :Id AND domain = :Domain`,
			map[string]interface{}{
				"Id": profile.Id,

				"Name":       profile.Name,
				"TypeId":     profile.Type.Id,
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
