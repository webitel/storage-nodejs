package model

import "time"

type FileTranscript struct {
	Id         int64     `json:"id" db:"id"`
	File       Lookup    `json:"file" db:"file"`
	Transcript string    `json:"transcript" db:"transcript"`
	Log        []byte    `json:"log" db:"log"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
