package service

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go-common/app/job/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestActionAddDM(t *testing.T) {
	id := int64(time.Now().UnixNano())
	dm := &model.DM{
		ID:       id,
		Type:     1,
		Oid:      1221,
		Mid:      4780461,
		Progress: 111,
		State:    0,
		Pool:     0,
		Ctime:    1533804859,
		Content: &model.Content{
			ID:       id,
			Mode:     4,
			IP:       123,
			FontSize: 25,
			Color:    12345,
			Msg:      "testtddddddddddddd",
			Ctime:    1533804859,
		},
	}
	Convey("", t, func() {
		data, err := json.Marshal(dm)
		So(err, ShouldBeNil)
		act := &model.Action{
			Action: model.ActAddDM,
			Data:   data,
		}
		err = svr.actionAct(context.TODO(), act)
		So(err, ShouldBeNil)
	})
}
