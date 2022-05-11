package main

// #cgo LDFLAGS: -lavcodec -lavformat -lavutil -lavdevice -lswresample
//#include "audio.c"
import "C"
import (
	"reflect"
	"unsafe"
)

type audioBuffer struct {
	channels C.int
	len      C.int
	data     *C.float
}

type transcoding struct {
	buffer *audioBuffer
}

func (e *transcoding) Decode(path string) {
	var rate C.int = 16000

	buf := &audioBuffer{}
	C.transcode(unsafe.Pointer(buf), C.CString(path), rate)

	e.buffer = buf
}

func (e *transcoding) Close() {
	if e.buffer != nil {
		C.clean_audio_buffer(unsafe.Pointer(e.buffer))
		e.buffer = nil
	}
}

func (e *transcoding) Channels() [][]float32 {
	res := make([][]float32, e.buffer.channels, e.buffer.channels)
	channels := (int)(e.buffer.channels)
	bufLen := (int)(e.buffer.len)
	length := bufLen / channels

	var list []C.float
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&list))
	sliceHeader.Cap = bufLen
	sliceHeader.Len = bufLen
	sliceHeader.Data = uintptr(unsafe.Pointer(e.buffer.data))

	for i := 0; i < channels; i++ {
		res[i] = make([]float32, length, length)
	}

	for k, v := range list {
		res[k%channels][(k-(k%channels))/channels] = (float32)(v)
	}

	list = nil

	return res
}

var Transcoding transcoding

//export ChunkCallback
func ChunkCallback(ctx *C.void, samples *C.float, len int, channels int) {
	var list []C.float
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&list))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(samples))

	//r := make([]C.float, 0, len/2)
	//for i := 1; i < len; i = i + 2 {
	//	r = append(r, list[i])
	//}

	//binary.Write(gfile, binary.LittleEndian, list)
}

func main() {
	C.initTranscoding()

	//var t transcoding
	//t.Decode("https://cloud.webitel.ua/api/storage/recordings/444552/stream?access_token=gdpiujx3oiftxg5j8ewrheb13h")
	//c := t.Channels()
	//buf := new(bytes.Buffer)
	//binary.Write(buf, binary.LittleEndian, c[1])
	//gfile.Write(buf.Bytes())
	//
	//gfile.Close()
	//runtime.GC()
}
