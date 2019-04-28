package retry

import (
	"go-common/app/job/main/archive/model/result"
)

const (
	FailList          = "arc_job_fail_list"
	FailUpCache       = "up_cache"
	FailUpVideoCache  = "up_video_cache"
	FailDelVideoCache = "del_video_cache"
	FailDatabus       = "up_databus"
	FailResultAdd     = "result_add"
)

// Info retry data
type Info struct {
	Action string `json:"action"`
	Data   struct {
		Aid        int64                 `json:"aid"`
		State      int                   `json:"state"`
		DatabusMsg *result.ArchiveUpInfo `json:"dbus_msg"`
		Cids       []int64               `json:"cid"`
	} `json:"data"`
}
