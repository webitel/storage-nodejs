package model

import "encoding/json"

type TtsProfile struct {
	Enabled    bool            `json:"enabled" db:"enabled"`
	Provider   string          `json:"provider" db:"provider"`
	Properties json.RawMessage `json:"properties" db:"properties"`
}
