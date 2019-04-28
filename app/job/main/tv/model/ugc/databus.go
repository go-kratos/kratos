package ugc

import (
	"time"

	arcmdl "go-common/app/service/main/archive/api"
)

// ArcMsg reprensents the archive Notify-T message structure
type ArcMsg struct {
	Action string       `json:"action"`
	Table  string       `json:"table"`
	Old    *ArchDatabus `json:"old"`
	New    *ArchDatabus `json:"new"`
}

// ArchDatabus model ( we pick the fields that we need )
type ArchDatabus struct {
	Aid       int64  `json:"aid"`
	Mid       int64  `json:"mid"`
	TypeID    int32  `json:"typeid"`
	Videos    int64  `json:"videos"`
	Duration  int    `json:"duration"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Content   string `json:"content"`
	Attribute int32  `json:"attribute"`
	Copyright int32  `json:"copyright"`
	State     int32  `json:"state"`
	Access    int    `json:"access"`
	PubTime   string `json:"pubtime"`
}

// VideoDiff reprensents the result of videos comparison
type VideoDiff struct {
	Aid     int64
	Equal   []int64 // totally equal
	New     []int64 // new added videos
	Updated []*arcmdl.Page
	Removed []int64
}

// DatabusVideo is the struct of message for the modification of ugc_Video
type DatabusVideo struct {
	New *MarkVideo `json:"new"`
	Old *MarkVideo `json:"old"`
}

// DatabusArc is the struct of message for the modification of ugc_archive
type DatabusArc struct {
	Old *MarkArc `json:"old"`
	New *MarkArc `json:"new"`
}

// MarkVideo contains the main fields that we want to pick up from databus message
type MarkVideo struct {
	Mark       int    `json:"mark"`
	Deleted    int    `json:"deleted"`
	CID        int64  `json:"cid"`
	AID        int64  `json:"aid"`
	EPTitle    string `json:"eptitle"`
	IndexOrder int    `json:"index_order"`
	Valid      int    `json:"valid"`
	Result     int    `json:"result"`
	Submit     int    `json:"submit"`
	Transcoded int    `json:"transcoded"`
	Retry      int64  `json:"retry"`
}

// MarkArc contains the main fields that we want to pick up from databus message
type MarkArc struct {
	ID         int    `json:"id"`
	AID        int64  `json:"aid"`
	MID        int    `json:"mid"`
	TypeID     int32  `json:"typeid"`
	Videos     int    `json:"videos"`
	Title      string `json:"title"`
	Cover      string `json:"cover"`
	Content    string `json:"content"`
	Duration   int    `json:"duration"`
	Copyright  int    `json:"copyright"`
	Pubtime    string `json:"pubtime"`
	Ctime      string `json:"ctime"`
	Mtime      string `json:"mtime"`
	State      int    `json:"state"`
	Manual     int    `json:"manual"`
	Valid      int    `json:"valid"`
	Submit     int    `json:"submit"`
	Retry      int    `json:"retry"`
	Result     int    `json:"result"`
	Deleted    int    `json:"deleted"`
	InjectTime string `json:"inject_time"`
	Reason     string `json:"reason"`
}

// IsPass returns whether the arc is able to play
func (a MarkArc) IsPass() bool {
	return a.Deleted == 0 && a.Valid == 1 && a.Result == 1
}

// CidReq reprensents the structure for reporting cid
type CidReq struct {
	CID int64 `json:"cid"`
}

// CidResp represents the structure of cid reporting API's response
type CidResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToReport distinguishes whether the CID need to be reported to video cloud
func (vm *MarkVideo) ToReport(criCID int64) bool {
	return vm.Deleted == 0 && vm.Mark == 0 && vm.CID > criCID
}

// ToAudit distinguishes whether the CID need to be reported to the license owner
func (vm *MarkVideo) ToAudit(criCID int64) bool {
	return vm.Submit == 1 && (vm.Transcoded == 1 || vm.CID <= criCID) && vm.Retry < time.Now().Unix() && vm.Deleted == 0
}

// CanPlay tells whether a video can play or not
func (vm *MarkVideo) CanPlay() bool {
	return vm.Result == 1 && vm.Deleted == 0 && vm.Valid == 1
}

// ToCMS transforms a databus video to CMS info
func (vm *MarkVideo) ToCMS() *VideoCMS {
	return &VideoCMS{
		CID:        int(vm.CID),
		Title:      vm.EPTitle,
		AID:        int(vm.AID),
		IndexOrder: vm.IndexOrder,
		Valid:      vm.Valid,
		Deleted:    vm.Deleted,
		Result:     vm.Result,
	}
}
