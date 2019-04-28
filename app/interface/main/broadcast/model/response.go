package model

import (
	"encoding/json"
	"go-common/library/ecode"
)

// Response .
type Response struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// Message .
func Message(raw map[string]interface{}, e error) (bs []byte) {
	res := &Response{
		Code:    ecode.Cause(e).Code(),
		Message: ecode.Cause(e).Message(),
		Data:    raw,
	}
	bs, _ = json.Marshal(res)
	return
}
