package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"go-common/app/job/main/passport-game-data/conf"
	"go-common/app/job/main/passport-game-data/dao"
	"go-common/app/job/main/passport-game-data/model"
	"go-common/library/log"
)

const (
	_defaultInitCloudOffsetFilePath = "/data/passport-game-data-job.initcloud.offset"

	_defaultInitCloudSleep = time.Second
)

type initCloudConfig struct {
	OffsetFilePath string
	UseOldOffset   bool

	Start, End int64

	Batch int

	Sleep time.Duration
}

func newInitCloudConfigFrom(c *conf.Config) (ic *initCloudConfig) {
	ic = &initCloudConfig{
		OffsetFilePath: c.InitCloud.OffsetFilePath,
		UseOldOffset:   c.InitCloud.UseOldOffset,

		Start: c.InitCloud.Start,
		End:   c.InitCloud.End,

		Batch: c.InitCloud.Batch,

		Sleep: time.Duration(c.InitCloud.Sleep),
	}

	ic.fix()

	if ic.UseOldOffset {
		data, err := ioutil.ReadFile(ic.OffsetFilePath)
		if err != nil {
			log.Error("failed to read old offset, skip")
			return
		}

		oldOffset, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			log.Error("failed to parse offset, strconv.ParseInt(%s, 10, 64)", string(data), err)
			return
		}

		if oldOffset > 0 {
			ic.Start = oldOffset
		}
	}
	return
}

func (ic *initCloudConfig) fix() {
	if len(ic.OffsetFilePath) == 0 {
		ic.OffsetFilePath = _defaultInitCloudOffsetFilePath
	}
	if ic.Start < 0 {
		ic.Start = 0
	}
	if ic.End < 0 {
		ic.End = 0
	}

	if ic.Batch <= 0 {
		ic.Batch = _defaultBatchSize
	}

	if int64(ic.Sleep) < 0 {
		ic.Sleep = _defaultInitCloudSleep
	}
}

// NewInitCloud new a service for initiating cloud.
func NewInitCloud(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		d: dao.New(c),

		ic: newInitCloudConfigFrom(c),
	}
	return
}

// InitCloud init cloud.
func (s *Service) InitCloud(c context.Context) {
	var err error
	ic := s.ic

	dstFile, err := os.Create(ic.OffsetFilePath)
	if err != nil {
		log.Error("failed to open file %s, error(%v)", ic.OffsetFilePath, err)
		return
	}
	defer dstFile.Close()

	for i := ic.Start; i <= ic.End; {
		time.Sleep(ic.Sleep)
		if err = ioutil.WriteFile(ic.OffsetFilePath, []byte(strconv.FormatInt(i, 10)), os.ModeAppend); err != nil {
			log.Error("failed to record offset, offsetFilePath: %s, offset: %d, error(%v)", ic.OffsetFilePath, i, err)
			continue
		}
		st := i
		ed := i + int64(ic.Batch)
		if ed > ic.End {
			ed = ic.End
		}

		mids := make([]int64, 0)
		for j := st; j <= ed; j++ {
			mids = append(mids, j)
		}

		var as []*model.OriginAsoAccount
		if as, err = s.d.AsoAccountsLocal(c, mids); err != nil {
			log.Error("failed to get local aso accounts by mids, service.dao.AsoAccountsLocal(%v) error(%v)", mids, err)
			continue
		}

		cloudAs := make([]*model.AsoAccount, 0)
		for _, a := range as {
			cloudAs = append(cloudAs, model.Default(a))
		}
		if err = s.d.AddAsoAccountsCloud(c, cloudAs); err != nil {
			str, _ := json.Marshal(cloudAs)
			log.Error("failed to add aso accounts to cloud, service.dao.AddAsoAccountsCloud(%v) error(%v)", str, err)
			continue
		}
		i += int64(ic.Batch)
	}
}
