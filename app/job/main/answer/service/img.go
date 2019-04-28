package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"strconv"

	"go-common/app/job/main/answer/model"
	"go-common/library/log"
	"go-common/library/text/translate/chinese"

	"github.com/golang/freetype"
)

var (
	platform   = map[string]bool{"H5": true, "PC": false}
	fileFormat = "%s_A-%s_B-%s_%s"

	language = []string{"zh-CN", "zh-TW"}
)

// createBFSImg create bfs img.
func (s *Service) createBFSImg(c context.Context, que *model.LabourQs) (err error) {
	log.Info("createBFSImg(%v)", que)
	as := [2]string{model.ExtraAnsA, model.ExtraAnsB}
	for _, langv := range language {
		if langv == "zh-TW" {
			que.Question = chinese.Convert(c, que.Question)
			as = [2]string{chinese.Convert(c, model.ExtraAnsA), chinese.Convert(c, model.ExtraAnsB)}
		}
		for ps, pb := range platform {
			quec := s.dao.QueConf(pb)
			imgh := s.dao.Height(quec, que.Question, 2)
			r := s.dao.Board(imgh)
			imgc := s.dao.Context(r, s.c.Properties.FontFilePath)
			pt := freetype.Pt(0, int(quec.Fontsize))
			s.dao.DrawQue(imgc, que.Question, quec, &pt)
			s.dao.DrawAns(imgc, quec, as, &pt)

			buf := new(bytes.Buffer)
			jpeg.Encode(buf, r, nil)
			bufReader := bufio.NewReader(buf)

			m := md5.New()
			m.Write([]byte(fmt.Sprintf(fileFormat, strconv.FormatInt(que.ID, 10), as[0], as[1], ps)))
			fname := hex.EncodeToString(m.Sum(nil)) + ".jpg"

			location, err := s.dao.Upload(c, "image/jpeg", fname, bufReader)
			if err != nil {
				log.Error("question (%v) bfs upload failed error(%s)", que, err)
				continue
			}
			log.Info("upload img success.fname:%s,location:%s,lang:%s", fname, location, langv)
		}
	}
	return
}
