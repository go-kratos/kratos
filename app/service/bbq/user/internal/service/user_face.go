package service

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"go-common/app/service/bbq/user/internal/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/json-iterator/go"
)

func (s *Service) subBfsUserFace() {
	mch := s.userFaceSub.Messages()
	for {
		msg, ok := <-mch
		if !ok {
			continue
		}
		s.bfsUserFaceConsumer(msg)
	}
}

func (s *Service) bfsUserFaceConsumer(msg *databus.Message) {
	if msg == nil {
		return
	}
	defer msg.Commit()

	inst := new(model.UserFaceBFS)
	err := jsoniter.Unmarshal(msg.Value, inst)
	if err != nil {
		log.Error("subBfsUserFace error: %+v", err)
		return
	}

	if inst.Bucket != "bbq" || !inst.IsYellow {
		return
	}
	// TODO: remove debug info
	log.Info("subBfsUserFace instance: %+v", inst)

	if inst.URL == "" {
		return
	}

	s1 := strings.TrimPrefix(inst.FileName, "video-image/userface/")
	if s1 == inst.FileName {
		log.Error("subBfsUserFace mid parse failed filename: %s", inst.FileName)
		return
	}

	midInfo := strings.Split(s1, "_")
	if len(midInfo) != 2 {
		log.Error("subBfsUserFace mid parse failed filename: %s", inst.FileName)
		return
	}

	// 推送到审核后台
	params := url.Values{
		"mid":         {midInfo[0]},
		"origin_face": {"https://i0.hdslb.com/bfs/bbq/video-image/userface/1558868601542006937.png"},
		"face":        {inst.URL},
	}
	resp, err := http.DefaultClient.PostForm("http://bbq-mng.bilibili.co/bbq/cms/user/setface", params)
	if err != nil {
		log.Error("subBfsUserFace post user_face to cms error: %+v", err)
		return
	}
	response, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	log.Info("subBfsUserFace result: %s, error: %+v", string(response), err)
}
