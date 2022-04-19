package main

// #cgo LDFLAGS: -lavcodec -lavformat -lavutil -lavdevice -lswresample
//#include "audio.c"
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"unsafe"
)

type audioBuffer struct {
	len  int
	data *C.int16_t
}

type transcoding struct {
	data []int16
	buf  *audioBuffer
}

var (
	gbuf     = new(bytes.Buffer)
	gfile, _ = os.Create("gfile.bin")
)

//export FuncInMain
func FuncInMain(_ *C.void, frameCount int, buffer *C.int16_t) {
	//t := (*transcoding)(unsafe.Pointer(&buffer))
	gbuf.Reset()
	var list []int16
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&list))
	sliceHeader.Cap = frameCount
	sliceHeader.Len = frameCount
	sliceHeader.Data = uintptr(unsafe.Pointer(buffer))

	err := binary.Write(gbuf, binary.LittleEndian, list)
	if err != nil {
		panic(err.Error())
	}

	binary.Write(gfile, binary.BigEndian, gbuf.Bytes())
}

func (e *transcoding) carray2slice() {
}

func (e *transcoding) Decode(path string) {
	b := &audioBuffer{}

	C.test_buffer(unsafe.Pointer(b), C.CString(path))
	e.buf = b
}

func (e *transcoding) Close() {
	if e.buf != nil {
		C.clean_audio_buffer(unsafe.Pointer(e.buf))
		e.buf = nil
	}
}

func (e *transcoding) Wave() []int16 {
	var list []int16
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&list))
	sliceHeader.Cap = e.buf.len
	sliceHeader.Len = e.buf.len
	sliceHeader.Data = uintptr(unsafe.Pointer(e.buf.data))
	return list
}

func (e *transcoding) Bytes() (error, []byte) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, e.Wave())
	if err != nil {
		return err, nil
	}
	return nil, buf.Bytes()
}

var Transcoding transcoding

func test() {
	d := &transcoding{}
	d.Decode("https://cloud.webitel.ua/api/storage/recordings/444552/download?access_token=gdpiujx3oiftxg5j8ewrheb13h")
	defer d.Close()

	arr := d.Wave()
	f, _ := os.Create("file.bin")
	binary.Write(f, binary.BigEndian, arr)
	f.Close()
	fmt.Println(len(arr))
	gfile.Close()
}

func main() {
	//var i int64
	//C.dd((*C.int64_t)(unsafe.Pointer(&i)))
	//fmt.Println(i)

	test()
	return

	b := &audioBuffer{}

	C.test_buffer(unsafe.Pointer(b), C.CString("https://dev.webitel.com/api/storage/recordings/58492/stream?access_token=7osjt9gqnjg85y1cj6q4ifqusw"))
	res := carray2slice(b.data, b.len)

	f, _ := os.Create("file.bin")
	binary.Write(f, binary.BigEndian, res)
	f.Close()
	C.clean_audio_buffer(unsafe.Pointer(b))
}

func carray2slice(array *C.int16_t, len int) []C.int16_t {
	var list []C.int16_t
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&list))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return list
}
