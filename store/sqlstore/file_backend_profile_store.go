package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/webitel/engine/auth_manager"
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

func (s SqlFileBackendProfileStore) CheckAccess(domainId, id int64, groups []int, access auth_manager.PermissionAccess) (bool, *model.AppError) {

	res, err := s.GetReplica().SelectNullInt(`select 1
		where exists(
          select 1
          from file_backend_profiles_acl a
          where a.dc = :DomainId
            and a.object = :Id
            and a.subject = any (:Groups::int[])
            and a.access & :Access = :Access
        )`, map[string]interface{}{"DomainId": domainId, "Id": id, "Groups": pq.Array(groups), "Access": access.Value()})

	if err != nil {
		return false, nil
	}

	return res.Valid && res.Int64 == 1, nil
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

func (s SqlFileBackendProfileStore) GetAllPage(domainId int64, offset, limit int) ([]*model.FileBackendProfile, *model.AppError) {
	var profiles []*model.FileBackendProfile
	_, err := s.GetMaster().Select(&profiles, `select p.id, call_center.cc_get_lookup(c.id, c.name) as created_by, p.created_at, call_center.cc_get_lookup(u.id, u.name) as updated_by,
       p.updated_at, p.name, p.description, p.expire_day, p.priority, p.disabled, p.max_size_mb, p.properties,
       call_center.cc_get_lookup(t.id, t.name) as type, p.data_size, p.data_count
from file_backend_profiles p
    inner join file_backend_profile_type t on t.id = p.type_id
    left join directory.wbt_user c on c.id = p.created_by
    left join directory.wbt_user u on u.id = p.updated_by
    where p.domain_id = :DomainId
    order by p.priority
	limit :Limit
	offset :Offset`, map[string]interface{}{
		"DomainId": domainId,
		"Limit":    limit,
		"Offset":   offset,
	})
	if err != nil {
		return nil, model.NewAppError("SqlBackendProfileStore.GetAllPage", "store.sql_file_backend_profile.get_all.finding.app_error",
			nil, err.Error(), extractCodeFromErr(err))
	}

	return profiles, nil
}

func (s SqlFileBackendProfileStore) GetAllPageByGroups(domainId int64, groups []int, offset, limit int) ([]*model.FileBackendProfile, *model.AppError) {
	var profiles []*model.FileBackendProfile
	_, err := s.GetMaster().Select(&profiles, `select p.id, call_center.cc_get_lookup(c.id, c.name) as created_by, p.created_at, call_center.cc_get_lookup(u.id, u.name) as updated_by,
       p.updated_at, p.name, p.description, p.expire_day, p.priority, p.disabled, p.max_size_mb, p.properties,
       call_center.cc_get_lookup(t.id, t.name) as type, p.data_size, p.data_count
from file_backend_profiles p
    inner join file_backend_profile_type t on t.id = p.type_id
    left join directory.wbt_user c on c.id = p.created_by
    left join directory.wbt_user u on u.id = p.updated_by
    where p.domain_id = :DomainId and (
		exists(select 1
		  from file_backend_profiles_acl a
		  where a.dc = p.domain_id and a.object = p.id and a.subject = any(:Groups::int[]) and a.access&:Access = :Access)
	  )
    order by p.priority
	limit :Limit
	offset :Offset`, map[string]interface{}{
		"DomainId": domainId,
		"Limit":    limit,
		"Offset":   offset,
		"Groups":   pq.Array(groups),
		"Access":   auth_manager.PERMISSION_ACCESS_READ.Value(),
	})
	if err != nil {
		return nil, model.NewAppError("SqlBackendProfileStore.GetAllPageByGroups", "store.sql_file_backend_profile.get_all.finding.app_error",
			nil, err.Error(), extractCodeFromErr(err))
	}

	return profiles, nil
}

//FIXME
func (s SqlFileBackendProfileStore) Get(id, domainId int64) (*model.FileBackendProfile, *model.AppError) {
	var profile *model.FileBackendProfile
	err := s.GetMaster().SelectOne(&profile, `select p.id, call_center.cc_get_lookup(c.id, c.name) as created_by, p.created_at, call_center.cc_get_lookup(u.id, u.name) as updated_by,
       p.updated_at, p.name, p.description, p.expire_day, p.priority, p.disabled, p.max_size_mb, p.properties,
       call_center.cc_get_lookup(t.id, t.name) as type, p.data_size, p.data_count, p.domain_id
from file_backend_profiles p
    inner join file_backend_profile_type t on t.id = p.type_id
    left join directory.wbt_user c on c.id = p.created_by
    left join directory.wbt_user u on u.id = p.updated_by
    where p.id = :Id and p.domain_id = :DomainId`, map[string]interface{}{
		"Id":       id,
		"DomainId": domainId,
	})

	if err != nil {
		return nil, model.NewAppError("SqlBackendProfileStore.Get", "store.sql_file_backend_profile.get.app_error", nil,
			fmt.Sprintf("id=%d, domain=%d, %s", id, domainId, err.Error()), extractCodeFromErr(err))
	}

	return profile, nil
}

func (s SqlFileBackendProfileStore) Update(profile *model.FileBackendProfile) (*model.FileBackendProfile, *model.AppError) {
	err := s.GetMaster().SelectOne(&profile, `with p as (
    update file_backend_profiles
    set name = :Name,
        expire_day = :ExpireDay,
        priority = :Priority,
        disabled = :Disabled,
        max_size_mb = :MaxSize,
        properties = :Properties,
        description = :Description,
        updated_at = :UpdatedAt,
        updated_by = :UpdatedBy
    where id = :Id and domain_id = :DomainId
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
		"UpdatedAt":   profile.UpdatedAt,
		"UpdatedBy":   profile.UpdatedBy.Id,
		"DomainId":    profile.DomainId,
		"Description": profile.Description,
		"Id":          profile.Id,
	})

	if err != nil {
		return nil, model.NewAppError("SqlBackendProfileStore.Update", "store.sql_file_backend_profile.update.app_error", nil,
			fmt.Sprintf("id=%d, domain=%d, %s", profile.Id, profile.DomainId, err.Error()), extractCodeFromErr(err))
	}

	return profile, nil
}

func (s SqlFileBackendProfileStore) Delete(domainId, id int64) *model.AppError {
	if _, err := s.GetMaster().Exec(`delete from file_backend_profiles p where id = :Id and domain_id = :DomainId`,
		map[string]interface{}{"Id": id, "DomainId": domainId}); err != nil {
		return model.NewAppError("SqlFileBackendProfileStore.Delete", "store.sql_file_backend_profile.delete.app_error", nil,
			fmt.Sprintf("Id=%v, %s", id, err.Error()), extractCodeFromErr(err))
	}
	return nil
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
