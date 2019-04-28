package service

import (
	"context"
	"fmt"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/log"
	"os"
	"time"
)

//taskSyncUserDmg 同步用户画像
func (s *Service) taskSyncUserDmg() {
	jobURL, err := s.dao.QueryUserDmg(context.Background())
	if err != nil {
		log.Error("get user dmg err(%v)", err)
		return
	}
	urls, err := s.dao.QueryJobStatus(context.Background(), jobURL)
	if err != nil {
		log.Error("get user dmg job result err(%v)", err)
		return
	}
	for _, url := range urls {
		go func(url string) {
			fpath, err := s.dao.Download(url, "")
			if err != nil {
				return
			}
			s.dao.ReadLine(fpath, s.dao.HandlerUserDmg)
			os.RemoveAll(fpath)
		}(url)
	}
}

//UserProfileUpdate bbq_user_profile
func (s *Service) UserProfileUpdate() {
	_, err := s.dao.UserProfileGet(context.Background())
	if err != nil {
		log.Error("get user dmg err(%v)", err)
	}
}

//taskSyncUserDmg 同步天马推荐用户画像
func (s *Service) taskSyncPegasusUserBasic() {
	_, err := s.dao.QueryUserBasic(context.Background())
	if err != nil {
		log.Error("get user dmg err(%v)", err)
	}
}

//taskSyncUpUserDmg 同步up主画像
func (s *Service) taskSyncUpUserDmg() {
	log.Infov(context.Background(), log.KV("event", "sync_up_user_dmg"))
	var mid = int64(0)
	for {
		log.Infov(context.Background(), log.KV("event", "one_sync_up_user_dmg"), log.KV("mid", mid))
		upUserDmgs, err := s.dao.QueryUpUserDmg(context.Background(), mid)
		if err != nil {
			log.Error("get up user dmg err(%v)", err)
			return
		}

		if len(upUserDmgs) == 0 {
			break
		}

		for _, upUserDmg := range upUserDmgs {
			mid = upUserDmg.MID
			fmt.Println(upUserDmg.MID)
			if err = s.dao.InsertOnDup(context.Background(), upUserDmg); err != nil {
				log.Error("user dmg insert on dup failed,mid :%v, err:%v", upUserDmg.MID, err)
			}
		}
	}
	// TODO: 考虑个好方法
	// 万一对方接口有问题，那就都完了
	//s.dao.DelUpUserDmg(context.Background())
}

// 同步up主画像 从hive更新到user_statistics_hive
func (s *Service) taskSyncUsrStaFromHive() {
	log.Infov(context.Background(), log.KV("event", "taskSyncUsrStaFromHive"))
	var (
		err    error
		jobURL string
		urls   []string
		url    string
		fpath  string
		try    int
		date   = time.Now().AddDate(0, 0, -1).Format("20060102")
	)
	for try = 1; try <= 3; try++ {
		//发起hive查询,拿到url
		if jobURL, err = s.dao.QueryUpMid(context.Background(), date); err != nil {
			log.Warn("taskSyncUsrStaFromHive try and init query hive failed, err(%v)", err)
			continue
		}
		log.Info("taskSyncUsrStaFromHive init query hive success")
		//查询job状态
		if urls, err = s.dao.QueryJobStatus(context.Background(), jobURL); err != nil {
			log.Warn("taskSyncUsrStaFromHive try and get hive query status failed, err(%v)", err)
			continue
		}
		break
	}
	if err != nil {
		log.Error("taskSyncUsrStaFromHive init and get hive query status failed, err(%v)", err)
		return
	}
	for _, url = range urls {
		for try = 0; try <= 3; try++ {
			if fpath, err = s.dao.Download(url, ""); err != nil {
				log.Warn("taskSyncUsrInfoFromHive try and download file (%v) failed, err(%v)", url, err)
				time.Sleep(time.Duration(try*10) * time.Second)
				continue
			}
			s.dao.ReadLines(fpath, s.dao.HandlerMids)
			os.RemoveAll(fpath)
			return
		}
		if err != nil {
			log.Error("taskSyncUsrInfoFromHive download file (%v) failed, err(%v)", url, err)
		}
	}
}

