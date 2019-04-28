package service

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"go-common/app/admin/main/member/conf"
	"go-common/app/admin/main/member/dao"
	"go-common/app/admin/main/member/model"
	"go-common/app/admin/main/member/service/block"
	acccrypto "go-common/app/interface/main/account/service/realname/crypto"
	account "go-common/app/service/main/account/api"
	coinrpc "go-common/app/service/main/coin/api/gorpc"
	rpcfigure "go-common/app/service/main/figure/rpc/client"
	memberrpc "go-common/app/service/main/member/api/gorpc"
	"go-common/app/service/main/member/service/crypto"
	rpcrelation "go-common/app/service/main/relation/rpc/client"
	rpcspy "go-common/app/service/main/spy/rpc/client"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
	"go-common/library/stat/prom"
	xtime "go-common/library/time"

	"github.com/robfig/cron"
)

// Service struct
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	block          *block.Service
	auditHandlers  map[string]auditHandler
	coinRPC        *coinrpc.Service
	memberRPC      *memberrpc.Service
	spyRPC         *rpcspy.Service
	figureRPC      *rpcfigure.Service
	accountClient  account.AccountClient
	cron           *cron.Cron
	relationRPC    *rpcrelation.Service
	realnameCrypto *crypto.Realname
	mainCryptor    *acccrypto.Main
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:              c,
		dao:            dao.New(c),
		coinRPC:        coinrpc.New(c.RPCClient.Coin),
		memberRPC:      memberrpc.New(c.RPCClient.Member),
		figureRPC:      rpcfigure.New(c.RPCClient.Figure),
		spyRPC:         rpcspy.New(c.RPCClient.Spy),
		relationRPC:    rpcrelation.New(c.RPCClient.Relation),
		auditHandlers:  make(map[string]auditHandler),
		cron:           cron.New(),
		realnameCrypto: crypto.NewRealname(string(c.Realname.RsaPub), string(c.Realname.RsaPriv)),
		mainCryptor:    acccrypto.NewMain(string(c.Realname.RsaPub), string(c.Realname.RsaPriv)),
	}
	var err error
	if s.accountClient, err = account.NewClient(c.RPCClient.Account); err != nil {
		panic(err)
	}
	s.block = block.New(c, s.dao.BlockImpl(), s.spyRPC, s.figureRPC, s.accountClient, databus.New(c.AccountNotify))
	s.initAuditHandler()
	s.initCron()
	s.cron.Start()
	return s
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
	s.block.Close()
}

func (s *Service) initCron() {
	s.cron.AddFunc("0 */5 * * * *", func() { s.notifyAudit(context.Background()) })                                  // 用于发送审核数据给目标用户
	s.cron.AddFunc("0 */5 * * * *", func() { s.promAuditTotal(context.Background()) })                               // 用于上报审核数据给promethues
	s.cron.AddFunc("0 */1 * * * *", func() { s.cacheRecentRealnameImage(context.Background()) })                     // 用于缓存实名认证的图片数据
	s.cron.AddFunc("0 */2 * * * *", func() { s.faceCheckproc(context.Background(), -10*time.Minute, "two minute") }) // 用于AI头像审核：首次：每隔2分钟审核10分钟内的头像
	s.cron.AddFunc("0 */60 * * * *", func() { s.faceCheckproc(context.Background(), -6*time.Hour, "per hour") })     // 用于AI头像审核：重新审核：每小时重新审核一下6小时内的头像
	s.cron.AddFunc("0 */5 * * * *", func() { s.faceAutoPassproc(context.Background()) })                             // 头像自动审核：每隔5分钟检查一次超过48小时未处理的头像并自动通过
}

