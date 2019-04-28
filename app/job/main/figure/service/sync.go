package service

import (
	"bufio"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"go-common/library/log"
)

func (s *Service) SyncUserVIP(ctx context.Context, file string) {
	log.Info("Start sync user VIP")
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error("os.OpenFile(%s) , err [%s]", file, err)
		return
	}
	br := bufio.NewReader(f)
	var (
		str   string
		mid   int64
		count int
	)
	for {
		if str, err = br.ReadString('\n'); err != nil {
			if err == io.EOF {
				break
			}
			log.Error("br.ReadString error [%s]", err)
			return
		}
		if mid, err = strconv.ParseInt(strings.TrimSpace(str), 10, 64); err != nil {
			log.Error("Parse midstr [%s] error [%s]", str, err)
			continue
		}
		if err = s.figureDao.UpdateVipStatus(ctx, mid, 1); err != nil {
			log.Error("s.figureDao.UpdateVipStatus(%d,1) err [%s]", mid, err)
			return
		}
		count++
		if count%10000 == 0 {
			time.Sleep(time.Second)
		}
	}
	log.Info("End sync user VIP count [%d]", count)
}
