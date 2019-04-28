package model

import (
	xtime "go-common/library/time"

	avmdl "go-common/app/interface/main/app-view/model"
)

// ads const
const (
	// ads plat
	VdoAdsPC      = int8(0)
	VdoAdsIPhone  = int8(1)
	VdoAdsAndroid = int8(2)
	VdoAdsIPad    = int8(3)
	// ads type
	VdoAdsTypeBangumi = int8(0)
	VdoAdsTypeNologin = int8(1)
	VdoAdsTypeNothing = int8(2)
	VdoAdsTypeOther   = int8(3)
	// ads target
	VdoAdsTargetArchive = int8(1)
	VdoAdsTargetBangumi = int8(2)
	VdoAdsTargetType    = int8(3)
)

// VideoAD is Ads of videos.
type VideoAD struct {
	Name          string     `json:"name"`
	ContractID    string     `json:"contract_id"`
	Aid           int64      `json:"aid"`
	SeasonID      int        `json:"season_id"`
	TypeID        int16      `json:"type _id"`
	AdCid         int64      `json:"ad_cid"`
	AdStrategy    int        `json:"ad_strategy"`
	AdURL         string     `json:"ad_url"`
	AdOrder       int        `json:"ad_order"`
	Skipable      int8       `json:"skipable"`
	Note          string     `json:"note"`
	AgencyName    string     `json:"agency_name"`
	AgencyCountry int        `json:"agency_country"`
	AgencyArea    int        `json:"agency_area"`
	Price         float32    `json:"price"`
	Verified      int        `json:"verified"`
	State         int        `json:"state"`
	FrontAid      int64      `json:"front_aid"`
	Target        int8       `json:"target"`
	Platform      int8       `json:"platform"`
	Type          int8       `json:"type"`
	UserSet       int8       `json:"user_set"`
	PlayCount     int64      `json:"play_count"`
	MTime         xtime.Time `json:"mtime"`
	Aids          string     `json:"-"`
}

// Paster struct
type Paster struct {
	AID       int64  `json:"aid"`
	CID       int64  `json:"cid"`
	Duration  int64  `json:"duration"`
	Type      int8   `json:"type"`
	AllowJump int8   `json:"allow_jump"`
	URL       string `json:"url"`
}

// PasterPlat exchange plat to video_ads
func PasterPlat(plat int8) int8 {
	switch plat {
	case PlatWEB:
		return VdoAdsPC
	case avmdl.PlatIPad, avmdl.PlatIpadHD, avmdl.PlatIPadI: // 2、9、6 -> 3
		return VdoAdsIPad
	case avmdl.PlatIPhone, avmdl.PlatIPhoneI: // 1、5 -> 1
		return VdoAdsIPhone
	case avmdl.PlatAndroid, avmdl.PlatAndroidG, avmdl.PlatAndroidI, avmdl.PlatAndroidTV, avmdl.PlatWPhone: // 0、4、8、7、3 -> 2
		return VdoAdsAndroid
	}
	return VdoAdsIPhone // 1
}