// notifyAudit
func (s *Service) notifyAudit(ctx context.Context) {
	now := time.Now()
	log.Info("start notify audit at: %+v", now)
	locked, err := s.dao.TryLockReviewNotify(ctx, now)
	if err != nil {
		log.Error("Failed to lock review notify at: %+v: %+v", now, err)
		return
	}
	if !locked {
		log.Warn("Already locked by other instance at: %+v", now)
		return
	}
	stime := now.Add(-time.Hour * 24 * 7) // 只计算 7 天内的数据
	// 绝对锁上了
	faceNotify := func() error {
		total, firstAt, err := s.faceAuditNotifyContent(ctx, stime)
		if err != nil {
			log.Error("Failed to fetch face audit notify content: %+v", err)
			return err
		}
		log.Info("faceAuditNotifyContent success: total(%v),firstAt(%v)", total, firstAt)
		title := fmt.Sprintf("头像审核提醒；消息时间：%s", now.Format("2006-01-02 15:04:05"))
		firstAtStr := "null"
		if firstAt != nil {
			firstAtStr = firstAt.Format("2006-01-02 15:04:05")
		}
		content := fmt.Sprintf(
			"头像审核提醒；消息时间：%s\n头像审核积压：%d 条；最早进审时间：%s",
			now.Format("2006-01-02 15:04:05"),
			total,
			firstAtStr,
		)
		return s.dao.MerakNotify(ctx, title, content)
	}
	if err := faceNotify(); err != nil {
		log.Error("Failed to notify face review stat: %+v", err)
	}

	monitorNotify := func() error {
		total, firstAt, err := s.monitorAuditNotifyContent(ctx, stime)
		if err != nil {
			log.Error("Failed to fetch monitor audit notify content: %+v", err)
			return err
		}
		log.Info("monitorAuditNotifyContent success: total(%v),firstAt(%v)", total, firstAt)
		title := fmt.Sprintf("用户信息监控提醒；消息时间：%s", now.Format("2006-01-02 15:04:05"))
		firstAtStr := "null"
		if firstAt != nil {
			firstAtStr = firstAt.Format("2006-01-02 15:04:05")
		}
		content := fmt.Sprintf(
			"用户信息监控提醒；消息时间：%s\n用户信息监控积压：%d 条；最早进审时间：%s",
			now.Format("2006-01-02 15:04:05"),
			total,
			firstAtStr,
		)
		return s.dao.MerakNotify(ctx, title, content)
	}
	if err := monitorNotify(); err != nil {
		log.Error("Failed to notify monitor review stat: %+v", err)
	}

	// 实名认证待审核通知
	realnameNotify := func() error {
		total, firstAt, err := s.realnameAuditNotifyContent(ctx, stime)
		if err != nil {
			log.Error("Failed to fetch realname audit notify content: %+v", err)
			return err
		}
		log.Info("realnameAuditNotifyContent success: total(%v),firstAt(%v)", total, firstAt)
		title := fmt.Sprintf("实名认证审核提醒；消息时间：%s", now.Format("2006-01-02 15:04:05"))
		firstAtStr := "null"
		if firstAt != nil {
			firstAtStr = firstAt.Format("2006-01-02 15:04:05")
		}
		content := fmt.Sprintf(
			"实名认证审核提醒；消息时间：%s\n实名认证审核积压：%d 条；最早进审时间：%s",
			now.Format("2006-01-02 15:04:05"),
			total,
			firstAtStr,
		)
		return s.dao.MerakNotify(ctx, title, content)
	}
	if err := realnameNotify(); err != nil {
		log.Error("Failed to notify realname list stat: %+v", err)
	}

	log.Info("end notify audit at: %+v", now)
}

// promAuditTotal
func (s *Service) promAuditTotal(ctx context.Context) {
	stime := time.Now().Add(-time.Hour * 24 * 7) // 只计算 7 天内的数据
	log.Info("promAuditTotal start %+v", time.Now())
	faceAudit := func() {
		faceTotal, _, err := s.faceAuditNotifyContent(ctx, stime)
		if err != nil {
			log.Error("Failed to fetch face audit notify content: %+v", err)
			return
		}
		prom.BusinessInfoCount.State("faceAudit-needAudit", int64(faceTotal))
	}

	monitorAudit := func() {
		monitorTotal, _, err := s.monitorAuditNotifyContent(ctx, stime)
		if err != nil {
			log.Error("Failed to fetch monitor audit notify content: %+v", err)
			return
		}
		prom.BusinessInfoCount.State("monitorAudit-needAudit", int64(monitorTotal))
	}

	realnameAudit := func() {
		realnameTotal, _, err := s.realnameAuditNotifyContent(ctx, stime)
		if err != nil {
			log.Error("Failed to fetch realname audit notify content: %+v", err)
			return
		}
		prom.BusinessInfoCount.State("realnameAudit-needAudit", int64(realnameTotal))
	}

	faceAudit()
	monitorAudit()
	realnameAudit()
	log.Info("promAuditTotal end %+v", time.Now())

}

func (s *Service) faceAuditNotifyContent(ctx context.Context, stime time.Time) (int, *time.Time, error) {
	arg := &model.ArgReviewList{
		State:     []int8{0},
		IsMonitor: false,
		Property:  []int8{model.ReviewPropertyFace},
		IsDesc:    false,
		Pn:        1,
		Ps:        1,
		STime:     xtime.Time(stime.Unix()),
		ForceDB:   false,
	}
	reviews, total, err := s.Reviews(ctx, arg)
	if err != nil {
		return 0, nil, err
	}
	if len(reviews) <= 0 {
		return 0, nil, nil
	}
	firstAt := reviews[0].CTime.Time()
	return total, &firstAt, nil
}

func (s *Service) monitorAuditNotifyContent(ctx context.Context, stime time.Time) (int, *time.Time, error) {
	arg := &model.ArgReviewList{
		State:     []int8{0},
		IsMonitor: true,
		IsDesc:    false,
		Pn:        1,
		Ps:        1,
		STime:     xtime.Time(stime.Unix()),
		ForceDB:   false,
	}
	reviews, total, err := s.Reviews(ctx, arg)
	if err != nil {
		return 0, nil, err
	}
	if len(reviews) <= 0 {
		return 0, nil, nil
	}
	firstAt := reviews[0].CTime.Time()
	return total, &firstAt, nil
}

