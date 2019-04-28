package cms

import "go-common/app/interface/main/tv/model"

const (
	// auth failure cases
	_pgcOffline    = "pgc_offline"
	_cmsInvalid    = "cms_invalid"
	_licenseReject = "license_reject"
	// auth correct values in DB
	_contentOnline = 0
	_snPass        = 1
	_epPass        = 3
	_cmsValid      = 1
)

// AuthMsg returns the error message in the config according to the error condition
func (d *Dao) authMsg(cond string) (msg string) {
	msgCfg := d.conf.Cfg.AuthMsg
	if cond == _licenseReject {
		return msgCfg.LicenseReject
	}
	if cond == _cmsInvalid {
		return msgCfg.CMSInvalid
	}
	if cond == _pgcOffline {
		return msgCfg.PGCOffline
	}
	return "err msg config load error"
}

// SnErrMsg returns the season auth result and the error message in case of auth failure
func (d *Dao) SnErrMsg(season *model.SnAuth) (bool, string) {
	if season.IsDeleted == _contentOnline && season.Check == _snPass && season.Valid == _cmsValid {
		return true, ""
	}
	if season.Check != _snPass {
		return false, d.authMsg(_licenseReject)
	}
	if season.IsDeleted != _contentOnline {
		return false, d.authMsg(_pgcOffline)
	}
	return false, d.authMsg(_cmsInvalid) // must be cms valid logically
}

// UgcErrMsg returns the arc auth result and the error message in case of auth failure
func (d *Dao) UgcErrMsg(deleted int, result int, valid int) (bool, string) {
	if deleted == _contentOnline && result == _snPass && valid == _cmsValid {
		return true, ""
	}
	if result != _snPass {
		return false, d.authMsg(_licenseReject)
	}
	if deleted != _contentOnline {
		return false, d.authMsg(_pgcOffline)
	}
	return false, d.authMsg(_cmsInvalid) // must be cms valid logically
}

// AuditingMsg returns the msg for the archive whose all videos are being audited
func (d *Dao) AuditingMsg() (bool, string) {
	return false, d.authMsg(_licenseReject) // must be cms valid logically
}

// EpErrMsg returns the ep auth result and the error message in case of auth failure
func (d *Dao) EpErrMsg(ep *model.EpAuth) (bool, string) {
	if ep.IsDeleted == _contentOnline && ep.State == _epPass && ep.Valid == _cmsValid {
		return true, ""
	}
	if ep.State != _epPass {
		return false, d.authMsg(_licenseReject)
	}
	if ep.IsDeleted != _contentOnline {
		return false, d.authMsg(_pgcOffline)
	}
	return false, d.authMsg(_cmsInvalid) // must be cms valid logically
}
