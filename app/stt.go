package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/webitel/storage/stt"

	"github.com/webitel/storage/stt/microsoft"

	"github.com/webitel/storage/model"
	"github.com/webitel/wlog"
)

func (app *App) GetSttProfileById(id int) (*model.CognitiveProfile, *model.AppError) {
	return app.Store.CognitiveProfile().GetById(int64(id))
}

func (app *App) GetSttProfile(id *int, syncTime *int64) (p *model.CognitiveProfile, appError *model.AppError) {
	var ok bool
	var cache interface{}

	if id == nil {
		return nil, model.NewAppError("GetSttProfile", "", nil, "", http.StatusInternalServerError)
	}

	cache, ok = app.sttProfilesCache.Get(*id)
	if ok {
		p = cache.(*model.CognitiveProfile)
		if syncTime != nil && p.GetSyncTime() == *syncTime {
			return
		}
	}

	if p == nil || syncTime == nil {
		p, appError = app.GetSttProfileById(*id)
		if appError != nil {
			return
		}
	}

	if appError != nil {
		return
	}

	switch p.Provider {
	case microsoft.ClientName:
		var err error

		if p.Instance, err = microsoft.NewClient(microsoft.ConfigFromJson(p.JsonProperties())); err != nil {
			// TODO
		}
	default:
		//todo error
	}

	app.sttProfilesCache.Add(*id, p)
	wlog.Info("[stt] Added to cache", wlog.String("name", p.Name))
	return p, nil
}

func (a *App) TranscriptFile(fileId int64, profileId int, locale string) (*model.FileTranscript, *model.AppError) {
	var fileUri string
	p, err := a.GetSttProfile(&profileId, nil) //todo
	if err != nil {
		return nil, err
	}

	if !p.Enabled {
		return nil, model.NewAppError("TranscriptFile", "app.stt.transcript.valid", nil, "Profile is disabled", http.StatusInternalServerError)
	}

	stt, ok := p.Instance.(stt.Stt)
	if !ok {
		return nil, model.NewAppError("TranscriptFile", "app.stt.transcript.valid", nil, "Bad client interface", http.StatusInternalServerError)
	}

	fileUri, err = a.GeneratePreSignetResourceSignature(model.AnyFileRouteName, "download", fileId, p.DomainId)
	if err != nil {
		return nil, err
	}

	//fmt.Println(fmt.Sprintf("https://dev.webitel.com%s", fileUri))

	ctx, _ := context.WithTimeout(context.Background(), time.Hour*2)

	if transcript, e := stt.Transcript(ctx, fmt.Sprintf("https://dev.webitel.com%s", fileUri), locale); e != nil {
		return nil, model.NewAppError("TranscriptFile", "app.stt.transcript.err", nil, e.Error(), http.StatusInternalServerError)
	} else {
		transcript.File = model.Lookup{
			Id: int(fileId),
		}
		transcript.Profile = model.Lookup{
			Id: profileId,
		}
		transcript.Locale = locale

		return a.Store.TranscriptFile().Store(&transcript)
	}
}