func (s *Service) realnameAuditNotifyContent(ctx context.Context, stime time.Time) (int, *time.Time, error) {
	arg := &model.ArgRealnameList{
		Channel: "main", //main : 主站  alipay : 支付宝
		TSFrom:  stime.Unix(),
		State:   model.RealnameApplyStatePending,
		IsDesc:  false,
		PN:      1,
		PS:      1,
	}
	mainList, total, err := s.realnameMainList(ctx, arg)
	if err != nil {
		return 0, nil, err
	}
	if len(mainList) <= 0 {
		return 0, nil, nil
	}
	firstAt := time.Unix(mainList[0].CreateTS, 0)
	return total, &firstAt, nil
}

func (s *Service) faceAutoPassproc(ctx context.Context) {
	now := time.Now()
	log.Info("faceAutoPassproc start %+v", now)
	etime := now.AddDate(0, 0, -2)
	if err := s.faceAutoPass(ctx, etime); err != nil {
		log.Error("Failed to face auto pass, error: %+v", err)
	}
}

func (s *Service) faceAutoPass(ctx context.Context, etime time.Time) error {
	property := []int{model.ReviewPropertyFace}
	state := []int{model.ReviewStateWait, model.ReviewStateQueuing}
	result, err := s.dao.SearchUserPropertyReview(ctx, 0, property,
		state, false, false, "", "", etime.Format("2006-01-02 15:04:05"), 1, 100)
	if err != nil {
		return err
	}
	ids := result.IDs()
	if len(ids) == 0 {
		log.Info("face auto pass empty result list, end time: %v", etime)
		return nil
	}
	if err = s.dao.FaceAutoPass(ctx, ids, xtime.Time(etime.Unix())); err != nil {
		return err
	}
	return nil
}

func (s *Service) faceCheckproc(ctx context.Context, duration time.Duration, tag string) {
	now := time.Now()
	stime := now.Add(duration).Unix()
	etime := now.Unix()
	log.Info("faceCheckproc:%v start %+v", tag, now)
	if err := s.faceAuditAI(ctx, stime, etime); err != nil {
		log.Error("Failed to check face, error: %+v", err)
	}
}

func (s *Service) faceAuditAI(ctx context.Context, stime, etime int64) error {
	rws, err := s.dao.QueuingFaceReviewsByTime(ctx, xtime.Time(stime), xtime.Time(etime))
	if err != nil {
		log.Warn("Failed to get recent user_property_review image: %+v", err)
		return err
	}
	for _, rw := range rws {
		fcr, err := s.faceCheckRes(ctx, path.Base(rw.New))
		if err != nil {
			log.Error("Failed to get face check res, rw: %+v, error: %+v", rw, err)
			continue
		}
		state := int8(model.ReviewStateWait)
		if fcr.Valid() {
			state = model.ReviewStatePass
		}
		remark := fmt.Sprintf("AI: %s", fcr.String())
		if err = s.dao.AuditQueuingFace(ctx, rw.ID, remark, state); err != nil {
			log.Error("Failed to audit queuing face, rw: %+v, error: %+v", rw, err)
			continue
		}
		log.Info("face check success, rw: %+v", rw)
	}
	log.Info("faceCheckproc end")
	return nil
}

func (s *Service) faceCheckRes(ctx context.Context, fileName string) (*model.FaceCheckRes, error) {
	res, err := s.dao.SearchFaceCheckRes(ctx, fileName)
	if err != nil {
		return nil, err
	}
	if len(res.Result) == 0 {
		return nil, ecode.NothingFound
	}
	userLog := res.Result[0]
	fcr, err := parseFaceCheckRes(userLog.Extra)
	if err != nil {
		log.Error("Failed to parse faceCheckRes, userLog: %+v error: %+v", userLog, err)
		return nil, err
	}
	return fcr, nil
}

func parseFaceCheckRes(in string) (*model.FaceCheckRes, error) {
	res := &model.FaceCheckRes{}
	err := json.Unmarshal([]byte(in), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// BlockImpl is
func (s *Service) BlockImpl() *block.Service {
	return s.block
}

func (s *Service) cacheRecentRealnameImage(ctx context.Context) {
	images, err := s.dao.RecentRealnameApplyImg(ctx, time.Minute*2)
	if err != nil {
		log.Warn("Failed to get recent realname apply image: %+v", err)
		return
	}
	for _, image := range images {
		data, _ := s.dao.GetRealnameImageCache(ctx, image.IMGData)
		if len(data) > 0 {
			log.Info("This image has already been cached: %s", image.IMGData)
			continue
		}
		data, err := s.FetchRealnameImage(ctx, asIMGToken(image.IMGData))
		if err != nil {
			log.Warn("Failed to fetch realname image to cache: %s: %+v", image.IMGData, err)
			continue
		}
		if err := s.dao.SetRealnameImageCache(ctx, image.IMGData, data); err != nil {
			log.Warn("Failed to set realname image cache: %s: %+v", image.IMGData, err)
			continue
		}
		log.Info("Succeeded to cache realname image: %s", image.IMGData)
	}
}
