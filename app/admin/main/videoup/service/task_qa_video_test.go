package service

import (
	"context"
	"encoding/json"
	"go-common/app/admin/main/videoup/model/archive"
	"testing"
)

func getQAVideo(s *Service) (qav *archive.QAVideo, err error) {
	vp := &archive.VideoParam{
		ID:       10005358,
		Aid:      10098493,
		Cid:      10109201,
		Status:   0,
		Encoding: 1,
		UID:      521,
		Oname:    "chenxi01",
		Title:    "测试添加qavideo",
		Filename: "j180311at3c5g5me6nt4h3zkbiltif35",
		RegionID: 23,
		Attrs: &archive.AttrParam{
			NoRank:      1,
			NoDynamic:   1,
			NoRecommend: 1,
			NoSearch:    1,
			OverseaLock: 1,
			PushBlog:    1,
		},
	}

	qav, err = s.fetchQAVideo(context.TODO(), vp)
	return
}

func Test_qavideoaddsingle(t *testing.T) {
	WithService(func(s *Service) {
		qav, err := getQAVideo(s)
		t.Logf("qavideo(%+v) err(%v)\r\n", qav, err)

		ctx := context.TODO()
		err = s.addQAVideo(ctx, qav)
		t.Logf("error(%v)", err)
	})()

}

func Test_qavideoadd(t *testing.T) {
	WithService(func(s *Service) {
		qav, err := getQAVideo(s)
		t.Logf("qavideo(%+v) error(%v)", qav, err)

		ctx := context.TODO()
		task, _ := json.Marshal(qav)

		err = s.arc.SendQAVideoAdd(ctx, task)
		if err != nil {
			t.Logf("s.arc.SendQAVideoAdd error(%v)", err)
		} else {
			t.Logf("s.arc.SendQAVideoAdd success")
		}
	})()

}
