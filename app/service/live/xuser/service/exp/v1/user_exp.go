package v1

import (
	"context"
	"github.com/pkg/errors"
	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/conf"
	"go-common/app/service/live/xuser/dao/exp"
	XanchorDao "go-common/app/service/live/xuser/dao/xanchor"
	expm "go-common/app/service/live/xuser/model/exp"
	"go-common/library/cache"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/stat/prom"
	"go-common/library/sync/pipeline/fanout"
)

// UserExpService struct
type UserExpService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao         *exp.Dao
	addExpCache *cache.Cache
	dbUpdater   *fanout.Fanout
	xuserDao    *XanchorDao.Dao
}

const (
	_errorServiceLogPrefix = "xuser.exp.service"
	_promCacheMissed       = "xuser_exp_mc:用户经验cache miss"
	_promCacheHitAll       = "xuser_exp_mc:用户经验cache全部命中"
	_retryAddExpTimes      = 3
)

// NewUserExpService init
func NewUserExpService(c *conf.Config) (s *UserExpService) {
	s = &UserExpService{
		conf:        c,
		dao:         exp.NewExpDao(c),
		addExpCache: cache.New(1, 10240),
		dbUpdater:   fanout.New("jibamao", fanout.Worker(1), fanout.Buffer(10240)),
		xuserDao:    XanchorDao.New(c),
	}
	return s
}

// Close close exp service
func (s *UserExpService) Close() {
	s.dao.Close()
}

func (s *UserExpService) adaptResultFromMemCache(expList map[int64]*expm.LevelInfo) (resp *v1pb.GetUserExpResp) {
	resp = &v1pb.GetUserExpResp{}
	resp.Data = make(map[int64]*v1pb.LevelInfo)
	for _, v := range expList {
		resp.Data[v.UID] = &v1pb.LevelInfo{}
		resp.Data[v.UID] = s.adaptAPIModel(v)
	}
	return
}

func (s *UserExpService) adaptResultFromDBMemCache(expFromMc map[int64]*expm.LevelInfo, expFromDB map[int64]*expm.LevelInfo) (resp *v1pb.GetUserExpResp) {
	resp = &v1pb.GetUserExpResp{}
	resp.Data = make(map[int64]*v1pb.LevelInfo)
	if len(expFromMc) <= 0 && len(expFromDB) <= 0 {
		return

	} else if len(expFromMc) <= 0 {
		for _, v := range expFromDB {
			resp.Data[v.UID] = &v1pb.LevelInfo{}
			resp.Data[v.UID] = s.adaptAPIModel(v)
		}
	} else if len(expFromDB) <= 0 {
		for _, v := range expFromMc {
			resp.Data[v.UID] = &v1pb.LevelInfo{}
			resp.Data[v.UID] = s.adaptAPIModel(v)
		}
	} else {
		for _, v := range expFromMc {
			resp.Data[v.UID] = &v1pb.LevelInfo{}
			resp.Data[v.UID] = s.adaptAPIModel(v)
		}
		for _, v := range expFromDB {
			resp.Data[v.UID] = &v1pb.LevelInfo{}
			resp.Data[v.UID] = s.adaptAPIModel(v)
		}
	}
	return
}

// 直播用户经验gRPC接口

// GetUserExp ...
// 获取用户经验与等级信息,支持批量
func (s *UserExpService) GetUserExp(ctx context.Context, req *v1pb.GetUserExpReq) (resp *v1pb.GetUserExpResp, err error) {
	reqStartTime := s.RecordTimeCost()
	// s.RecordTimeCostLog(reqStartTime, "G0|获取经验入口耗时")
	resp = &v1pb.GetUserExpResp{}
	resp.Data = make(map[int64]*v1pb.LevelInfo)
	cacheHealth := true
	expResultFromMc, missedUIDs, err := s.dao.GetExpFromMemCache(ctx, req.Uids)

	// 回源db原则,仅在get成功且miss时回源db!!!
	if err != nil {
		reqAfterQueryMCTime := s.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_QUERY_CACHE_FAIL+"|GetUserExp|查询缓存失败,暂不回源db,接口返回err|%dms", err, reqAfterQueryMCTime-reqStartTime)
		err = errors.WithMessage(ecode.XUserExpGetExpMcFail, "获取用户经验缓存失败")
		cacheHealth = false
		return
	} else if len(missedUIDs) == 0 {
		resp = s.adaptResultFromMemCache(expResultFromMc)
		exp.PromCacheHit(_promCacheHitAll)
		return
	}
	exp.PromCacheMiss(_promCacheMissed)
	prom.CacheMiss.Incr("expList_mc")
	resultDB, err := s.dao.MultiExp(ctx, missedUIDs)
	if err != nil {
		reqAfterQueryDBTime := s.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_GETEXP_DB_FAIL+"|GetUserExp|获取用户经验回源DB失败error|(%v),missedUIDs(%v)|耗时:%dms", err, missedUIDs, reqAfterQueryDBTime-reqStartTime)
		err = errors.WithMessage(ecode.XUserExpGetExpDBFail, "获取用户经验回源DB失败")
		return
	}
	expResultFromDB := s.FormatLevel(resultDB)
	// 写入缓存
	if cacheHealth {
		s.asyncSetExpCache(ctx, expResultFromDB)
	}
	resp = s.adaptResultFromDBMemCache(expResultFromMc, expResultFromDB)
	return
}

