package service

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"time"

	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/library/log"
)

// VideoReports get video report record from DB
func (s *Service) VideoReports(c context.Context, t int8, stime, etime time.Time) (reports []*archive.Report, err error) {
	if reports, err = s.arc.Reports(c, t, stime, etime); err != nil {
		log.Error("s.arc.Reports(%d) err(%v)", t, err)
		return
	}
	return
}

// hdlVideoUpdateBinLog handle bilibili_archive's video table update bin log
func (s *Service) hdlVideoUpdateBinLog(nMsg, oMsg []byte) {
	var (
		nv  = &archive.Video{}
		ov  = &archive.Video{}
		err error
	)
	if err = json.Unmarshal(nMsg, nv); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nMsg, err)
		return
	}
	if err = json.Unmarshal(oMsg, ov); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", oMsg, err)
		return
	}
	if nv.Status != ov.Status {
		s.hdlVideoAudit(*nv, *ov)
	}
	if ov.XcodeState != nv.XcodeState {
		s.hdlXcodeTime(*nv, *ov)
	}
	// 视频状态变为待审核(视频信息改动或者一转完成)
	if nv.Status != ov.Status {
		if nv.Status == archive.VideoStatusWait { //待审核
			s.hdlVideoTask(context.TODO(), nv.Filename)
		}
		if nv.Status == archive.VideoStatusDelete { //视频删除
			s.arc.DelDispatch(context.TODO(), nv.Aid, nv.Cid)
		}
	}
}

// hdlVideoAudit handle video audit stats
func (s *Service) hdlVideoAudit(video, oldVideo archive.Video) {
	var (
		err error
		arc = &archive.Archive{}
	)
	if arc, err = s.arc.ArchiveByAid(context.TODO(), video.Aid); err != nil {
		log.Error("s.arc.ArchiveByAid(%d) error(%v)", video.Aid, err)
		return
	}
	s.videoAuditCache.Lock()
	defer s.videoAuditCache.Unlock()
	if _, ok := s.videoAuditCache.Data[arc.TypeID]; !ok {
		s.videoAuditCache.Data[arc.TypeID] = make(map[string]int)
	}
	switch video.Status {
	case archive.VideoStatusWait:
		s.videoAuditCache.Data[arc.TypeID]["auditing"]++
	case archive.VideoStatusOpen:
		s.videoAuditCache.Data[arc.TypeID]["audited"]++
	}
}

// hdlVideoAuditCount handle audit stats count
func (s *Service) hdlVideoAuditCount() {
	var (
		err    error
		report *archive.Report
		ctime  = time.Now()
		mtime  = ctime
		bs     []byte
	)
	if report, err = s.arc.ReportLast(context.TODO(), archive.ReportTypeVideoAudit); err != nil {
		log.Error("s.arc.ReportLast(%d) error(%v)", archive.ReportTypeVideoAudit, err)
		return
	}
	if report != nil && time.Now().Unix()-report.CTime.Unix() < 60*5 {
		log.Info("s.arc.ReportLast(%d) 距离上一次写入还没过5分钟!", archive.ReportTypeVideoAudit)
		return
	}
	s.videoAuditCache.Lock()
	defer s.videoAuditCache.Unlock()
	if bs, err = json.Marshal(s.videoAuditCache.Data); err != nil {
		log.Error("json.Marshal(%v) error(%v)", s.videoAuditCache.Data, err)
		return
	}
	if _, err = s.arc.ReportAdd(context.TODO(), archive.ReportTypeVideoAudit, string(bs), ctime, mtime); err != nil {
		log.Error("s.arc.ReportAdd(%d,%s,%v,%v) error(%v)", archive.ReportTypeVideoAudit, string(bs), ctime, mtime, err)
		return
	}
	s.videoAuditCache.Data = make(map[int16]map[string]int)
}

// VideoAudit get video audit by typeid
func (s *Service) VideoAudit(c context.Context, stime, etime time.Time) (reports []*archive.Report, err error) {
	if reports, err = s.arc.Reports(c, archive.ReportTypeVideoAudit, stime, etime); err != nil {
		log.Error("s.arc.Reports(%d) err(%v)", archive.ReportTypeVideoAudit, err)
		return
	}
	return
}

