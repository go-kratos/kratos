package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"time"

	coinclient "go-common/app/service/main/coin/api"
	"go-common/library/log"
)

// FormatCSV  format csv data.
func (s *Service) FormatCSV(records [][]string) (res []byte) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Error("error(%+v) writing record to csv:", err)
			return
		}
	}
	w.Flush()
	res = buf.Bytes()
	return
}

// Upload http upload file.
func (s *Service) Upload(c context.Context, fileName, fileType string, expire int64, body io.Reader) (location string, err error) {
	if location, err = s.uploadDao.Upload(c, fileName, fileType, expire, body); err != nil {
		log.Error("s.upload.Upload() error(%v)", err)
	}
	return
}

// AnnualCoins .
func (s *Service) AnnualCoins(c context.Context, reader *csv.Reader) (fmid []int64) {
	for {
		record, err := reader.Read()
		if err == io.EOF {
			log.Warn("AnnualCoins is over!")
			err = nil
			break
		}
		if err != nil {
			log.Error("AnnualCoins(%+v) Error(%+v)", record, err)
			err = nil
			continue
		}
		log.Info("AnnualCoins record(%+v)", record)
		if len(record) < 2 {
			log.Error("AnnualCoins record(%+v) len<2", record)
			continue
		}
		mid, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			log.Error("AnnualCoins strconv.ParseInt mid(%+v) err(%+v)", record, err)
			err = nil
			continue
		}
		coins, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Error("AnnualCoins strconv.ParseInt coins(%+v) err(%+v)", record, err)
			err = nil
			continue
		}
		arg := &coinclient.ModifyCoinsReq{
			Mid:       mid,
			Count:     float64(coins),
			Reason:    "风纪委奖励",
			IP:        "",
			Operator:  "credit",
			CheckZero: 1,
			Ts:        time.Now().Unix(),
		}
		_, err = s.coinClient.ModifyCoins(context.Background(), arg)
		if err != nil {
			fmid = append(fmid, mid)
			log.Error("ModifyCoins arg(%+v), err(%+v)", arg, err)
			continue
		}
		if err = s.msgDao.SendSysMsg(context.Background(), mid, "您的风纪委硬币礼包已到账", "Hi 风纪委员你好，你的2018年风纪委硬币礼包已到账！快进入你的硬币账户看看吧！"); err != nil {
			log.Error("SendSysMsg mid(%d), err(%+v)", mid, err)
		}
		time.Sleep(time.Duration(50) * time.Millisecond)
	}
	return
}
