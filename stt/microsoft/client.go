package microsoft

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	ClientName = "microsoft"
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
	CombinedRecognizedPhrases []struct {
		Lexical string `json:"lexical"`
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

func (c *client) Transcript(fileUri, locale string) (string, []byte, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)
	var data []byte

	task, err := c.TranscriptJob(fileUri, locale)
	if err != nil {
		return "", nil, err
	}

	if _, err = c.WaitFoSuccess(ctx, task); err != nil {
		return "", nil, err
	}

	if task.Properties.Error != nil {
		return "", nil, errors.New(task.Properties.Error.Message)
	}

	data, err = c.LoadTranscript(task)
	if err != nil {
		return "", nil, err
	}

	return getText(data), data, nil
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

func getText(data []byte) string {
	var n Transcript
	json.Unmarshal(data, &n)
	if len(n.CombinedRecognizedPhrases) > 0 {
		return n.CombinedRecognizedPhrases[0].Lexical
	}

	return ""
}
