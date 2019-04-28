package model

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"go-common/library/log"
	"go-common/library/xstr"
)

// ParamHandler fn
type ParamHandler func(target *map[string]interface{}, args ...interface{})

// BatchParam str
type BatchParam struct {
	Params  map[string]interface{}
	Handler ParamHandler
}

// NewBatchParam func
func NewBatchParam(p map[string]interface{}, h ParamHandler) *BatchParam {
	if h == nil {
		h = BaseParamHandler
	}
	return &BatchParam{
		Params:  p,
		Handler: h,
	}
}

var (
	// BaseParamHandler fn
	BaseParamHandler = func(target *map[string]interface{}, args ...interface{}) {}
	// PushParamHandler fn
	PushParamHandler = func(target *map[string]interface{}, args ...interface{}) {
		var (
			ok   bool
			arc  *Archive
			fans []int64
		)
		if target == nil || (*target)["archive"] == nil {
			log.Warn("PushParamHandler target(%+v)/target[archive] nil, args(%+v)", target, args)
		} else if arc, ok = (*target)["archive"].(*Archive); !ok || arc == nil {
			log.Warn("PushParamHandler target[archive]=%+v parse failed/nil, args(%+v), target(%+v)", (*target)["archive"], args, target)
		}
		if arc == nil {
			arc = &Archive{}
		}

		if len(args) != 1 || args[0] == nil {
			log.Warn("PushParamHandler args(%+v) less than 1 or nil, target(%+v) archive(%+v)", args, target, arc)
		} else if fans, ok = args[0].([]int64); !ok {
			log.Warn("PushParamHandler args(%+v) parse failed, target(%+v) archive(%+v)", args, target, arc)
		}

		b := md5.Sum([]byte(fmt.Sprintf("%d%s", arc.ID, xstr.JoinInts(fans))))
		(*target)["uuid"] = fmt.Sprintf("%s%d", hex.EncodeToString(b[:]), time.Now().UnixNano())
	}
)