// AddUserExp ...
// 增加用户经验,支持批量
func (s *UserExpService) AddUserExp(ctx context.Context, req *v1pb.AddUserExpReq) (resp *v1pb.AddUserExpResp, err error) {
	resp = &v1pb.AddUserExpResp{}
	if req == nil {
		err = ecode.XUserAddUserExpParamsEmptyError
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ADD_UEXP_PARA_NIL+"|AddUserExp|s.dao.Exp|error(%v)|入参校验错误", err)
		err = errors.WithMessage(ecode.ResourceParamErr, "添加用户经验入参异常")
		return
	}
	userInfo := req.UserInfo
	allow, _ := s.checkReqBiz(userInfo.ReqBiz)
	if !allow {
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ADD_UEXP_REQ_BIZ_FAIL+"|AddUserExp|error(%v),reqBiz(%d)is not register yet", err, userInfo.ReqBiz)
		err = ecode.XUserAddUserExpReqBizNotAllow
		return
	}
	typeAllow, _ := s.checkAddType(userInfo.Type)
	if !typeAllow {
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ADD_UEXP_TYPE_FAIL+"|AddUserExp|error(%v),addType(%d)is not allow yet", err, userInfo.Type)
		err = ecode.XUserAddUserExpTypeNotAllow
		return
	}
	numAllow, _ := s.checkAddNum(userInfo.Num)
	if !numAllow {
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ADD_UEXP_NUM_FAIL+"|AddUserExp|error(%v),addNum(%d)is not allow yet", err, userInfo.Num)
		err = ecode.XUserAddUserExpNumNotAllow
	}
	s.doAddExp(ctx, userInfo)

	return
}

func (s *UserExpService) doAddExp(ctx context.Context, req *v1pb.UserExpChunk) (err error) {
	switch req.Type {
	case _addUserExpType:
		{
			return s.asyncUpdateDBUexp(ctx, req)
		}
	case _addAnchorExpType:
		{
			return s.asyncUpdateDBRexp(ctx, req)
		}
	}
	return
}

func (s *UserExpService) adaptAPIModel(dbModel *expm.LevelInfo) (apiModel *v1pb.LevelInfo) {
	apiModel = &v1pb.LevelInfo{}
	apiModel.UserLevel = &v1pb.UserLevelInfo{}
	apiModel.AnchorLevel = &v1pb.AnchorLevelInfo{}

	apiModel.Uid = dbModel.UID

	apiModel.UserLevel.UserExpNextRight = dbModel.UserLevel.UserExpNextRight
	apiModel.UserLevel.IsLevelTop = dbModel.UserLevel.IsLevelTop
	apiModel.UserLevel.UserExpNextLeft = dbModel.UserLevel.UserExpNextLeft
	apiModel.UserLevel.Level = dbModel.UserLevel.Level
	apiModel.UserLevel.NextLevel = dbModel.UserLevel.NextLevel
	apiModel.UserLevel.UserExpLeft = dbModel.UserLevel.UserExpLeft
	apiModel.UserLevel.UserExpRight = dbModel.UserLevel.UserExpRight
	apiModel.UserLevel.UserExp = dbModel.UserLevel.UserExp
	apiModel.UserLevel.UserExpNextLevel = dbModel.UserLevel.UserExpNextLevel
	apiModel.UserLevel.Color = dbModel.UserLevel.Color

	apiModel.AnchorLevel.UserExpNextRight = dbModel.AnchorLevel.UserExpNextRight
	apiModel.AnchorLevel.IsLevelTop = dbModel.AnchorLevel.IsLevelTop
	apiModel.AnchorLevel.UserExpNextLeft = dbModel.AnchorLevel.UserExpNextLeft
	apiModel.AnchorLevel.Level = dbModel.AnchorLevel.Level
	apiModel.AnchorLevel.NextLevel = dbModel.AnchorLevel.NextLevel
	apiModel.AnchorLevel.UserExpLeft = dbModel.AnchorLevel.UserExpLeft
	apiModel.AnchorLevel.UserExpRight = dbModel.AnchorLevel.UserExpRight
	apiModel.AnchorLevel.UserExp = dbModel.AnchorLevel.UserExp
	apiModel.AnchorLevel.UserExpNextLevel = dbModel.AnchorLevel.UserExpNextLevel
	apiModel.AnchorLevel.Color = dbModel.AnchorLevel.Color
	apiModel.AnchorLevel.AnchorScore = dbModel.AnchorLevel.AnchorScore

	return
}
