package liverpc

import "encoding/json"

const (
	_magic     = 2233
	_headerLen = 60

	_cmdReqType = byte('0')
)

type protoHeader struct {
	magic     uint32
	timestamp uint32
	checkSum  uint32
	version   uint32
	reserved  uint32
	seq       uint32
	length    uint32
	cmd       []byte
}

type protoReq struct {
	Header protoHeader
	Body   []byte
}

type protoResp struct {
	Header protoHeader
	Body   []byte
}

// Args .
type Args struct {
	Header *Header     `json:"header"`
	Body   interface{} `json:"body"`
	HTTP   interface{} `json:"http"`
}

// Reply .
type Reply struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}
