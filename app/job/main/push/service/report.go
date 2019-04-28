package service

import (
	"context"
	"runtime"
	"time"

	"go-common/app/job/main/push/dao"
	pb "go-common/app/service/main/push/api/grpc/v1"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/log"
)

const (
	_dbBatch    = 100000
	_cacheBatch = 50
)

func (s *Service) delInvalidReportsproc() {
	for {
		arg := &pb.DelInvalidReportsRequest{Type: pushmdl.DelMiFeedback}
		if _, err := s.pushRPC.DelInvalidReports(context.Background(), arg); err != nil {
			log.Error("s.pushRPC.DelInvalidReports(%d) error(%v)", arg.Type, err)
			dao.PromError("report:删除mi无效上报")
		}
		// arg = &pushmdl.ArgDelInvalidReport{Type: pushmdl.DelMiUninstalled}
		// if err := s.pushRPC.DelInvalidReports(context.Background(), arg); err != nil {
		// 	log.Error("s.pushRPC.DelInvalidReports(%d) error(%v)", arg.Type, err)
		// 	dao.PromError("report:删除mi卸载token")
		// }
		time.Sleep(time.Duration(s.c.Job.DelInvalidReportInterval))
	}
}

func (s *Service) reportproc() {
	defer s.waiter.Done()
	var err error
	for {
		msg, ok := <-s.reportCh
		if !ok {
			log.Warn("s.reportproc() closed")
			return
		}
		for _, v := range msg {
			if v == nil {
				continue
			}
			arg := &pb.AddReportRequest{
				Report: &pb.ModelReport{
					APPID:        int32(v.APPID),
					PlatformID:   int32(v.PlatformID),
					Mid:          v.Mid,
					Buvid:        v.Buvid,
					DeviceToken:  v.DeviceToken,
					Build:        int32(v.Build),
					TimeZone:     int32(v.TimeZone),
					NotifySwitch: int32(v.NotifySwitch),
					DeviceBrand:  v.DeviceBrand,
					DeviceModel:  v.DeviceModel,
					OSVersion:    v.OSVersion,
					Extra:        v.Extra,
				},
			}
			for i := 0; i < _retry; i++ {
				if _, err = s.pushRPC.AddReport(context.Background(), arg); err == nil {
					break
				}
				time.Sleep(20 * time.Millisecond)
			}
			if err != nil {
				log.Error("s.pushRPC.AddReport(%+v) error(%v)", v, err)
				dao.PromError("report:新增上报数据")
			}
			time.Sleep(time.Millisecond)
		}
	}
}

func (s *Service) refreshTokensproc() {
	for {
		now := time.Now()
		if int(now.Weekday()) != s.c.Job.SyncReportCacheWeek || int(now.Hour()) != s.c.Job.SyncReportCacheHour {
			time.Sleep(time.Minute)
			continue
		}
		s.RefreshTokenCache()
		time.Sleep(time.Hour)
	}
}

