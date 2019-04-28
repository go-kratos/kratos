package service

import (
	"bufio"
	"context"
	"fmt"
	"go-common/app/job/bbq/video/dao"
	searchv1 "go-common/app/service/bbq/search/api/grpc/v1"
	"go-common/library/log"
	"io"
	"net/url"
	"os"

	"strconv"
	"time"
)

const (
	_retryTimes = 3

	_selection = 5 //运营精选状态
)

//taskCheckVideoDBVSES 视频全量脚本
func (s *Service) taskSyncVideo2ES() {
	var step time.Duration
	var id int64
	for {
		ids, videos, err := s.dao.VideoList(context.Background(), id)
		if err != nil {
			log.Error("sync video err(%v)", err)
			return
		}
		if videos == nil {
			return
		}
		videoStatisticsHive, _ := s.dao.VideoStatisticsHiveList(context.Background(), ids)
		videoStatistics, _ := s.dao.VideoStatisticsList(context.Background(), ids)
		videoTags, _ := s.dao.VideoTagsList(context.Background(), ids)
		req := new(searchv1.SaveVideoRequest)
		for _, v := range videos {
			fmt.Println(v.SVID)
			id = v.SVID
			tmp := &searchv1.VideoESInfo{
				SVID:         v.SVID,
				Title:        v.Title,
				Content:      v.Content,
				MID:          v.MID,
				CID:          v.CID,
				Pubtime:      int64(v.Pubtime),
				Ctime:        int64(v.Ctime),
				Mtime:        int64(v.Mtime),
				Duration:     v.Duration,
				Original:     v.Original,
				State:        v.State,
				VerID:        v.VerID,
				Ver:          v.Ver,
				From:         v.From,
				AVID:         v.AVID,
				Tid:          v.Tid,
				SubTid:       v.SubTid,
				ISFullScreen: v.ISFullScreen,
				Score:        v.Score,
			}
			if videoStatisticsHive[id] != nil {
				tmp.PlayHive = videoStatisticsHive[id].PlayHive
				tmp.FavHive = videoStatisticsHive[id].FavHive
				tmp.CoinHive = videoStatisticsHive[id].CoinHive
				tmp.SubtitlesHive = videoStatisticsHive[id].SubtitlesHive
				tmp.LikesHive = videoStatisticsHive[id].LikesHive
				tmp.ShareHive = videoStatisticsHive[id].ShareHive
				tmp.ReportHive = videoStatisticsHive[id].ReportHive
				tmp.DurationDailyHive = videoStatisticsHive[id].DurationDailyHive
				tmp.DurationAllHive = videoStatisticsHive[id].DurationAllHive
				tmp.ReplyHive = videoStatisticsHive[id].ReplyHive
				tmp.ShareDailyHive = videoStatisticsHive[id].ShareDailyHive
				tmp.PlayDailyHive = videoStatisticsHive[id].PlayDailyHive
				tmp.SubtitlesDailyHive = videoStatisticsHive[id].SubtitlesDailyHive
				tmp.LikesDailyHive = videoStatisticsHive[id].LikesDailyHive
				tmp.FavDailyHive = videoStatisticsHive[id].FavDailyHive
				tmp.ReplyDailyHive = videoStatisticsHive[id].ReplyDailyHive
				tmp.AccessHive = videoStatisticsHive[id].AccessHive
			}
			if videoStatistics[id] != nil {
				tmp.Play = videoStatistics[id].Play
				tmp.Subtitles = videoStatistics[id].Subtitles
				tmp.Like = videoStatistics[id].Like
				tmp.Share = videoStatistics[id].Share
				tmp.Report = videoStatistics[id].Report
			}
			if videoTags[id] != nil {
				tmp.Tags = videoTags[id]
			}
			req.List = append(req.List, tmp)
		}
		step = 1
		for {
			if _, err := s.dao.SearchClient.SaveVideo(context.Background(), req); err != nil {
				time.Sleep(step * time.Second)
				step++
				continue
			}
			break
		}
	}

}

