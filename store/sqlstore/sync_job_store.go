package sqlstore

import (
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
	"net/http"
)

type SqlSyncFileStore struct {
	SqlStore
}

func NewSqlSyncFileStore(sqlStore SqlStore) store.SyncFileStore {
	us := &SqlSyncFileStore{sqlStore}
	return us
}

func (s SqlSyncFileStore) FetchRemoveJobs(limit int) ([]*model.SyncJob, *model.AppError) {
	var res []*model.SyncJob
	_, err := s.GetMaster().Select(&res, `update storage.remove_file_jobs u
set state = 1
from (
    select j.id, j.file_id, f.domain_id, f.properties, f.profile_id, p.updated_at as profile_updated_at, f.name, f.size, f.mime_type, f.instance
    from storage.remove_file_jobs j
        inner join storage.files f on f.id = j.file_id
        left join storage.file_backend_profiles p on p.id = f.profile_id
    where j.state = 0
    order by j.created_at asc
    limit :Limit
    for update OF j skip locked
) j
where u.id = j.id and u.state = 0
returning j.*`, map[string]interface{}{
		"Limit": limit,
	})

	if err != nil {
		return nil, model.NewAppError("SqlSyncFileStore.FetchRemoveJobs", "store.sql_sync_file_job.save.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return res, nil
}

func (s SqlSyncFileStore) SetRemoveJobs(localExpDay int) *model.AppError {
	_, err := s.GetMaster().Exec(`insert into storage.remove_file_jobs (file_id)
select f.id
from (
    select (extract(epoch from now() - (p.expire_day || ' day')::interval) * 1000)::int8 as exp, p.id::int4
    from storage.file_backend_profiles p
    where p.expire_day > 0
    union all
    select (extract(epoch from now() - (:LocalExpire || ' day')::interval) * 1000)::int8 as exp, 0::int4
    where :LocalExpire > 0
) p
inner join lateral (
   select  f.id as id
    from storage.files f
    where coalesce(f.profile_id, 0) = p.id
        and f.created_at < p.exp
        and not exists(select 1 from storage.remove_file_jobs j where j.file_id = f.id)
    order by f.created_at asc
) f on true
union all
select id
from (
    select f.id
    from storage.files f
    where f.removed
        and not exists(select 1 from storage.remove_file_jobs j where j.file_id = f.id)
    order by f.created_at
) t`, map[string]interface{}{
		"LocalExpire": localExpDay,
	})

	if err != nil {
		return model.NewAppError("SqlSyncFileStore.SetRemoveJobs", "store.sql_sync_file_job.set_removed.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s SqlSyncFileStore) Clean(jobId int64) *model.AppError {
	_, err := s.GetMaster().Exec(`with del as (
    delete
    from storage.remove_file_jobs rj
    where rj.id = :Id
    returning rj.file_id
)
delete
from storage.files f
where f.id = (select del.file_id from del )`, map[string]interface{}{
		"Id": jobId,
	})

	if err != nil {
		return model.NewAppError("SqlSyncFileStore.Remove", "store.sql_sync_file_job.clean.app_error", nil, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