// hdlXcodeTime Stats video xcode spend time.
func (s *Service) hdlXcodeTime(nv, ov archive.Video) {
	if nv.XcodeState != archive.VideoXcodeSDFinish && nv.XcodeState != archive.VideoXcodeHDFinish && nv.XcodeState != archive.VideoDispatchFinish {
		return
	}
	var (
		nMt time.Time
		oMt time.Time
		err error
	)
	s.xcodeTimeCache.Lock()
	defer s.xcodeTimeCache.Unlock()
	if nMt, err = time.ParseInLocation("2006-01-02 15:04:05", nv.MTime, time.Local); err != nil {
		log.Error("time.ParseInLocation(%s) err(%v)", nv.MTime, err)
		return
	}
	if oMt, err = time.ParseInLocation("2006-01-02 15:04:05", ov.MTime, time.Local); err != nil {
		log.Error("time.ParseInLocation(%s) err(%v)", ov.MTime, err)
		return
	}
	t := int(nMt.Unix() - oMt.Unix())
	if t <= 0 {
		log.Info("warning: xcode spend time: %d", t)
		return
	}
	s.xcodeTimeCache.Data[nv.XcodeState] = append(s.xcodeTimeCache.Data[nv.XcodeState], t)
}

// hdlXcodeStats handle calculate and save hdlXcodeTime() stats result
func (s *Service) hdlXcodeStats() {
	var (
		c          = context.TODO()
		states     = []int8{archive.VideoXcodeSDFinish, archive.VideoXcodeHDFinish, archive.VideoDispatchFinish} //xcode states need stats
		levels     = []int8{50, 60, 80, 90}
		xcodeStats = make(map[int8]map[string]int)
		bs         []byte
		err        error
		ctime      = time.Now()
		mtime      = ctime
	)

	for _, st := range states {
		if _, ok := s.xcodeTimeCache.Data[st]; !ok {
			continue
		}
		sort.Ints(s.xcodeTimeCache.Data[st])
		seconds := s.xcodeTimeCache.Data[st]
		if len(seconds) < 1 {
			continue
		}
		for _, l := range levels {
			m := "m" + strconv.Itoa(int(l))
			o := int(math.Floor(float64(len(seconds))*(float64(l)/100)+0.5)) - 1 //seconds offset
			if o < 0 {
				continue
			}
			if o < 0 || o >= len(seconds) {
				log.Error("s.hdlVideoXcodeStats() index out of range. seconds(%d)", o)
				continue
			}
			if _, ok := xcodeStats[st]; !ok {
				xcodeStats[st] = make(map[string]int)
			}
			xcodeStats[st][m] = seconds[o]
		}
	}
	if bs, err = json.Marshal(xcodeStats); err != nil {
		log.Error("s.hdlVideoXcodeStats() json.Marshal error(%v)", err)
		return
	}
	log.Info("s.hdlVideoXcodeStats() end xcode stats xcodeStats:%s", bs)
	if len(xcodeStats) < 1 {
		log.Info("s.hdlVideoXcodeStats() end xcode stats ignore empty data")
		return
	}
	if _, err = s.arc.ReportAdd(c, archive.ReportTypeXcode, string(bs), ctime, mtime); err != nil {
		log.Error("s.hdlVideoXcodeStats() s.arc.ReportAdd error(%v)", err)
		return
	}
	s.xcodeTimeCache.Lock()
	defer s.xcodeTimeCache.Unlock()
	s.xcodeTimeCache.Data = make(map[int8][]int)
}

