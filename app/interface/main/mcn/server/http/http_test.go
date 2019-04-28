package http

import (
	"testing"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/net/http/blademaster/binding"
)

func TestValidater(t *testing.T) {
	var req = model.MCNSignEntryReq{
		MCNMID:    1,
		BeginDate: "2018-01-01",
		EndDate:   "2018-01-01",
	}
	var err = binding.Validator.ValidateStruct(&req)
	if err == nil {
		t.FailNow()
	} else {
		t.Logf("err=%s", err)
	}
}
