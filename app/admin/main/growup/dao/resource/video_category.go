package resource

import (
	"context"

	"go-common/library/ecode"
	"go-common/library/log"
)

// videoInfo video category
type videoInfo struct {
	ID   int64  `json:"id"`
	PID  int64  `json:"pid"`
	Name string `json:"name"`
}

func videoCategory(c context.Context) (data map[int16]videoInfo, err error) {
	var res struct {
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    map[int16]videoInfo `json:"data"`
	}
	url := videoCategoryURL
	if err = client.Get(c, url, "", nil, &res); err != nil {
		log.Error("resource.videoCategory GET error(%v) | uri(%s)", err, url)
		return
	}
	if res.Code != 0 {
		log.Error("resource.videoCategory code != 0. res.Code(%d) | uri(%s) res(%v)", res.Code, url, res)
		err = ecode.GrowupGetTypeError
		return
	}
	data = res.Data
	return
}

// VideoCategoryIDToName .
func VideoCategoryIDToName(c context.Context) (res map[int64]string, err error) {
	data, err := videoCategory(c)
	if err != nil {
		return
	}
	res = make(map[int64]string, len(data))
	for _, v := range data {
		if v.PID == 0 {
			res[v.ID] = v.Name
		}
	}
	return
}

// VideoCategoryNameToID .
func VideoCategoryNameToID(c context.Context) (res map[string]int64, err error) {
	data, err := videoCategory(c)
	if err != nil {
		return
	}
	res = make(map[string]int64, len(data))
	for _, v := range data {
		if v.PID == 0 {
			res[v.Name] = v.ID
		}
	}
	return
}
