package microsoft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/webitel/storage/model"

	"github.com/pkg/errors"
)

const (
	ClientName = "Microsoft"
)

var (
	ErrNotFound = errors.New("Not found files")
)

type client struct {
	key    string
	region string
	http   http.Client
}

type Config struct {
	Key    string `json:"key"`
	Region string `json:"region"`
}

type transcriptRequest struct {
	ContentUrls []string `json:"contentUrls"`
	Properties  struct {
		WordLevelTimestampsEnabled bool `json:"wordLevelTimestampsEnabled"`
	} `json:"properties"`
	Locale      string `json:"locale"`
	DisplayName string `json:"displayName"`
}

type File struct {
	Name  string `json:"name"`
	Links struct {
		ContentUrl string `json:"contentUrl"`
	} `json:"links"`
}

type Files struct {
	Values []*File `json:"values"`
}

type Transcript struct {
	RecognizedPhrases []struct {
		Offset   string `json:"offset"`
		Duration string `json:"duration"`
		Channel  int    `json:"channel"`
		NBest    []struct {
			Words []struct {
				Word     string `json:"word"`
				Offset   string `json:"offset"`
				Duration string `json:"duration"`
			} `json:"words"`
			Itn     string `json:"itn"`
			Display string `json:"display"`
			Lexical string `json:"lexical"`
		} `json:"nBest"`
	} `json:"recognizedPhrases"`
	CombinedRecognizedPhrases []struct {
		Lexical string `json:"lexical"`
		Channel int    `json:"channel"`
		Display string `json:"display"`
	} `json:"combinedRecognizedPhrases"`
}

type Task struct {
	Id     string `json:"id"`
	Self   string `json:"self"`
	Status string `json:"status"`
	Links  struct {
		Files string `json:"files"`
	} `json:"links"`
	Properties struct {
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	} `json:"properties"`
}

func NewClient(config Config) (*client, error) {
	return &client{
		key:    config.Key,
		region: config.Region,
		http:   http.Client{},
	}, nil
}

func (c *client) Transcript(ctx context.Context, fileUri, locale string) (model.FileTranscript, error) {
	var data []byte

	task, err := c.TranscriptJob(fileUri, locale)
	if err != nil {
		return model.FileTranscript{}, err
	}

	if _, err = c.WaitFoSuccess(ctx, task); err != nil {
		return model.FileTranscript{}, err
	}

	if task.Properties.Error != nil {
		return model.FileTranscript{}, errors.New(task.Properties.Error.Message)
	}

	data, err = c.LoadTranscript(task)
	if err != nil {
		return model.FileTranscript{}, err
	}

	ph, cs := getTranscript(data)

	res := model.FileTranscript{
		Log:       data,
		Phrases:   ph,
		Channels:  cs,
		CreatedAt: time.Now(),
	}

	return res, nil
}

func (t Task) Finished() bool {
	return t.Status == "Succeeded" || t.Status == "Failed"
}

func (c *client) TranscriptJob(fileUrl string, locale string) (*Task, error) {

	tr := &transcriptRequest{
		ContentUrls: []string{fileUrl},
		Properties: struct {
			WordLevelTimestampsEnabled bool `json:"wordLevelTimestampsEnabled"`
		}{
			true,
		},
		Locale:      locale,
		DisplayName: fmt.Sprintf("Transcription using default model for %s", locale),
	}

	data, _ := json.Marshal(tr)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s.api.cognitive.microsoft.com/speechtotext/v3.0/transcriptions", c.region), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.key)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var t Task
	err = json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *client) Finished(t *Task) (bool, error) {
	// todo or error ?
	var data []byte
	if t.Finished() {
		return true, nil
	}

	req, err := http.NewRequest("GET", t.Self, nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.key)

	res, err := c.http.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(data, &t)
	if err != nil {
		return false, err
	}

	return t.Finished(), nil
}

func (c *client) LoadTranscript(t *Task) ([]byte, error) {
	files, err := c.GetFiles(t)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNotFound
	}

	file := files[len(files)-1]
	var data []byte

	req, err := http.NewRequest("GET", file.Links.ContentUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.key)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if data, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *client) GetFiles(t *Task) ([]*File, error) {
	var data []byte

	req, err := http.NewRequest("GET", t.Links.Files, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ocp-Apim-Subscription-Key", c.key)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if data, err = ioutil.ReadAll(res.Body); err != nil {
		return nil, err
	}

	var result Files
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result.Values, nil
}

func (c *client) WaitFoSuccess(ctx context.Context, t *Task) (ok bool, err error) {

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 2):
			if ok, err = c.Finished(t); ok || err != nil {
				return
			}
		}
	}
}

func getTranscript(data []byte) ([]model.TranscriptPhrase, []model.TranscriptChannel) {
	var n Transcript
	if err := json.Unmarshal(data, &n); err != nil {
		//TODO error
		return nil, nil
	}

	res := make([]model.TranscriptPhrase, 0, len(n.RecognizedPhrases))

	for _, v := range n.RecognizedPhrases {
		if len(v.NBest) < 1 || len(v.NBest[0].Words) < 1 {
			continue
		}

		words := make([]model.TranscriptWord, 0, len(v.NBest[0].Words))

		for _, w := range v.NBest[0].Words {
			words = append(words, model.TranscriptWord{
				Word: w.Word,
				TranscriptRange: model.TranscriptRange{
					StartSec: (ParseDuration(w.Offset)).Seconds(),
					EndSec:   (ParseDuration(w.Offset) + ParseDuration(w.Duration)).Seconds(),
				},
			})
		}

		res = append(res, model.TranscriptPhrase{
			TranscriptRange: model.TranscriptRange{
				StartSec: (ParseDuration(v.Offset)).Seconds(),
				EndSec:   (ParseDuration(v.Offset) + ParseDuration(v.Duration)).Seconds(),
			},
			Channel: v.Channel,
			Itn:     v.NBest[0].Itn,
			Display: v.NBest[0].Display,
			Lexical: v.NBest[0].Lexical,
			Words:   words,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].StartSec < res[j].StartSec
	})

	cs := make([]model.TranscriptChannel, 0, len(n.CombinedRecognizedPhrases))
	for _, v := range n.CombinedRecognizedPhrases {
		cs = append(cs, model.TranscriptChannel{
			Channel: v.Channel,
			Display: v.Display,
			Lexical: v.Lexical,
		})
	}

	return res, cs
}

var durationRegex = regexp.MustCompile(`P([\d\.]+Y)?([\d\.]+M)?([\d\.]+D)?T?([\d\.]+H)?([\d\.]+M)?([\d\.]+?S)?`)

// ParseDuration converts a ISO8601 duration into a time.Duration
func ParseDuration(str string) time.Duration {
	matches := durationRegex.FindStringSubmatch(str)

	years := parseDurationPart(matches[1], time.Hour*24*365)
	months := parseDurationPart(matches[2], time.Hour*24*30)
	days := parseDurationPart(matches[3], time.Hour*24)
	hours := parseDurationPart(matches[4], time.Hour)
	minutes := parseDurationPart(matches[5], time.Second*60)
	seconds := parseDurationPart(matches[6], time.Second)

	return time.Duration(years + months + days + hours + minutes + seconds)
}

func parseDurationPart(value string, unit time.Duration) time.Duration {
	if len(value) != 0 {
		if parsed, err := strconv.ParseFloat(value[:len(value)-1], 64); err == nil {
			return time.Duration(float64(unit) * parsed)
		}
	}
	return 0
}
