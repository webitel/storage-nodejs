package app

import (
	"fmt"
	"net/http"

	"github.com/webitel/storage/stt"

	"github.com/webitel/storage/stt/microsoft"

	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
)

func (app *App) GetSttProfileById(id int) (*model.SttProfile, *model.AppError) {
	panic("TYODO")
	//return app.Store.CognitiveProfile().GetById(id)
}

func (app *App) GetSttProfile(id *int, syncTime *int64) (p *model.SttProfile, appError *model.AppError) {
	var ok bool
	var cache interface{}

	if id == nil {
		return nil, model.NewAppError("GetSttProfile", "", nil, "", http.StatusInternalServerError)
	}

	cache, ok = app.sttProfilesCache.Get(*id)
	if ok {
		p = cache.(*model.SttProfile)
		if syncTime != nil && p.GetSyncTime() == *syncTime {
			return
		}
	}

	if p == nil {
		p, appError = app.GetSttProfileById(*id)
		if appError != nil {
			return
		}
	}

	if appError != nil {
		return
	}

	switch p.Type {
	case microsoft.ClientName:
		var err error

		if p.Instance, err = microsoft.NewClient(microsoft.ConfigFromJson(p.Config)); err != nil {
			// TODO
		}
	default:
		//todo error
	}

	app.sttProfilesCache.Add(*id, p)
	wlog.Info("[STT] Added to cache", wlog.String("name", p.Name))
	return p, nil
}

func (a *App) TranscriptFile(fileId int64, profileId int, locale string) *model.AppError {
	var fileUri string
	p, err := a.GetSttProfile(&profileId, nil) //todo
	if err != nil {
		return err
	}
	stt, ok := p.Instance.(stt.Stt)
	if !ok {
		return model.NewAppError("TranscriptFile", "app.stt.transcript.valid", nil, "Bad client interface", http.StatusInternalServerError)
	}

	fileUri, err = a.GeneratePreSignetResourceSignature(model.AnyFileRouteName, "download", fileId, p.DomainId)
	if err != nil {
		return err
	}

	if transcript, logs, e := stt.Transcript(fmt.Sprintf("https://dev.webitel.com%s", fileUri), locale); e != nil {
		return model.NewAppError("TranscriptFile", "app.stt.transcript.err", nil, e.Error(), http.StatusInternalServerError)
	} else {
		fmt.Println(transcript, "\n", string(logs))
	}

	return nil
}
