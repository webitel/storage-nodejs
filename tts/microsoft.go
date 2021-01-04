package tts

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	genderFemale = "Female"
)

func Microsoft(req TTSParams) (io.ReadCloser, *string, error) {
	var request *http.Request
	token, err := microsoftToken(req.Key, req.Region)
	if err != nil {
		return nil, nil, err
	}

	data := fmt.Sprintf(`<speak version='1.0' xml:lang='%s'>
	<voice xml:lang='%s' xml:gender='%s' name='Microsoft Server Speech Text to Speech Voice (%s)'>
	%s
	 </voice>
</speak>
`, req.Language, req.Language, req.Voice, microsoftLocalesNameMapping(req.Language, req.Voice), req.Text)

	request, err = http.NewRequest("POST", fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", req.Region), bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, nil, err
	}

	request.Header.Set("Content-Type", "application/ssml+xml")
	request.Header.Set("User-Agent", "WebitelACR")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	if strings.Index(req.Format, "wav") > -1 {
		request.Header.Set("X-Microsoft-OutputFormat", "riff-8khz-8bit-mono-mulaw")
	} else {
		request.Header.Set("X-Microsoft-OutputFormat", "audio-16khz-32kbitrate-mono-mp3")
	}

	result, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, nil, err
	}

	contentType := result.Header.Get("Content-Type")

	return result.Body, &contentType, nil
}

func microsoftToken(key, region string) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s.api.cognitive.microsoft.com/sts/v1.0/issueToken", region), nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Context-Type", "text/plain")
	req.Header.Set("Ocp-Apim-Subscription-Key", key)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var data []byte
	data, err = ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func isFemale(gender string) bool {
	return genderFemale == gender
}

func microsoftLocalesNameMapping(locale, gender string) string {
	switch locale {
	case "ar-EG":
		return "ar-EG, Hoda"

	case "ar-SA":
		return "ar-SA, Naayf"

	case "ca-ES":
		return "ca-ES, HerenaRUS"

	case "cs-CZ":
		return "cs-CZ, Vit"

	case "da-DK":
		return "da-DK, HelleRUS"

	case "de-AT":
		return "de-AT, Michael"

	case "de-CH":
		return "de-CH, Karsten"

	case "de-DE":
		if isFemale(gender) {
			return "de-DE, Hedda"
		} else {
			return "de-DE, Stefan, Apollo"
		}

	case "el-GR":
		return "el-GR, Stefanos"

	case "en-AU":
		return "en-AU, Catherine"

	case "id-ID":
		return "id-ID, Andika"

	case "en-CA":
		return "en-CA, Linda"

	case "en-GB":
		if isFemale(gender) {
			return "en-GB, Susan, Apollo"
		} else {
			return "en-GB, George, Apollo"
		}

	case "en-IE":
		return "en-IE, Shaun"

	case "en-IN":
		if isFemale(gender) {
			return "en-IN, Heera, Apollo"
		} else {
			return "en-IN, Ravi, Apollo"
		}

	case "en-US":
		if isFemale(gender) {
			return "en-US, ZiraRUS"
		} else {
			return "en-US, BenjaminRUS"
		}

	case "es-ES":
		if isFemale(gender) {
			return "es-ES, Laura, Apollo"
		} else {

			return "es-ES, Pablo, Apollo"
		}

	case "es-MX":
		if isFemale(gender) {
			return "es-MX, HildaRUS"
		} else {
			return "es-MX, Raul, Apollo"
		}

	case "fi-FI":
		return "fi-FI, HeidiRUS"

	case "fr-CA":
		if isFemale(gender) {
			return "fr-CA, Caroline"
		} else {
			return "fr-CH, Guillaume"
		}

	case "fr-FR":
		if isFemale(gender) {
			return "fr-FR, Julie, Apollo"
		} else {

			return "fr-FR, Paul, Apollo"
		}

	case "he-IL":
		return "he-IL, Asaf"

	case "hi-IN":
		if isFemale(gender) {
			return "hi-IN, Kalpana, Apollo"
		} else {
			return "hi-IN, Hemant"
		}

	case "hu-HU":
		return "hu-HU, Szabolcs"

	case "it-IT":
		return "it-IT, Cosimo, Apollo"

	case "ja-JP":
		if isFemale(gender) {
			return "ja-JP, Ayumi, Apollo"
		} else {
			return "ja-JP, Ichiro, Apollo"
		}

	case "ko-KR":
		return "ko-KR, HeamiRUS"

	case "nb-NO":
		return "nb-NO, HuldaRUS"

	case "nl-NL":
		return "nl-NL, HannaRUS"

	case "pl-PL":
		return "pl-PL, PaulinaRUS"

	case "pt-BR":
		if isFemale(gender) {
			return "pt-BR, HeloisaRUS"
		} else {
			return "pt-BR, Daniel, Apollo"
		}

	case "pt-PT":
		return "pt-PT, HeliaRUS"

	case "ro-RO":
		return "ro-RO, Andrei"

	case "ru-RU":
		if isFemale(gender) {
			return "ru-RU, Irina, Apollo"
		} else {

			return "ru-RU, Pavel, Apollo"
		}

	case "sk-SK":
		return "sk-SK, Filip"

	case "sv-SE":
		return "sv-SE, HedvigRUS"

	case "th-TH":
		return "th-TH, Pattara"

	case "tr-TR":
		return "tr-TR, SedaRUS"

	case "zh-CN":
		if isFemale(gender) {
			return "zh-CN, Yaoyao, Apollo"
		} else {
			return "zh-CN, Kangkang, Apollo"
		}

	case "zh-HK":
		if isFemale(gender) {
			return "zh-HK, Tracy, Apollo"
		} else {
			return "zh-HK, Danny, Apollo"
		}

	case "zh-TW":
		if isFemale(gender) {
			return "zh-TW, Yating, Apollo"
		} else {
			return "zh-TW, Zhiwei, Apollo"
		}

	case "vi-VN":
		return "vi-VN, An"

	default:
		return ""

	}
}