//从video表同步up画像
//func (s *Service) taskSyncUsrBaseFromVideo(c context.Context) {
//	log.Infov(context.Background(), log.KV("event", "taskSyncUsrInfoFromVideo"))
//	//get mids
//	var (
//		mids *[]int64
//		err  error
//		i    int8
//		mid  int64
//		req  *video.SyncUserBaseResponse
//	)
//	for i := 0; i <= 3; i++ {
//		mids, err = s.dao.SelMidFromVideo()
//		if err != nil {
//			log.Info("taskSyncUsrInfoFromVideo try and get up mid failed , err(%v)", err)
//		} else {
//			break
//		}
//	}
//	if err != nil {
//		log.Error("taskSyncUsrInfoFromVideo get up mid failed, err(%v)", err)
//	}
//	//get userinfo and update
//	for _, mid = range *mids {
//		//重试
//		for i = 0; i <= 3; i++ {
//			if req, err = s.dao.VideoClient.SyncUserBase(c, &video.SyncMidRequset{MID: mid}); err != nil || req.Affc == -1 {
//				log.Info("taskSyncUsrInfoFromVideo try and failed, mid(%v), err(%v)", mid, err)
//			} else {
//				break
//			}
//		}
//		if err != nil {
//			log.Error("taskSyncUsrInfoFromVideo failed, mid(%v), err(%v)", mid, err)
//		}
//	}
//}

//从video表同步up画像
func (s *Service) taskSyncUsrBaseFromVideo(c context.Context) {
	fmt.Println("taskSyncUsrBaseFromVideo start")
	log.Infov(context.Background(), log.KV("event", "taskSyncUsrInfoFromVideo"))
	//get mids
	var (
		mids []int64
		err  error
		req  *video.SyncUserBaseResponse
	)
	for try := 0; try <= 3; try++ {
		mids, err = s.dao.SelMidFromVideo()
		if err != nil {
			log.Info("taskSyncUsrInfoFromVideo try and get up mid failed , err(%v)", err)
		} else {
			break
		}
	}
	if err != nil {
		log.Error("taskSyncUsrInfoFromVideo get up mid failed, err(%v)", err)
	}
	i := len(mids) / 50
	//get userinfo and update
	for j := 1; j <= i; j++ {
		for try := 0; try <= 3; try++ {
			if req, err = s.dao.VideoClient.SyncUserBases(c, &video.SyncMidsRequset{MIDS: mids[(j-1)*50 : j*50]}); err != nil {
				log.Info("taskSyncUsrInfoFromVideo try and failed, err(%v)", err)
			} else {
				break
			}
		}
		if err != nil {
			log.Error("taskSyncUsrInfoFromVideo failed, err(%v)", err)
		} else {
			log.Info("taskSyncUsrInfoFromVideo success ,affected %v rows", req.Affc)
		}
	}
	if i*50 < len(mids) {
		for try := 0; try <= 3; try++ {
			if req, err = s.dao.VideoClient.SyncUserBases(c, &video.SyncMidsRequset{MIDS: mids[i*50:]}); err != nil {
				log.Info("taskSyncUsrInfoFromVideo try and failed, err(%v)", err)
			} else {
				break
			}
		}
		if err != nil {
			log.Error("taskSyncUsrInfoFromVideo failed, err(%v)", err)
		} else {
			log.Info("taskSyncUsrInfoFromVideo success ,affected %v rows", req.Affc)
		}
	}
}

// UpdateUsrBaseFace 更新user_base里的face字段
func (s *Service) UpdateUsrBaseFace() (err error) {
	var (
		mids []int64
	)
	log.Infov(context.Background(), log.KV("event", "UpdateUsrBaseFace"))
	for i := 0; ; i++ {
		mids, err = s.dao.SelMidFromUserBase(i * 1000)
		if err != nil {
			log.Error("UpdateUsrBaseFace select mid failed")
			return
		}
		if len(mids) == 0 {
			break
		}
		s.dao.UpUserBases(context.Background(), mids)
	}
	return
}