// hdlTraffic Calculate how long it took to check video flow in ten minutes.
// Stats result include sd_xcode,video check,hd_xcode,dispatch time.
func (s *Service) hdlTraffic() {
	var (
		err        error
		ctx        = context.TODO()
		report     *archive.Report                                                                                                        //Single report type
		reports    []*archive.Report                                                                                                      //Report type slice
		tooks      []*archive.TaskTook                                                                                                    //Task took time stats
		statsCache = make(map[int8]map[string][]int)                                                                                      //Event took time list
		traffic    = make(map[int8]map[string]int)                                                                                        //Event took time stats result
		bs         []byte                                                                                                                 //Json byte
		ctime      = time.Now()                                                                                                           //Stats create time
		mtime      = ctime                                                                                                                //Stats modify time
		states     = []int8{archive.VideoUploadInfo, archive.VideoXcodeSDFinish, archive.VideoXcodeHDFinish, archive.VideoDispatchFinish} //xcode states need stats
	)

	//0.Get the last report write time. If less than 10 minutes, then return.
	if report, err = s.arc.ReportLast(ctx, archive.ReportTypeTraffic); err != nil {
		log.Error("s.arc.ReportLast(%d) error(%v)", archive.ReportTypeTraffic, err)
		return
	}
	if report != nil && time.Now().Unix()-report.CTime.Unix() < 60*6 {
		log.Info("s.arc.ReportLast(%d) 距离上一次写入还没过6分钟!", archive.ReportTypeTraffic)
		return
	}
	now := time.Now()
	stime := now.Add(-10 * time.Minute)

	//1.Get video task time stats.
	if tooks, err = s.arc.TaskTooks(ctx, stime); err != nil {
		log.Error("s.arc.TaskTooks(%v) error(%v)", stime, err)
		return
	}
	statsCache[archive.VideoUploadInfo] = make(map[string][]int)
	for _, v := range tooks {
		statsCache[archive.VideoUploadInfo]["m50"] = append(statsCache[archive.VideoUploadInfo]["m50"], v.M50)
		statsCache[archive.VideoUploadInfo]["m60"] = append(statsCache[archive.VideoUploadInfo]["m60"], v.M60)
		statsCache[archive.VideoUploadInfo]["m80"] = append(statsCache[archive.VideoUploadInfo]["m80"], v.M80)
		statsCache[archive.VideoUploadInfo]["m90"] = append(statsCache[archive.VideoUploadInfo]["m90"], v.M90)
	}

	//2.Get sd_xcode,hd_xcode,dispatch time stats.
	if reports, err = s.arc.Reports(ctx, archive.ReportTypeXcode, stime, now); err != nil {
		log.Error("s.arc.Reports(%d) err(%v)", archive.ReportTypeXcode, err)
		return
	}
	xcodeStats := make(map[int8]map[string]int)
	for _, v := range reports {
		err = json.Unmarshal([]byte(v.Content), &xcodeStats)
		if err != nil {
			log.Error("json.Unmarshal(%s) err(%v)", v.Content, err)
			continue
		}
		for state, stats := range xcodeStats {
			if _, ok := statsCache[state]; !ok {
				statsCache[state] = make(map[string][]int)
			}
			totalTime := 0
			for level, val := range stats {
				totalTime += val
				statsCache[state][level] = append(statsCache[state][level], val)
			}
		}
	}

	//3.Calculate total time stats.
	for state, stats := range statsCache {
		for level, vals := range stats {
			total := 0
			for _, v := range vals {
				total += v
			}
			if _, ok := traffic[state]; !ok {
				traffic[state] = make(map[string]int)
			}
			traffic[state][level] = total / len(vals)
		}
	}

	//4.Save stats result
	if len(traffic) < 1 {
		log.Info("s.hdlTraffic() end traffic stats ignore empty data")
		return
	}
	if bs, err = json.Marshal(traffic); err != nil {
		log.Error("s.hdlTraffic() json.Marshal error(%v)", err)
		return
	}
	log.Info("s.hdlTraffic() end traffic stats traffic:%s", bs)
	if _, err = s.arc.ReportAdd(ctx, archive.ReportTypeTraffic, string(bs), ctime, mtime); err != nil {
		log.Error("s.hdlVideoXcodeStats() s.arc.ReportAdd error(%v)", err)
		return
	}

	//5.Update video traffic jam time
	jamTime := 0
	stateOk := true
	for _, s := range states {
		if _, ok := traffic[s]; !ok {
			stateOk = false
			break
		}
		if _, ok := traffic[s]["m60"]; !ok {
			stateOk = false
			break
		}
		if _, ok := traffic[s]["m80"]; !ok {
			stateOk = false
			break
		}
		jamTime += traffic[s]["m60"]
		jamTime += traffic[s]["m80"]
	}
	if !stateOk {
		log.Error("s.hdlTraffic() 一审耗时计算失败！traffic：%v", traffic)
	} else {
		err = s.redis.SetVideoJam(ctx, jamTime)
		log.Info("s.hdlTraffic() s.redis.SetVideoJam(%d)", jamTime)
		if err != nil {
			log.Error("s.hdlTraffic() 更新Redis失败！error(%v)", err)
		}
	}
}

func (s *Service) putVideoChan(action string, nwMsg []byte, oldMsg []byte) {
	var (
		err      error
		chanSize = int64(s.c.ChanSize)
	)
	nw := &archive.Video{}
	if err = json.Unmarshal(nwMsg, nw); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nwMsg, err)
		return
	}
	switch action {
	case _insertAct:
		s.videoUpInfoChs[nw.Aid%chanSize] <- &archive.VideoUpInfo{Nw: nw, Old: nil}
	case _updateAct:
		old := &archive.Video{}
		if err = json.Unmarshal(oldMsg, old); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", oldMsg, err)
			return
		}
		s.videoUpInfoChs[nw.Aid%chanSize] <- &archive.VideoUpInfo{Nw: nw, Old: old}
	}
}

func (s *Service) upVideoproc(k int) {
	defer s.waiter.Done()
	for {
		var (
			ok     bool
			upInfo *archive.VideoUpInfo
		)
		if upInfo, ok = <-s.videoUpInfoChs[k]; !ok {
			log.Info("s.videoUpInfoCh[%d] closed", k)
			return
		}
		s.trackVideo(upInfo.Nw, upInfo.Old)
		go s.hdlMonitorVideo(upInfo.Nw, upInfo.Old)
	}
}
