package model

import (
	"encoding/json"
	"time"
)

type TranscriptRange struct {
	StartSec float64 `json:"start_sec"`
	EndSec   float64 `json:"end_sec"`
}

type TranscriptWord struct {
	Word string `json:"word"`
	TranscriptRange
}

type TranscriptPhrase struct {
	TranscriptRange
	Channel int              `json:"channel"`
	Itn     string           `json:"itn"`
	Display string           `json:"display"`
	Lexical string           `json:"lexical"`
	Words   []TranscriptWord `json:"words"`
}

type TranscriptChannel struct {
	Channel int    `json:"channel"`
	Display string `json:"display"`
	Lexical string `json:"lexical"`
}

type FileTranscript struct {
	Id         int64               `json:"id" db:"id"`
	File       Lookup              `json:"file" db:"file"`
	Profile    Lookup              `json:"profile" db:"profile"`
	Transcript string              `json:"transcript" db:"transcript"`
	Log        json.RawMessage     `json:"log" db:"log"`
	CreatedAt  time.Time           `json:"created_at" db:"created_at"`
	Locale     string              `json:"locale" db:"locale"`
	Phrases    []TranscriptPhrase  `json:"phrases" db:"phrases"`
	Channels   []TranscriptChannel `json:"channels" db:"channels"`
}

func (f *FileTranscript) TidyTranscript() string {
	t := ""

	for k, v := range f.Channels {
		if k > 0 {
			t += " "
		}

		t += v.Display
	}

	return t
}

func (f *FileTranscript) JsonPhrases() []byte {
	if f.Phrases != nil {
		d, _ := json.Marshal(f.Phrases)
		return d
	}

	return nil
}

func (f *FileTranscript) JsonChannels() []byte {
	if f.Channels != nil {
		d, _ := json.Marshal(f.Channels)
		return d
	}

	return nil
}
