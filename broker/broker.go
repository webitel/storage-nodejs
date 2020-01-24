package broker

import (
	"encoding/json"
	"github.com/webitel/storage/model"
)

type Broker interface {
	Close()
	RPC(request *Request) (*Response, *model.AppError)
}

type Handler func(req *Request) []byte

type Request struct {
	Api  string                 `json:"api"`
	Args map[string]interface{} `json:"args"`
}

type Response struct {
	Data interface{} `json:"data"`
}

func NewRequestFromBytes(data []byte) (err error, req Request) {
	err = json.Unmarshal(data, &req)
	return
}
