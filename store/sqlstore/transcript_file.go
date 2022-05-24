package sqlstore

import (
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
