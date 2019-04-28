package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/coupon/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	xtime "go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

// ActivitySalaryCoupon activity salary coupon.
func (s *Service) ActivitySalaryCoupon(c context.Context, req *model.ArgBatchSalaryCoupon) (err error) {
	var (
		ok bool
		r  *model.CouponBatchInfo
	)
	if r, err = s.dao.BatchInfo(c, req.BranchToken); err != nil {
		return
	}
	if r == nil {
		return ecode.CouPonBatchNotExistErr
	}
	if r.State != model.BatchStateNormal {
		return ecode.CouPonHadBlockErr
	}
	if r.ExpireDay != -1 {
		return ecode.CouponTypeNotSupportErr
	}
	if ok = s.dao.AddGrantUniqueLock(c, req.BranchToken, _lockseconds); !ok {
		return ecode.CouponSalaryHadRunErr
	}
	// mids index
	s.runSalary(metadata.String(c, metadata.RemoteIP), r, req)
	return
}

// runSalary run batch salary.
func (s *Service) runSalary(ip string, b *model.CouponBatchInfo, req *model.ArgBatchSalaryCoupon) (err error) {
	go func() {
		var (
			err    error
			fails  []int64
			tmp    []int64
			mids   = []int64{}
			total  int64
			midmap map[int64][]int64
			c      = context.Background()
		)
		defer func() {
			if x := recover(); x != nil {
				x = errors.WithStack(x.(error))
				log.Error("runtime error caught: %+v", x)
			}
			//Write fail mids
			if len(fails) > 0 {
				s.OutFile(c, []byte(xstr.JoinInts(fails)), req.BranchToken+"_failmids.csv")
			}
		}()
		log.Info("batch salary analysis slary file[req(%v)]", req)
		// analysis slary file
		if midmap, total, err = s.AnalysisFile(c, req.FileURL); err != nil {
			return
		}
		if total != req.Count {
			log.Error("[batch salary count err req(%v),filetotal(%d)]", req, total)
			s.dao.DelGrantUniqueLock(c, b.BatchToken)
			return
		}
		log.Info("batch salary start[count(%d),token(%s)]", total, req.BranchToken)
		var crrent int
		for _, allmids := range midmap {
			for i, v := range allmids {
				mids = append(mids, v)
				if len(mids)%req.SliceSize != 0 && i != len(allmids)-1 {
					continue
				}
				crrent = crrent + len(mids)
				log.Info("batch salary ing[mid(%d),i(%d),token(%s),crrent(%d)]", v, i, req.BranchToken, crrent)
				if tmp, err = s.batchSalary(c, mids, ip, b); err != nil && len(tmp) > 0 {
					fails = append(fails, tmp...)
					log.Error("batch salary err[mids(%v),i(%d),token(%s)] err(%v)", mids, i, req.BranchToken, err)
				}
				mids = []int64{}
			}
		}
		log.Info("batch salary end[count(%d),token(%s),faillen(%d)]", total, req.BranchToken, len(fails))
	}()
	return
}

