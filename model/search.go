package model

import (
	"encoding/json"
	"io"
)

type SearchEngineRequest struct {
	Index           string
	Type            string
	Columns         []string                  `json:"columnsDate"` //TODO rename to columns
	Includes        []string                  `json:"includes"`
	Excludes        []string                  `json:"excludes"`
	ScrollKeepAlive *string                   `json:"scroll"`
	Size            int                       `json:"limit"`
	Page            int                       `json:"pageNumber"` //TODO rename to page
	Query           string                    `json:"query"`
	Filter          SearchEngineRequestFilter `json:"filter"`
	Sort            SearchEngineRequestSort   `json:"sort"`
}

type SearchEngineScroll struct {
	ScrollKeepAlive string `json:"scroll"`
	ScrollId        string `json:"scrollId"` //TODO rename to scroll_id
}

type SearchEngineRequestSort struct {
	Data map[string]interface{}
}

type SearchEngineRequestFilter struct {
	Data []interface{}
}

type BaseHitsResponseSearchEngine struct {
	Hits  []*SearchEngineHitsResponse `json:"hits"`
	Total int64                       `json:"total"`
}

type SearchEngineHitsResponse struct {
	Id     string                 `json:"_id"`
	Index  string                 `json:"_index"`
	Fields map[string]interface{} `json:"fields"`
	Source interface{}            `json:"_source"`
}

type SearchEngineResponse struct {
	Hits     BaseHitsResponseSearchEngine `json:"hits"`
	Timeout  bool                         `json:"timeout"`
	Shards   map[string]int               `json:"_shards"`
	ScrollId *string                      `json:"_scroll_id,omitempty"`
}

type SortSearchEngine struct {
	Data []byte
}

func (s *SearchEngineResponse) ToJson() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (info SortSearchEngine) Source() (interface{}, error) {
	source := make(map[string]interface{})
	err := json.Unmarshal(info.Data, &source)
	if err != nil {
		return nil, err
	}

	return source, nil
}

func (f SearchEngineRequestSort) Source() (interface{}, error) {
	return f.Data, nil
}

func (f *SearchEngineRequestSort) IsEmpty() bool {
	for range f.Data {
		return false
	}
	return true
}

func (f *SearchEngineRequestSort) UnmarshalJSON(b []byte) (err error) {
	if len(b) == 2 {
		//TODO set default
	}
	return json.Unmarshal(b, &f.Data)
}

func (f SearchEngineRequestFilter) Source() (interface{}, error) {
	return f.Data, nil
}

func (f *SearchEngineRequestFilter) UnmarshalJSON(b []byte) (err error) {
	return json.Unmarshal(b, &f.Data)
}

func (f *SearchEngineRequestFilter) AddFilter(q map[string]interface{}) {
	f.Data = append(f.Data, q)
}

func SearchEngineRequestFromJson(data io.Reader) *SearchEngineRequest {
	var req SearchEngineRequest
	if err := json.NewDecoder(data).Decode(&req); err == nil {
		//TODO old bug
		if req.Page > 0 {
			req.Page--
		}
		return &req
	} else {
		return nil
	}
}
func SearchEngineScrollFromJson(data io.Reader) *SearchEngineScroll {
	var req SearchEngineScroll
	if err := json.NewDecoder(data).Decode(&req); err == nil {
		return &req
	} else {
		return nil
	}
}
