package sqlstore

import (
	"github.com/lib/pq"
	"github.com/webitel/storage/model"
	"github.com/webitel/storage/store"
)

type SqlTranscriptFileStore struct {
	SqlStore
}

func NewSqlTranscriptFileStore(sqlStore SqlStore) store.TranscriptFileStore {
	us := &SqlTranscriptFileStore{sqlStore}
	return us
}

func (s SqlTranscriptFileStore) GetByFileId(fileId int64, profileId int64) (*model.FileTranscript, *model.AppError) {
	var t model.FileTranscript
	err := s.GetReplica().SelectOne(&t, `select t.id,
       storage.get_lookup(f.id, f.name) as file,
       storage.get_lookup(p.id, p.name) as profile,
       t.transcript,
       t.log,
       t.locale,
       t.created_at
from storage.file_transcript t
    left join storage.files f on f.id = t.file_id
    left join storage.cognitive_profile_services p on p.id = t.profile_id
where t.id = :Id and t.profile_id = :ProfileId`, map[string]interface{}{
		"Id":        fileId,
		"ProfileId": profileId,
	})

	if err != nil {
		return nil, model.NewAppError("SqlTranscriptFileStore.GetByFileId", "store.sql_stt_file.get.app_error", nil, err.Error(), extractCodeFromErr(err))
	}

	return &t, nil
}

func (s SqlTranscriptFileStore) Store(t *model.FileTranscript) (*model.FileTranscript, *model.AppError) {
	err := s.GetMaster().SelectOne(&t, `with t as (
    insert into storage.file_transcript (file_id, transcript, log, profile_id, locale, phrases, channels)
    values (:FileId, :Transcript, :Log, :ProfileId, :Locale, :Phrases, :Channels)
    returning *
)
select t.id,
       storage.get_lookup(f.id, f.name) as file,
       storage.get_lookup(p.id, p.name) as profile,
       t.transcript,
       t.log,
       t.created_at,
	   t.phrases,
	   t.channels	
from t
    left join storage.files f on f.id = t.file_id
    left join storage.cognitive_profile_services p on p.id = t.profile_id`, map[string]interface{}{
		"FileId":     t.File.Id,
		"Transcript": t.TidyTranscript(),
		"Log":        t.Log,
		"Locale":     t.Locale,
		"ProfileId":  t.Profile.Id,
		"Phrases":    t.JsonPhrases(),
		"Channels":   t.JsonChannels(),
	})

	if err != nil {
		return nil, model.NewAppError("SqlTranscriptFileStore.Store", "store.sql_stt_file.store.app_error", nil, err.Error(), extractCodeFromErr(err))
	}

	return t, nil
}

func (s SqlTranscriptFileStore) CreateJobs(domainId int64, fileIds []int64, params model.TranscriptOptions) ([]*model.FileTranscriptJob, *model.AppError) {
	var jobs []*model.FileTranscriptJob
	_, err := s.GetMaster().Select(&jobs, `insert into storage.file_jobs (state, file_id, action, config)
select 0 as state,
       fid,
       p.service,
       json_build_object('locale', :Locale::varchar,
           'profile_id', p.id,
           'profile_sync_time', (extract(epoch from p.updated_at) * 1000 )::int8) as config
from storage.cognitive_profile_services p,
     unnest((:FileIds)::int8[]) fid
where p.domain_id = :DomainId::int8
    and p.id = :Id::int4
    and p.enabled
    and p.service = :Service::varchar
returning storage.file_jobs.id,
    storage.file_jobs.file_id,
    (extract(epoch from storage.file_jobs.created_at) * 1000)::int8 as created_at,
    storage.file_jobs.state`, map[string]interface{}{
		"DomainId": domainId,
		"FileIds":  pq.Array(fileIds),
		"Id":       params.ProfileId,
		"Locale":   params.Locale,
		"Service":  model.SyncJobSTT,
	})

	if err != nil {
		return nil, model.NewAppError("SqlTranscriptFileStore.CreateJobs", "store.sql_stt_file.create.jobs.app_error", nil, err.Error(), extractCodeFromErr(err))
	}

	return jobs, nil
}