//SaveVideo2ES 保存视频到es
func (s *Service) SaveVideo2ES(ids string) (res bool) {
	res = true
	if len(ids) == 0 {
		return
	}
	videos, err := s.dao.VideoListByIDs(context.Background(), ids)
	if err != nil || videos == nil {
		res = false
		return
	}
	videoStatisticsHive, _ := s.dao.VideoStatisticsHiveList(context.Background(), ids)
	videoStatistics, _ := s.dao.VideoStatisticsList(context.Background(), ids)
	// videoTags, _ := s.dao.VideoTagsList(context.Background(), ids)
	var step time.Duration
	var id int64
	req := new(searchv1.SaveVideoRequest)
	for _, v := range videos {
		id = v.SVID
		fmt.Println(id)
		tmp := &searchv1.VideoESInfo{
			SVID:         v.SVID,
			Title:        v.Title,
			Content:      v.Content,
			MID:          v.MID,
			CID:          v.CID,
			Pubtime:      int64(v.Pubtime),
			Ctime:        int64(v.Ctime),
			Mtime:        int64(v.Mtime),
			Duration:     v.Duration,
			Original:     v.Original,
			State:        v.State,
			VerID:        v.VerID,
			Ver:          v.Ver,
			From:         v.From,
			AVID:         v.AVID,
			Tid:          v.Tid,
			SubTid:       v.SubTid,
			ISFullScreen: v.ISFullScreen,
			Score:        v.Score,
		}
		if videoStatisticsHive[id] != nil {
			tmp.PlayHive = videoStatisticsHive[id].PlayHive
			tmp.FavHive = videoStatisticsHive[id].FavHive
			tmp.CoinHive = videoStatisticsHive[id].CoinHive
			tmp.SubtitlesHive = videoStatisticsHive[id].SubtitlesHive
			tmp.LikesHive = videoStatisticsHive[id].LikesHive
			tmp.ShareHive = videoStatisticsHive[id].ShareHive
			tmp.ReportHive = videoStatisticsHive[id].ReportHive
			tmp.DurationDailyHive = videoStatisticsHive[id].DurationDailyHive
			tmp.DurationAllHive = videoStatisticsHive[id].DurationAllHive
			tmp.ReplyHive = videoStatisticsHive[id].ReplyHive
			tmp.ShareDailyHive = videoStatisticsHive[id].ShareDailyHive
			tmp.PlayDailyHive = videoStatisticsHive[id].PlayDailyHive
			tmp.SubtitlesDailyHive = videoStatisticsHive[id].SubtitlesDailyHive
			tmp.LikesDailyHive = videoStatisticsHive[id].LikesDailyHive
			tmp.FavDailyHive = videoStatisticsHive[id].FavDailyHive
			tmp.ReplyDailyHive = videoStatisticsHive[id].ReplyDailyHive
			tmp.AccessHive = videoStatisticsHive[id].AccessHive
		}
		if videoStatistics[id] != nil {
			tmp.Play = videoStatistics[id].Play
			tmp.Subtitles = videoStatistics[id].Subtitles
			tmp.Like = videoStatistics[id].Like
			tmp.Share = videoStatistics[id].Share
			tmp.Report = videoStatistics[id].Report
		}
		// if videoTags[id] != nil {
		// 	tmp.Tags = videoTags[id]
		// }
		req.List = append(req.List, tmp)
	}
	step = 1
	for {
		if _, err := s.dao.SearchClient.SaveVideo(context.Background(), req); err != nil {
			if step == 11 {
				log.Error("save es err(%v) ids(%s)", err, ids)
				res = false
				break
			}
			time.Sleep(step * time.Second)
			step++
			continue
		}
		break
	}
	return
}

func formArrayString(arr []int64) string {
	var res string
	for i, v := range arr {
		if i != 0 {
			res += ","
		}
		res += strconv.FormatInt(v, 10)
	}
	return res
}

//deltaSync2ES 为不同表进行增量同步的脚本，baseTableQuery指明不同表的查询语句
func (s *Service) deltaSync2ES(taskName string, baseTableQuery string) {
	task, err := s.dao.RawCheckTask(context.Background(), taskName)
	if err != nil {
		log.Error("get last_chek_time fail: task=%s", taskName)
		return
	}
	log.Info("get last_chek_time succ: task=%s, last_check_time=%d", taskName, task.LastCheck)

	// 获得所有变更的svid
	ids, mtime, err := s.dao.RawGetIDByMtime(baseTableQuery, task.LastCheck)
	if err != nil {
		log.Error("get raw id by mtime fail: task=%s, last_mtime=%d, base_table_query=%s",
			taskName, mtime, baseTableQuery)
		return
	}
	idsNum := len(ids)
	log.Info("get changed svids: task=%s, id_num=%d", taskName, idsNum)
	if idsNum == 0 {
		return
	}
	task.LastCheck = mtime

	// 对所有变更的svid分批次进行同步到es
	for i := 0; i < idsNum; i += dao.MaxSyncESNum {
		last := i + dao.MaxSyncESNum
		if last > idsNum {
			last = idsNum
		}
		selectedIDs := ids[i:last]
		idsStr := formArrayString(selectedIDs)
		if res := s.SaveVideo2ES(idsStr); !res {
			log.Error("sync video 2 es fail: task=%s, offset=%d, id_num=%d, base_table_query=%s",
				taskName, i, idsNum, baseTableQuery)
			return
		}
		log.Info("one sync video 2 es: task=%s, offset=%d, id_num=%d", taskName, i, idsNum)
	}

	// 更新task最近check的时间点
	if _, err := s.dao.UpdateTaskLastCheck(context.Background(), taskName, task.LastCheck); err != nil {
		log.Error("update task last check time fail: task=%s, last_mtime=%d, base_table_query=%s",
			taskName, task.LastCheck, baseTableQuery)
		return
	}
	log.Info("sync video 2 es: task=%s, id_num=%d, last_mtime=%d", taskName, idsNum, task.LastCheck)
}

