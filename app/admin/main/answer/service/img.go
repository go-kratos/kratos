package service

import (
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image/jpeg"
	"net/http"
	"strconv"
	"time"

	"go-common/library/log"
	"go-common/library/text/translate/chinese"

	"github.com/golang/freetype"
)

var (
	allOrders  [][]int
	x          = []int{0, 1, 2, 3}
	bname      = "member"
	accessKey  = "3d34b1ea1dbbb0ca"
	secretKey  = "d4caa344f3b115e302033b05dd0aa4"
	_template  = "%s\n%s\n%s\n%d\n"
	fileFormat = "v3_%s_A-%s_B-%s_C-%s_D-%s_%s"
	bfsURL     = "http://bfs.bilibili.co/bfs/%s/%s"
	platform   = map[string]bool{"H5": true, "PC": false}
	language   = []string{"zh-CN", "zh-TW"}
)

// generate 全排列组合算法   4*3*2*1 种情况
func (s *Service) generate(c context.Context, a []int, begin int, end int) {
	if begin == end {
		allOrders = append(allOrders, []int{a[0], a[1], a[2], a[3]})
		return
	}
	for i := begin; i <= end; i++ {
		a[begin], a[i] = a[i], a[begin]
		s.generate(c, a, begin+1, end)
		a[begin], a[i] = a[i], a[begin]
	}
	return
}

// GenerateImage .
func (s *Service) GenerateImage(c context.Context) {
	// 把有效id放redis
	ids, _ := s.dao.IDsByState(c)
	s.CreateBFSImg(c, ids)
}

// CreateBFSImg .
func (s *Service) CreateBFSImg(c context.Context, ids []int64) {
	for _, qid := range ids {
		log.Error("qid_%d", qid)
		que, err := s.dao.QueByID(c, qid)
		log.Error("que:%+v", que)
		if err != nil {
			log.Error("get question(%d), err(%v)", qid, err)
		}
		for _, langv := range language {
			if langv == "zh-TW" {
				que.Question = chinese.Convert(c, que.Question)
				que.Ans[0] = chinese.Convert(c, que.Ans[0])
				que.Ans[1] = chinese.Convert(c, que.Ans[1])
				que.Ans[2] = chinese.Convert(c, que.Ans[2])
				que.Ans[3] = chinese.Convert(c, que.Ans[3])
			}
			for plat, platv := range platform {
				for _, order := range allOrders {
					log.Error("allOrders(len:%d) %s_%s_%s_%s_%s_%s, %+v", len(allOrders), strconv.FormatInt(qid, 10), que.Ans[order[0]], que.Ans[order[1]], que.Ans[order[2]], que.Ans[order[3]], plat, que)

					quec := s.dao.QueConf(platv)
					imgh := s.dao.Height(quec, que.Question, len(que.Ans))
					r := s.dao.Board(imgh)
					imgc := s.dao.Context(r, s.c.Answer.FontFilePath)
					pt := freetype.Pt(0, int(quec.Fontsize))
					s.dao.DrawQue(imgc, que.Question, quec, &pt)

					as := [4]string{que.Ans[order[0]], que.Ans[order[1]], que.Ans[order[2]], que.Ans[order[3]]}
					s.dao.DrawAns(imgc, quec, as, &pt)

					buf := new(bytes.Buffer)
					jpeg.Encode(buf, r, nil)
					bufReader := bufio.NewReader(buf)

					ts := time.Now().Unix()
					tm := time.Unix(ts, 0)
					m := md5.New()
					m.Write([]byte(fmt.Sprintf(fileFormat, strconv.FormatInt(qid, 10), que.Ans[order[0]], que.Ans[order[1]], que.Ans[order[2]], que.Ans[order[3]], plat)))
					fname := hex.EncodeToString(m.Sum(nil)) + ".jpg"
					if s.c.Answer.Debug {
						fname = fmt.Sprintf("debug_%s", fname)
					}
					//imgName = append(imgName, fname)
					client := &http.Client{}
					content := fmt.Sprintf(_template, "PUT", bname, fname, ts)
					mac := hmac.New(sha1.New, []byte(secretKey))
					mac.Write([]byte(content))
					sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
					req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf(bfsURL, bname, fname), bufReader)
					req.Host = "bfs.bilibili.co"
					req.Header.Add("Date", tm.Format("2006-01-02 03:04:05"))
					req.Header.Add("Authorization", accessKey+":"+sign+":"+strconv.FormatInt(ts, 10))
					req.Header.Set("Content-Type", "image/jpeg")
					resp, _ := client.Do(req)
					//defer resp.Body.Close()
					if err != nil {
						log.Error("qid %s bfs upload failed error(%s)", strconv.FormatInt(qid, 10), err)
						//return
					}
					log.Info("upload img(%d), %s, %+v", qid, fname, resp)
					time.Sleep(time.Millisecond * 50)
				}
			}
		}
	}
}
