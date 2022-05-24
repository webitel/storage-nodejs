package synchronizer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/webitel/storage/app"
	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
)

type SttJob struct {
	file model.SyncJob
	app  *app.App
}

type params struct {
	ProfileId int    `json:"profile_id"`
	Locale    string `json:"locale"`
}

func (s *SttJob) Execute() {
	var p params
	json.Unmarshal(s.file.Config, &p)

	n := time.Now()

	if p.Locale != "" && p.ProfileId > 0 {
		s.transcript(p.Locale, p.ProfileId)
	} else {
		wlog.Error(fmt.Sprintf("[stt] file %d, error: profile_id & locale is requered", s.file.FileId))
	}

	err := s.app.Store.SyncFile().Remove(s.file.Id)
	if err != nil {
		wlog.Error(fmt.Sprintf("[stt] file %d, error: %s", s.file.FileId, err.Error()))
	}
	wlog.Debug(fmt.Sprintf("[stt] job_id: %d, file_id: %d stop, time %v", s.file.Id, s.file.FileId, time.Since(n)))

}

func (s *SttJob) transcript(local string, profileId int) {
	wlog.Debug(fmt.Sprintf("[stt] job_id: %d, file_id: %d start transcript to %s", s.file.Id, s.file.FileId, local))
	t, err := s.app.TranscriptFile(s.file.FileId, profileId, local)
	if err != nil {
		wlog.Error(err.Error())
	} else {
		wlog.Debug(fmt.Sprintf("[stt] file %d, transcript: %s", s.file.FileId, t.Transcript))
	}
}
