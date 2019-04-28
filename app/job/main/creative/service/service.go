package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/creative/conf"
	"go-common/app/job/main/creative/dao/academy"
	"go-common/app/job/main/creative/dao/archive"
	"go-common/app/job/main/creative/dao/monitor"
	"go-common/app/job/main/creative/dao/newcomer"
	"go-common/app/job/main/creative/dao/weeklyhonor"
	"go-common/app/job/main/creative/model"
	"go-common/library/conf/env"
	"go-common/library/queue/databus"
	"go-common/library/xstr"

	"github.com/robfig/cron"
)

// Service is service.
type Service struct {
	c *conf.Config
	//arc  databus
	arcSub       *databus.Databus
	arcNotifySub *databus.Databus
	upPub        *databus.Databus
	// wait group
	wg sync.WaitGroup
	// monitor
	monitor     *monitor.Dao
	arc         *archive.Dao
	arcNotifyMo int64
	arcMo       int64
	// chan for mids
	midsChan chan map[int64]int
	//aca
	aca *academy.Dao
	// honDao
	honDao *weeklyhonor.Dao
	//task databus
	newc                                           *newcomer.Dao
	taskSub, shareSub, relationSub, statLikeSub    *databus.Databus
	statShareSub, statCoinSub, statFavSub          *databus.Databus
	statReplySub, statDMSub, statViewSub, newUpSub *databus.Databus
	taskSubQueue                                   []chan *databus.Message
	shareSubQueue                                  []chan *model.ShareMsg
	relationQueue                                  []chan *model.Relation //用户关注队列
	followerQueue                                  []chan *model.Stat     //粉丝数队列
	newUpQueue                                     []chan *model.Up       //新投稿
	oldUpQueue                                     []chan *model.Up       //进阶任务视频投稿超过5个
	mobileUpQueue                                  []chan *model.Up       //手机投稿

	databusQueueLen  int //消费databus 队列长度
	statViewQueueLen int //播放消费databus 队列长度
	statLikeQueueLen int //点赞消费databus 队列长度
	chanSize         int //chan 缓冲长度
	//单个稿件计数
	statViewSubQueue  []chan *model.StatView
	statLikeSubQueue  []chan *model.StatLike
	statShareSubQueue []chan *model.StatShare
	statCoinSubQueue  []chan *model.StatCoin
	statFavSubQueue   []chan *model.StatFav
	statReplySubQueue []chan *model.StatReply
	statDMSubQueue    []chan *model.StatDM
	//db
	taskQueue       []chan []*model.UserTask
	TaskCache       []*model.Task
	TaskMapCache    map[int64]*model.Task
	GiftRewardCache map[int8][]*model.GiftReward //gift-reward
	//notify
	taskNotifyQueue   []chan []int64
	rewardNotifyQueue []chan []int64
	testNotifyMids    map[int64]struct{}
}

// New is go-common/app/service/videoup service implementation.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:            c,
		arcSub:       databus.New(c.ArcSub),
		arcNotifySub: databus.New(c.ArcNotifySub),
		upPub:        databus.New(c.UpPub),
		monitor:      monitor.New(c),
		midsChan:     make(chan map[int64]int, c.ChanSize),
		arc:          archive.New(c),
		aca:          academy.New(c),
		honDao:       weeklyhonor.New(c),
		//task
		newc:            newcomer.New(c),
		taskQueue:       make([]chan []*model.UserTask, c.Task.TableConsumeNum),
		TaskMapCache:    make(map[int64]*model.Task),
		GiftRewardCache: make(map[int8][]*model.GiftReward),
		//notify
		taskNotifyQueue:   make([]chan []int64, c.Task.TaskTableConsumeNum),
		rewardNotifyQueue: make([]chan []int64, c.Task.RewardTableConsumeNum),
		testNotifyMids:    make(map[int64]struct{}),
		chanSize:          c.Task.ChanSize,
	}
	s.newTaskDatabus()
	if c.Consume {
		s.wg.Add(1)
		go s.arcCanalConsume()
		s.wg.Add(1)
		go s.arcNotifyCanalConsume()
		go s.monitorConsume()
	}
	if c.HotSwitch {
		go s.FlushHot(model.BusinessForArchvie) //计算视频稿件的hot
		go s.FlushHot(model.BusinessForArticle) //计算专栏稿件的hot
	}
	s.taskConsume()
	return
}

// InitCron init cron
func (s *Service) InitCron() {
	c := cron.New()
	if s.c.HonorSwitch {
		c.AddFunc(s.c.HonorFlushSpec, s.FlushHonor)
		c.AddFunc(s.c.HonorMSGSpec, s.SendMsg)
	}
	c.Start()
}

