package stt

type Stt interface {
	Transcript(fileUri, locale string) (string, []byte, error)
}
