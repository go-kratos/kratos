package mcnmodel

import (
	"encoding/json"
	"fmt"
	model2 "go-common/app/admin/main/mcn/model"
	"go-common/app/interface/main/mcn/model"
	"testing"
)

type replyCommon struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func createReply(data interface{}) *replyCommon {
	return &replyCommon{
		Message: "",
		Code:    0,
		Data:    data,
	}
}

func printReply(data interface{}) {
	var r = createReply(data)
	var result, _ = json.MarshalIndent(r, "", "    ")
	fmt.Printf(string(result) + "\n")
}

// var now = xtime.Time(time.Now().Unix())

func TestPrintMcnUpDataInfo(t *testing.T) {
	var r = McnGetUpListReply{
		Result:     []*McnUpDataInfo{{}},
		PageResult: model.PageResult{},
	}
	defer func() {}()
	printReply(r)
}

func TestPrintMcnGetUpPermitReply(t *testing.T) {
	var r = McnGetUpPermitReply{
		Old:          &model2.Permits{},
		New:          &model2.Permits{},
		ContractLink: "http://xxx.xx/adi.pdf",
	}
	defer func() {}()
	printReply(r)
}