func (s *Service) newTaskDatabus() {
	s.databusQueueLen = s.c.Task.DatabusQueueLen
	s.statViewQueueLen = s.c.Task.StatViewQueueLen //播放
	s.statLikeQueueLen = s.c.Task.StatLikeQueueLen //点赞

	s.taskSub = databus.New(s.c.TaskSub)
	s.shareSub = databus.New(s.c.ShareSub) //分享自己的稿件
	s.relationSub = databus.New(s.c.RelationSub)
	s.newUpSub = databus.New(s.c.NewUpSub) //新投稿up主

	//单个稿件计数
	s.statViewSub = databus.New(s.c.StatViewSub)
	s.statLikeSub = databus.New(s.c.StatLikeSub)
	s.statShareSub = databus.New(s.c.StatShareSub) //计数分享
	s.statCoinSub = databus.New(s.c.StatCoinSub)
	s.statFavSub = databus.New(s.c.StatFavSub)
	s.statReplySub = databus.New(s.c.StatReplySub)
	s.statDMSub = databus.New(s.c.StatDMSub)

	s.taskSubQueue = make([]chan *databus.Message, s.databusQueueLen) //设置水印、观看创作学院视频、参加激励计划、开通粉丝勋章
	s.shareSubQueue = make([]chan *model.ShareMsg, s.databusQueueLen) //分享自己的稿件
	s.relationQueue = make([]chan *model.Relation, s.databusQueueLen) //用户关注队列
	s.followerQueue = make([]chan *model.Stat, s.databusQueueLen)     //粉丝数队列
	s.newUpQueue = make([]chan *model.Up, s.databusQueueLen)          //新投稿up主
	s.oldUpQueue = make([]chan *model.Up, s.databusQueueLen)          //进阶任务视频投稿超过5个
	s.mobileUpQueue = make([]chan *model.Up, s.databusQueueLen)       //进阶任务手机投稿

	//单个稿件计数
	s.statViewSubQueue = make([]chan *model.StatView, s.statViewQueueLen)
	s.statLikeSubQueue = make([]chan *model.StatLike, s.statLikeQueueLen)
	s.statShareSubQueue = make([]chan *model.StatShare, s.databusQueueLen)
	s.statCoinSubQueue = make([]chan *model.StatCoin, s.databusQueueLen)
	s.statFavSubQueue = make([]chan *model.StatFav, s.databusQueueLen)
	s.statReplySubQueue = make([]chan *model.StatReply, s.databusQueueLen)
	s.statDMSubQueue = make([]chan *model.StatDM, s.databusQueueLen)
}

// TaskClose close task sub.
func (s *Service) TaskClose() {
	s.taskSub.Close()     //水印、激励、观看创作学院、开通粉丝勋章
	s.shareSub.Close()    //个人稿件分享
	s.relationSub.Close() //关注哔哩哔哩创组中心，新手和进阶粉丝数
	s.newUpSub.Close()    //投第一个稿件

	//计数
	s.statViewSub.Close()
	s.statLikeSub.Close()
	s.statShareSub.Close()
	s.statCoinSub.Close()
	s.statFavSub.Close()
	s.statReplySub.Close()
	s.statDMSub.Close()
}

func (s *Service) taskConsume() {
	s.loadTasks()       //定时缓存所有任务
	s.loadGiftRewards() //定时缓存所有奖励
	go s.loadProc()

	//非实时任务状态变更
	s.initTaskQueue()
	go s.commitTask()

	//过期任务通知
	if s.c.Task.SwitchMsgNotify {
		mids, _ := xstr.SplitInts(s.c.Task.TestNotifyMids)
		for _, mid := range mids {
			s.testNotifyMids[mid] = struct{}{} //test mids
		}
		s.initTaskNotifyQueue()
		go s.expireTaskNotify()

		//奖励领取通知
		s.initRewardNotifyQueue()
		go s.rewardReceiveNotify()
	}

	//实时任务状态变更
	if s.c.Task.SwitchHighQPS { //消息qps 较高的消费
		s.initStatViewQueue()
		s.initStatLikeQueue()
		s.wg.Add(2)
		go s.statView() //1
		go s.statLike() //2
	}

	if s.c.Task.SwitchDatabus { //消息qps 较少的消费
		s.wg.Add(9)
		s.initDatabusQueue()
		go s.task()      //1
		go s.share()     //2
		go s.relation()  //3
		go s.statShare() //4
		go s.statCoin()  //5
		go s.statFav()   //6
		go s.statReply() //7
		go s.statDM()    //8
		go s.newUp()     //9
	}
}

// Ping service
func (s *Service) Ping(c context.Context) (err error) {
	return
}

func (s *Service) monitorConsume() {
	if s.c.Env != env.DeployEnvProd {
		return
	}
	var arcNotifyMo, arcmo int64
	for {
		time.Sleep(time.Minute * 1)
		if s.arcNotifyMo-arcNotifyMo == 0 {
			s.monitor.Send(context.TODO(), s.c.Monitor.UserName, "creative-job did not consume within a minute, moni url"+s.c.Monitor.Moni)
		}
		arcNotifyMo = s.arcNotifyMo
		if s.arcMo-arcmo == 0 {
			s.monitor.Send(context.TODO(), s.c.Monitor.UserName, "creative-job did not consume within a minute, moni url"+s.c.Monitor.Moni)
		}
		arcmo = s.arcMo
	}
}

// Close sub.
func (s *Service) Close() {
	s.arcSub.Close()
	s.arcNotifySub.Close()
	s.upPub.Close()
	close(s.midsChan)
	s.TaskClose() //task databus close
	s.wg.Wait()
}
