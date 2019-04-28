package dm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	model "go-common/app/job/main/activity/model/dm"
	"go-common/library/log"
)

// Broadcast dm broadcast.
func (d *Dao) Broadcast(c context.Context, dm *model.Broadcast) (err error) {
	var (
		b   []byte
		req *http.Request
		res struct {
			Code int64 `json:"code"`
		}
	)
	url := fmt.Sprintf("%s?cids=%d", d.broadcastURL, dm.RoomID)
	for i := 0; i < 50; i++ {
		if b, err = json.Marshal(dm); err != nil {
			log.Error("json.Marshal(%+v) error(%v)", dm, err)
			return
		}

		req, err = http.NewRequest("POST", url, bytes.NewReader(b))
		if err != nil {
			log.Error("NewRequest.Do(%d)(%s) error(%v)", dm.RoomID, url, err)
		}
		req.Header.Set("Content-type", "application/json")
		err = d.httpCli.Do(c, req, &res)
		if err == nil {
			break
		}
		log.Error("http.Do(%d)(%s) error(%v)", dm.RoomID, url, err)
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		log.Error("http.Do(%s) error(%v)", url, err)
		return
	}
	if res.Code != 0 {
		log.Error("http.Do(%s) res code(%d)", url, res.Code)
	}
	return
}