//taskCheckVideo video表增量脚本
func (s *Service) taskCheckVideo() {
	taskName := "checkVideo"
	s.deltaSync2ES(taskName, dao.QueryVideoByMtime)
}

//taskCheckVideoStatistics video_statistics表增量脚本
func (s *Service) taskCheckVideoStatistics() {
	taskName := "checkVideoSt"
	s.deltaSync2ES(taskName, dao.QueryVideoStatisticsByMtime)
}

//taskCheckVideoStatisticsHive video_statistics_hive表增量脚本
func (s *Service) taskCheckVideoStatisticsHive() {
	taskName := "checkVideoStHv"
	s.deltaSync2ES(taskName, dao.QueryVideoStatisticsHiveByMtime)
}

//taskCheckVideoTag video_tag表增量脚本
func (s *Service) taskCheckVideoTag() {
	taskName := "checkVideoTag"
	s.deltaSync2ES(taskName, dao.QueryVideoTagByMtime)
}

//taskCheckTag tag表增量脚本
func (s *Service) taskCheckTag() {
	taskName := "checkTag"
	task, err := s.dao.RawCheckTask(context.Background(), taskName)
	if err != nil {
		log.Error("get last_chek_time fail: task=%s", taskName)
		return
	}
	log.Info("get last_chek_time succ: task=%s, last_check_time=%d", taskName, task.LastCheck)
	for {
		ids, mtime, err := s.dao.RawTagByMtime(context.Background(), task.LastCheck)
		if err != nil || len(ids) == 0 {
			return
		}
		id := int64(0)
		for {
			svids, temp, err := s.dao.RawVideoTagByIDs(context.Background(), ids, id)
			if err != nil {
				return
			}
			if len(svids) == 0 {
				break
			}
			if flag := s.SaveVideo2ES(svids); !flag {
				return
			}
			id = temp
		}
		if num, err := s.dao.UpdateTaskLastCheck(context.Background(), taskName, mtime); err != nil || num == 0 {
			return
		}
		task.LastCheck = mtime
	}
}

// taskRmInvalidES 删除es中多余的视频
func (s *Service) taskRmInvalidES() {
	fmt.Println("aaa")
	esReq := new(searchv1.ESVideoDataRequest)
	delReq := new(searchv1.DelVideoBySVIDRequest)
	svid := int64(0)
	query := `{"query":{"range":{"svid":{"gt":%d}}},"sort":[{"svid":"asc"}],"from":0,"size":10}`
	for {
		esReq.Query = fmt.Sprintf(query, svid)
		res, err := s.dao.SearchClient.ESVideoData(context.Background(), esReq)
		if err != nil {
			return
		}
		svids := make([]string, 0)
		for _, v := range res.List {
			svids = append(svids, strconv.Itoa(int(v.SVID)))
			svid = v.SVID
		}
		vs, err := s.dao.RawVideoBySVIDS(context.Background(), svids)
		if err != nil {
			return
		}
		notList := make([]int64, 0)
		for _, v := range res.List {
			if _, ok := vs[v.SVID]; !ok {
				fmt.Println(v.SVID)
				notList = append(notList, v.SVID)
			}
		}
		if len(notList) != 0 {
			delReq.SVIDs = notList
			s.dao.SearchClient.DelVideoBySVID(context.Background(), delReq)
		}
	}
}

func (s *Service) commitCID() {
	ctx := context.Background()
	path := s.c.URLs["bvc_push"]
	if path == "" {
		return
	}
	srcPath := s.c.Path["cids"]
	if srcPath == "" {
		return
	}
	if srcPath == "" {
		log.Error("sugsrc path is empty")
		return
	}
	src, err := os.Open(srcPath)
	if err != nil {
		log.Error("writeSug os.Open source sug error(%v)", err)
		return
	}
	defer src.Close()
	br := bufio.NewReader(src)
	i := 1
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		cid, err := strconv.ParseInt(string(a), 10, 64)
		if err != nil {
			log.Error("parse err [%v]", err)
			continue
		}
		svid, err := s.dao.GetSvidByCid(ctx, cid)
		if err != nil {
			continue
		}
		params := url.Values{}
		params.Set("svid", strconv.FormatInt(svid, 10))
		params.Set("cid", string(a))
		req, err := s.dao.HTTPClient.NewRequest("GET", path, "", params)
		if err != nil {
			log.Error("error(%v)", err)
			continue
		}
		var res struct {
			Code int    `json:"code"`
			Msg  string `json:"message"`
		}
		if err = s.dao.HTTPClient.Do(ctx, req, &res); err != nil {
			log.Errorv(ctx, log.KV("log", fmt.Sprintf("err[%v]", err)))
			continue
		}
		if res.Code != 0 {
			log.Errorv(ctx, log.KV("log", fmt.Sprintf("error(%v)", err)))
		} else {
			log.Info("commit svid:%d cid:%d success No.%d", svid, cid, i)
		}
		i++
	}
}
