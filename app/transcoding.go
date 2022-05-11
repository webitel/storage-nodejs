package app

type Transcoding interface {
	Decode(string)
	Close()
	Channels() [][]float32
}
