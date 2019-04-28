package service

import (
	"go-common/app/job/main/passport-game-data/model"
	"go-common/library/log"
)

const (
	_statusOK      = 0
	_statusPending = 1
	_statusNo      = 2
)

func doCompare(cloud *model.AsoAccount, local *model.OriginAsoAccount, pending bool) int {
	if cloud == nil || local == nil {
		log.Info("either cloud or local aso account is nil, cloud %+v, local: %+v", cloud, local)
		return _statusNo
	}
	if cloud.Mtime.After(local.Mtime) {
		if model.Default(local).Equals(cloud) {
			return _statusOK
		}
		return _statusNo
	}
	if pending {
		return _statusPending
	}
	return _statusNo
}
