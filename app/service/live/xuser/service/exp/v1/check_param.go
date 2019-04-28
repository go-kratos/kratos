package v1

import (
	"go-common/app/service/live/xuser/conf"
)

const (
	_addExpReqBizTypeFromPHPUser = 419
	// _addExpReqBizTypeCloseRoom   = 420
	// _addExpReqBizTypeUserOnline  = 421
)

const (
	_addUserExpType   = 1
	_addAnchorExpType = 2
)

var (
	_addExpReqBizMap = map[int64]string{
		// _addExpReqBizTypeCloseRoom:   "主播关播后经验结算",
		// _addExpReqBizTypeUserOnline:  "用户在线观看直播经验结算",
		_addExpReqBizTypeFromPHPUser: "用户经验结算",
	}
)

var (
	_addExpReqType = map[int64]string{
		_addUserExpType:   "增加用户经验",
		_addAnchorExpType: "增加主播经验",
	}
)

func (s *UserExpService) checkReqBiz(reqNum int64) (pass bool, rspBiz string) {
	rspBiz, pass = _addExpReqBizMap[reqNum]
	return
}

func (s *UserExpService) checkAddType(reqType int64) (pass bool, rspBiz string) {
	rspBiz, pass = _addExpReqType[reqType]
	return
}

func (s *UserExpService) checkAddNum(reqType int64) (pass bool, rspNum int64) {
	// 当前业务场景不支持减经验,若之后需要支持请修改本method以及相关逻辑,注意修改proto文件,使用了validate gt标志
	if reqType < 0 {
		pass = false
		rspNum = 0
	}
	pass = true
	rspNum = reqType
	return
}

func (s *UserExpService) getQueryStatus() (rspConf int) {
	if t := conf.Conf.Switch; t != nil {
		rspConf = t.QueryExp
	} else {
		rspConf = 0
	}
	return
}
