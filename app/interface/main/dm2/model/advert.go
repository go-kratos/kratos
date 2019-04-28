package model

import (
	"encoding/json"
	"fmt"
)

// resource id defined by advert
const (
	adRscIDIphone      = 2630
	adRscIDAndrod      = 2631
	adRscIDIphoneIcon  = 2642
	adRscIDAndroidIcon = 2643
)

// Resource get resource by mobi_app.
func Resource(mobiApp string) (rsc string) {
	if mobiApp == "iphone" || mobiApp == "ipad" || mobiApp == "iphone_i" {
		rsc = fmt.Sprintf("%d,%d", adRscIDIphone, adRscIDIphoneIcon)
	} else {
		rsc = fmt.Sprintf("%d,%d", adRscIDAndrod, adRscIDAndroidIcon)
	}
	return
}

// ADReq advert request params
type ADReq struct {
	Aid      int64  `json:"aid"`
	Oid      int64  `json:"oid"`
	Mid      int64  `json:"mid"`
	Build    int64  `json:"build"`
	Buvid    string `json:"buvid"`
	ClientIP string `json:"ip"`
	MobiApp  string `json:"mobi_app"`
	ADExtra  string `json:"ad_extra"`
}

// ADResp advert response
type ADResp struct {
	Icon *ADInfo   `json:"icon,omitempty"`
	ADs  []*ADInfo `json:"ads_info,omitempty"`
}

// AD advert struct
type AD struct {
	RequestID string                      `json:"request_id,omitempty"`
	ADsInfo   map[int64]map[int64]*ADInfo `json:"ads_info,omitempty"` // resource_id --> source_id --> adinfo
}

// ADInfo advert info.
type ADInfo struct {
	// filed response from advert api
	Index     int             `json:"index,omitempty"`
	IsAd      bool            `json:"is_ad,omitempty"`
	CmMark    int             `json:"cm_mark,omitempty"`
	CardIndex int             `json:"card_index,omitempty"`
	ADInfo    json.RawMessage `json:"ad_info,omitempty"`
	// filed used in app
	RequestID  string `json:"request_id,omitempty"`
	ResourceID int64  `json:"resource_id,omitempty"`
	SourceID   int64  `json:"source_id,omitempty"`
	ClientIP   string `json:"client_ip,omitempty"`
	IsADLoc    bool   `json:"is_ad_loc,omitempty"`
}

// Convert convert AD to ADResp.
func (a *AD) Convert(clientIP string) (res *ADResp) {
	res = new(ADResp)
	for rscID, adInfoMap := range a.ADsInfo {
		for srcID, adInfo := range adInfoMap {
			v := new(ADInfo)
			v.RequestID = a.RequestID
			v.ResourceID = rscID
			v.SourceID = srcID
			v.ClientIP = clientIP
			v.IsADLoc = true // 该字段服务端代码写死为true
			if adInfo != nil {
				v.Index = adInfo.Index
				v.IsAd = adInfo.IsAd
				v.CmMark = adInfo.CmMark
				v.CardIndex = adInfo.CardIndex
			}
			if len(adInfo.ADInfo) > 0 {
				v.ADInfo = adInfo.ADInfo
			}
			if v.ResourceID == adRscIDIphoneIcon || v.ResourceID == adRscIDAndroidIcon { // icon resouce id
				res.Icon = v
				continue
			}
			res.ADs = append(res.ADs, v)
		}
	}
	return
}
