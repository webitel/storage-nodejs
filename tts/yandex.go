package tts

import (
	"fmt"
	_ "github.com/yandex-cloud/go-genproto/yandex/cloud/ai/tts/v3"
	"io"
	"net/http"
	"net/url"
)

func Yandex(params TTSParams) (io.ReadCloser, *string, error) {
	api := fmt.Sprintf("https://tts.api.cloud.yandex.net/speech/v1/tts:synthesize?lang=%s", url.QueryEscape(params.Language))

	if params.Voice != "" {
		api += "&voice=" + params.Voice
	}

	if params.TextType == "ssml" {
		api += "&ssml=" + url.QueryEscape(params.Text)
	} else {
		api += "&text=" + url.QueryEscape(params.Text)
	}

	request, err := http.NewRequest("POST", api, nil)
	request.Header.Add("Authorization", fmt.Sprintf("Api-Key %s", params.Token))

	result, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, nil, err
	}

	return result.Body, nil, nil
}