// batchSalary batch salary.
func (s *Service) batchSalary(c context.Context, mids []int64, ip string, b *model.CouponBatchInfo) (falls []int64, err error) {
	var (
		tx  *sql.Tx
		aff int64
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	defer func() {
		if err != nil {
			falls = mids
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("batch salary rollback: %+v", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			falls = mids
			log.Error("batch salary commit: %+v", err)
			tx.Rollback()
			return
		}
		//del cache
		s.cache.Do(c, func(c context.Context) {
			for _, v := range mids {
				if err = s.dao.DelCouponAllowancesKey(context.Background(), v, model.NotUsed); err != nil {
					log.Error("batch salary cache del(%d) err: %+v", v, err)
				}
			}
		})
		//send msg
		if s.c.Prop.SalaryMsgOpen {
			s.sendMsg(mids, ip, 1, true, _msgMeng)
		}
		time.Sleep(time.Duration(s.c.Prop.SalarySleepTime))
	}()
	cps := make([]*model.CouponAllowanceInfo, len(mids))
	for i, v := range mids {
		cps[i] = &model.CouponAllowanceInfo{
			CouponToken: s.tokeni(i),
			Mid:         v,
			State:       model.NotUsed,
			StartTime:   b.StartTime,
			ExpireTime:  b.ExpireTime,
			Origin:      model.AllowanceSystemAdmin,
			CTime:       xtime.Time(time.Now().Unix()),
			BatchToken:  b.BatchToken,
			Amount:      b.Amount,
			FullAmount:  b.FullAmount,
			AppID:       b.AppID,
		}
	}
	// mid should be same index.
	if aff, err = s.dao.BatchAddAllowanceCoupon(c, tx, cps); err != nil {
		return
	}
	if len(mids) != int(aff) {
		err = fmt.Errorf("batch salary midslen(%d) != aff(%d)", len(mids), aff)
		return
	}
	if _, err = s.dao.UpdateBatchInfo(c, tx, b.BatchToken, len(mids)); err != nil {
		return
	}
	return
}

// AnalysisFile analysis slary file.
func (s *Service) AnalysisFile(c context.Context, fileURL string) (midmap map[int64][]int64, total int64, err error) {
	cntb, err := ioutil.ReadFile(fileURL)
	if err != nil {
		return
	}
	var (
		mid     int64
		records [][]string
	)
	r := csv.NewReader(strings.NewReader(string(cntb)))
	records, err = r.ReadAll()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	midmap = make(map[int64][]int64, s.c.Prop.AllowanceTableCount)
	for i, v := range records {
		if len(v) <= 0 {
			continue
		}
		if mid, err = strconv.ParseInt(v[0], 10, 64); err != nil {
			err = errors.Wrapf(err, "read csv line err(%d,%v)", i, v)
			break
		}
		index := mid % s.c.Prop.AllowanceTableCount
		if midmap[index] == nil {
			midmap[index] = []int64{}
		}
		midmap[index] = append(midmap[index], mid)
		total++
	}
	return
}

// OutFile out file.
func (s *Service) OutFile(c context.Context, bs []byte, url string) (err error) {
	var outFile *os.File
	outFile, err = os.Create(url)
	if err != nil {
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	_, err = b.Write(bs)
	if err != nil {
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		os.Exit(1)
	}
	return
}

// get coupon tokeni
func (s *Service) tokeni(i int) string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%05d", i))
	b.WriteString(fmt.Sprintf("%02d", s.r.Int63n(99)))
	b.WriteString(fmt.Sprintf("%03d", time.Now().UnixNano()/1e6%1000))
	b.WriteString(time.Now().Format("20060102150405"))
	return b.String()
}

func (s *Service) sendMsg(mids []int64, ip string, times int, b bool, msgType string) (err error) {
	msg, ok := msgMap[msgType]
	if !ok {
		log.Warn("sendMsg not support msgType(%s)", msgType)
		return
	}
	for i := 1; i <= times; i++ {
		if err = s.dao.SendMessage(context.Background(), xstr.JoinInts(mids), msg["title"], msg["content"], ip); err != nil {
			if i != times {
				continue
			}
			if b {
				log.Error("batch salary msg send mids(%v) err: %+v", mids, err)
			}
		}
		break
	}
	return
}

const (
	_msgMeng = "meng"
)

var (
	msgMap = map[string]map[string]string{
		"vip": {
			"title":   "大会员代金券到账通知",
			"content": "大会员代金券已到账，快到“我的代金券”看看吧！IOS端需要在网页使用。#{\"传送门→\"}{\"https://account.bilibili.com/account/big/voucher\"}",
		},
		"meng": {
			"title":   "大会员代金券到账提醒",
			"content": "您的专属大会员代金券即将在10月10日0点到账，年费半价基础上再享折上折，请注意查收！代金券有效期1天，暂不支持苹果支付  #{\"点击查看详情\"}{\"https://www.bilibili.com/read/cv1291213\"}",
		},
	}
)
