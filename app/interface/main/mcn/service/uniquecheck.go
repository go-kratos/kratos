package service

import (
	"sync"

	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/time"
)

var (
	lastLoadMcnUniqueTime time.Time
)

//UniqueCheck check unique
type UniqueCheck struct {
	// all the values is mcn id
	PhoneMap            map[string]int64
	IDCardMap           map[string]int64
	CompanyNameMap      map[string]int64
	CompanyLicenseIDMap map[string]int64
	lock                sync.Mutex
}

//NewUniqueCheck new checker
func NewUniqueCheck() *UniqueCheck {
	return &UniqueCheck{
		PhoneMap:            make(map[string]int64),
		IDCardMap:           make(map[string]int64),
		CompanyNameMap:      make(map[string]int64),
		CompanyLicenseIDMap: make(map[string]int64),
	}
}

//CheckIsUniqe check is unique
func (u *UniqueCheck) CheckIsUniqe(req *mcnmodel.McnApplyReq) (err error) {
	if req == nil {
		return
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	if v, ok := u.PhoneMap[req.ContactPhone]; ok {
		if req.McnMid != v {
			err = ecode.MCNUpBindUpDuplicatePhone
			return
		}
	}

	if v, ok := u.IDCardMap[req.ContactIdcard]; ok {
		if req.McnMid != v {
			err = ecode.MCNUpBindUpDuplicateIDCard
			return
		}
	}

	if v, ok := u.CompanyNameMap[req.CompanyName]; ok {
		if req.McnMid != v {
			err = ecode.MCNUpBindUpDuplicateCompanyName
			return
		}
	}

	if v, ok := u.CompanyLicenseIDMap[req.CompanyLicenseID]; ok {
		if req.McnMid != v {
			err = ecode.MCNUpBindUpDuplicateCompanyLicenseID
			return
		}
	}
	return
}

//AddItem add item from db
func (u *UniqueCheck) AddItem(sign *mcnmodel.McnSign) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.PhoneMap[sign.ContactPhone] = sign.McnMid
	u.IDCardMap[sign.ContactIdcard] = sign.McnMid
	u.CompanyNameMap[sign.CompanyName] = sign.McnMid
	u.CompanyLicenseIDMap[sign.CompanyLicenseID] = sign.McnMid
}

func (s *Service) loadMcnUniqueCache() {
	var list []*mcnmodel.McnSign
	var err = s.mcndao.GetMcnDB().
		Select("mcn_mid, company_name, company_license_id, contact_idcard, contact_phone, mtime").
		Where("mtime>?", lastLoadMcnUniqueTime).
		Find(&list).Error
	if err != nil {
		log.Warn("cannot get unique, err=%s", err)
		return
	}

	for _, v := range list {
		s.uniqueChecker.AddItem(v)
		if lastLoadMcnUniqueTime < v.Mtime {
			lastLoadMcnUniqueTime = v.Mtime
		}
	}
}