// RefreshTokenCache .
func (s *Service) RefreshTokenCache() {
	var (
		err   error
		maxid int64
		ctx   = context.Background()
	)
	for i := 0; i < _retry; i++ {
		if maxid, err = s.dao.ReportLastID(ctx); err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Error("s.refreshTokensproc() error(%v)", err)
		return
	}
	log.Info("refresh token start, maxid(%d)", maxid)
	var (
		updatedUsers  int64
		updatedTokens int64
		sli           []*pb.ModelReport
		pool          = make(map[int64][]*pb.ModelReport)
	)
	for i := int64(0); i <= maxid; i += _dbBatch {
		var rs []*pushmdl.Report
		for j := 0; j < _retry; j++ {
			if rs, err = s.dao.ReportsByRange(ctx, i, i+_dbBatch); err == nil {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}
		if err != nil {
			log.Error("s.dao.ReportsByRange(%d,%d) error(%v)", i, i+_dbBatch, err)
			continue
		}
		for _, r := range rs {
			if r.NotifySwitch == 0 {
				continue
			}
			nr := &pb.ModelReport{
				APPID:        int32(r.APPID),
				PlatformID:   int32(r.PlatformID),
				Mid:          r.Mid,
				Buvid:        r.Buvid,
				DeviceToken:  r.DeviceToken,
				Build:        int32(r.Build),
				TimeZone:     int32(r.TimeZone),
				NotifySwitch: int32(r.NotifySwitch),
				DeviceBrand:  r.DeviceBrand,
				DeviceModel:  r.DeviceModel,
				OSVersion:    r.OSVersion,
				Extra:        r.Extra,
			}
			sli = append(sli, nr)
			if len(sli) >= _cacheBatch {
				s.addTokensCache(sli)
				sli = []*pb.ModelReport{}
			}
			if r.Mid == 0 {
				continue
			}
			pool[r.Mid] = append(pool[r.Mid], nr)
			updatedTokens++
		}
		log.Info("refresh token sovled min(%d) max(%d)", i, i+_dbBatch)
		time.Sleep(time.Millisecond)
	}
	if len(sli) > 0 {
		s.addTokensCache(sli)
	}
	log.Info("refresh token data, users(%d) tokens(%d)", len(pool), updatedTokens)
	for mid, rs := range pool {
		arg := &pb.AddUserReportCacheRequest{Mid: mid, Reports: rs}
		for i := 0; i < _retry; i++ {
			if _, err = s.pushRPC.AddUserReportCache(ctx, arg); err == nil {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		if err != nil {
			log.Error("s.pushRPC.AddUserReportCache(%d) error(%v)", mid, err)
			continue
		}
		updatedUsers++
		delete(pool, mid)
	}
	pool = nil
	runtime.GC()
	log.Info("refresh token end, updated users(%d) tokens(%d)", updatedUsers, updatedTokens)
}

func (s *Service) addTokensCache(rs []*pb.ModelReport) (err error) {
	arg := new(pb.AddTokensCacheRequest)
	arg.Reports = append(arg.Reports, rs...)
	for i := 0; i < _retry; i++ {
		if _, err = s.pushRPC.AddTokensCache(context.Background(), arg); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		log.Error("s.pushRPC.AddTokensCache tokens(%d) error(%v)", len(rs), err)
		return
	}
	log.Info("s.pushRPC.AddTokensCache tokens(%d)", len(rs))
	return
}

func (s *Service) tokensByMids(task *pushmdl.Task, mids []int64) (res map[int][]string, valid int64, err error) {
	rs, _, err := s.dao.ReportsCacheByMids(context.Background(), mids)
	if err != nil {
		log.Error("s.dao.ReportsCacheByMids() error(%v)", err)
		return
	}
	var (
		exist = make(map[int64]bool, len(rs))
		// platformCount = len(task.Platform)
		buildCount = len(task.Build)
	)
	for mid := range rs {
		exist[mid] = true
	}
	for _, mid := range mids {
		if !exist[mid] {
			log.Warn("tokens by mid, task(%s) mid(%d)", task.ID, mid)
		}
	}
	res = make(map[int][]string)
	for _, rr := range rs {
		for _, r := range rr {
			if r.APPID != task.APPID {
				continue
			}
			if r.NotifySwitch == pushmdl.SwitchOff {
				continue
			}
			realTime := pushmdl.RealTime(r.TimeZone)
			if realTime.Unix() > int64(task.ExpireTime) {
				continue
			}
			// if platformCount > 0 && !validatePlatform(r.PlatformID, task.Platform) {
			// 	continue
			// }
			if buildCount > 0 && !pushmdl.ValidateBuild(r.PlatformID, r.Build, task.Build) {
				continue
			}
			res[r.PlatformID] = append(res[r.PlatformID], r.DeviceToken)
		}
		valid++
	}
	return
}
