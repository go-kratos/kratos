package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/live/resource/lrucache"
	"go-common/app/service/live/resource/model"
	"go-common/library/log"
	"strconv"
	"strings"
	"time"

	v1pb "go-common/app/service/live/resource/api/grpc/v1"
	v1hpb "go-common/app/service/live/resource/api/http/v1"
	"go-common/app/service/live/resource/conf"
	"go-common/app/service/live/resource/dao"
	"go-common/library/ecode"
)

// TitansService struct
type TitansService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
	//lruCache
	treeCache   *lrucache.SyncCache
	titansCache *lrucache.SyncCache
}

// NewTitansService init
func NewTitansService(c *conf.Config) (s *TitansService) {
	length := 10
	treeCache := lrucache.NewSyncCache(length, 100, 600)
	titansCache := lrucache.NewSyncCache(length, 100, 60)

	s = &TitansService{
		conf:        c,
		dao:         dao.New(c),
		treeCache:   treeCache,
		titansCache: titansCache,
	}
	//启动同步流程
	go s.SyncConfig()
	return s
}

const (
	_titansLruCache = "titans_cache_treeId:%d"
)

// GetConfigByKeyword implementation
// http获取team下某个keyword的配置 `internal:"true"`
func (s *TitansService) GetConfigByKeyword(ctx context.Context, req *v1pb.GetConfigReq) (resp *v1pb.GetConfigResp, err error) {
	resp = &v1pb.GetConfigResp{}
	team := req.GetTeam()
	keyword := req.GetKeyword()

	if 0 == req.GetId() {
		err = s.dao.CheckParams(ctx, team, keyword)
		if err != nil {
			err = ecode.Error(ecode.ParamInvalid, "必要参数不准为空")
			return
		}
	}

	res, err := s.dao.SelectByTeamIndex(ctx, team, keyword, req.GetId())
	if err != nil {
		//errors.Wrap(err, "系统错误")
		return
	}
	resp = &v1pb.GetConfigResp{
		Team:    res.Team,
		Keyword: res.Keyword,
		// 配置值
		Value: res.Value,
		// 配置解释
		Name: res.Name,
		// 创建时间
		Ctime: res.Ctime,
		//最近更新时间
		Mtime:  res.Ctime,
		Status: res.Status,
		Id:     res.Id,
	}
	return
}

// SetConfigByKeyword implementation
// http设置team下某个keyword配置 `internal:"true"`
func (s *TitansService) SetConfigByKeyword(ctx context.Context, req *v1pb.SetConfigReq) (resp *v1pb.SetConfigResp, err error) {
	resp = &v1pb.SetConfigResp{}
	if "" == req.GetKeyword() {
		err = ecode.Error(ecode.ParamInvalid, "必要参数不准为空")
		return
	}
	id, count, err := s.dao.InsertRecord(ctx, req.GetTeam(), req.GetKeyword(), req.GetValue(), req.GetName(), req.Status, req.GetId())
	if nil != err {
		err = ecode.Error(ecode.ServerErr, "系统错误")
		return
	}
	if 0 != count {
		err = ecode.Error(
			11,
			"分组"+strconv.Itoa(int(req.GetTeam()))+"内已存在"+req.GetKeyword()+"的配置",
		)
		return
	}
	resp.Id = id
	return
}

// GetConfigsByParams implementation
// http管理后台根据条件获取配置 `internal:"true"`
func (s *TitansService) GetConfigsByParams(ctx context.Context, req *v1pb.ParamsConfigReq) (resp *v1pb.ParamsConfigResp, err error) {
	resp = &v1pb.ParamsConfigResp{}
	list, count, err := s.dao.SelectByParams(ctx, req.GetId(), req.GetTeam(), req.GetKeyword(), req.GetName(), req.GetStatus(), req.GetPage(), req.GetPageSize())
	if nil != err {
		err = ecode.Error(ecode.ServerErr, "系统错误")
		return
	}
	resp.List = []*v1pb.List{}
	for _, v := range list {
		resp.List = append(resp.List, &v1pb.List{
			Id:      v.Id,
			Team:    v.Team,
			Keyword: v.Keyword,
			Name:    v.Name,
			Value:   v.Value,
			Ctime:   v.Ctime,
			Mtime:   v.Mtime,
			Status:  v.Status,
		})
	}
	resp.TotalNum = count
	return
}

// GetConfigsByLikes implementation
// grpc获取多个team或索引的的全部配置 `internal:"true"`
func (s *TitansService) GetConfigsByLikes(ctx context.Context, req *v1pb.LikesConfigReq) (resp *v1pb.LikesConfigResp, err error) {
	resp = &v1pb.LikesConfigResp{}
	params := req.GetParams()
	teams := make([]int64, 0)
	teamKeys := make([]*dao.TeamKeyword, 0)
	for _, v := range params {
		strArr := strings.Split(v, ".")
		team, _ := strconv.ParseInt(strArr[0], 10, 64)
		if team != 0 && len(strArr) == 1 {
			teams = append(teams, team)
		}
		if len(strArr) == 2 && team != 0 {
			teamKeys = append(teamKeys, &dao.TeamKeyword{Team: team, Keyword: strArr[1]})
		}
	}

	items, err := s.dao.SelectByLikes(ctx, teams, teamKeys)
	if err != nil {
		return
	}

	mapParent := make(map[int64]*v1pb.Child)
	for _, v := range items {
		parentSet(mapParent, v.Team, v.Keyword, v.Value)
	}
	resp.List = mapParent
	return
}

// parentSet format数据
func parentSet(parent map[int64]*v1pb.Child, index int64, keyword string, value string) {
	_, ok := parent[index]
	if !ok {
		parent[index] = &v1pb.Child{}
		parent[index].Keys = make(map[string]string)
	}
	parent[index].Keys[keyword] = value
}

// parentHSet format数据
func parentHSet(parent map[int64]*v1hpb.MChild, index int64, keyword string, value string) {
	_, ok := parent[index]
	if !ok {
		parent[index] = &v1hpb.MChild{}
		parent[index].Keys = make(map[string]string)
	}
	parent[index].Keys[keyword] = value
}

// GetMultiConfigs implementation
// http获取配置 请求参数逗号隔开的字符串 返回`internal:"true"`
func (s *TitansService) GetMultiConfigs(ctx context.Context, req *v1hpb.MultiConfigReq) (resp *v1hpb.MultiConfigResp, err error) {
	resp = &v1hpb.MultiConfigResp{}
	params := strings.Split(req.GetValues(), ",")
	teams := make([]int64, 0)
	teamKeys := make([]*dao.TeamKeyword, 0)
	for _, v := range params {
		strArr := strings.Split(v, ".")
		team, _ := strconv.ParseInt(strArr[0], 10, 64)
		if team != 0 && len(strArr) == 1 {
			teams = append(teams, team)
		}
		if len(strArr) == 2 && team != 0 {
			teamKeys = append(teamKeys, &dao.TeamKeyword{Team: team, Keyword: strArr[1]})
		}
	}

	items, err := s.dao.SelectByLikes(ctx, teams, teamKeys)
	if err != nil {
		return
	}

	mapParent := make(map[int64]*v1hpb.MChild)
	for _, v := range items {
		parentHSet(mapParent, v.Team, v.Keyword, v.Value)
	}
	resp.List = mapParent
	return
}

// GetServiceConfig implementation
// http获取服务tree_id对应的配置 `internal:"true"`
func (s *TitansService) GetServiceConfig(ctx context.Context, req *v1hpb.ServiceConfigReq) (resp *v1hpb.ServiceConfigResp, err error) {
	resp = &v1hpb.ServiceConfigResp{}
	treeId := req.GetTreeId()
	if 0 == treeId {
		err = ecode.Error(ecode.InvalidParam, "tree_id 为空")
		return
	}
	cacheKey := fmt.Sprintf(_titansLruCache, treeId)
	// 读取lruCache
	cacheValue, ok := s.titansCache.Get(cacheKey)
	if ok {
		resp.List = cacheValue.(map[string]string)
		return
	}
	value, err := s.dao.GetServiceConfig(ctx, treeId)
	if nil != err {
		err = ecode.Error(ecode.ServerErr, "内部错误")
		return
	}
	resp.List = value
	// 存lruCache
	s.titansCache.Put(cacheKey, value)
	return
}

// SetServiceConfig implementation
// http插入服务配置 `method:"POST", internal:"true"`
func (s *TitansService) SetServiceConfig(ctx context.Context, req *v1hpb.SetReq) (resp *v1hpb.SetResp, err error) {
	resp = &v1hpb.SetResp{}
	if 0 == req.GetTreeId() || "" == req.GetTreeName() || "" == req.GetKeyword() {
		err = ecode.Error(ecode.ParamInvalid, "服务的tree_name, tree_id, keyword必传")
		return
	}

	if len(req.GetKeyword()) > 16 {
		err = ecode.Error(ecode.ParamInvalid, "keyword长度不能超过16")
		return
	}
	if len(req.GetValue()) > 4096 {
		err = ecode.Error(ecode.ParamInvalid, "配置内容长度不能超过4096")
		return
	}

	if "" != req.GetValue() {
		if !json.Valid([]byte(req.GetValue())) {
			err = ecode.Error(ecode.InvalidParam, "配置的内容必须为json")
			return
		}
	}
	id, err := s.dao.InsertServiceConfig(
		ctx, req.GetId(),
		req.GetTreeName(),
		req.GetTreePath(),
		req.GetTreeId(),
		req.GetService(),
		req.GetKeyword(),
		req.GetTemplate(),
		req.GetName(),
		req.GetValue(),
		req.GetStatus())
	if nil != err {
		err = ecode.Error(ecode.ServerErr, "内部错误")
		return
	}
	if -1 == id {
		err = ecode.Error(ecode.InvalidParam, "同一个tree_id下keyword配置重复，请确认修改")
		return
	}
	resp.Id = id
	// 编辑操作后，清一下lrucache
	s.titansCache.Delete(fmt.Sprintf(_titansLruCache, req.GetTreeId()))
	return
}

// GetServiceConfigList implementation
// http管理后台获取服务级配置 `internal:"true"`
func (s *TitansService) GetServiceConfigList(ctx context.Context, req *v1hpb.ServiceListReq) (resp *v1hpb.ServiceListResp, err error) {
	resp = &v1hpb.ServiceListResp{}
	if "" == req.GetTreeName() {
		err = ecode.Error(ecode.InvalidParam, "服务树根名不能为空")
		return
	}
	if req.GetTreeId() == 0 && (req.GetName() != "" || req.GetService() != "") {
		err = ecode.Error(ecode.InvalidParam, "通过描述名称查询时，tree_id不准为空")
		return
	}
	page := req.GetPage()
	pageSize := req.GetPageSize()
	if pageSize == 0 {
		pageSize = 30
	}
	if page == 0 {
		page = 1
	}
	resp.List, resp.TotalNum, err = s.dao.GetServiceConfigList(ctx, req.GetTreeName(), req.GetTreeId(), req.GetKeyword(), req.GetService(), page, pageSize, req.GetName(), req.GetStatus())
	return
}

// GetTreeIds implementation
// http获取已配置的discoveryId `internal:"true"`
func (s *TitansService) GetTreeIds(ctx context.Context, req *v1hpb.TreeIdsReq) (resp *v1hpb.TreeIdsResp, err error) {
	resp = &v1hpb.TreeIdsResp{}
	if "" == req.GetTreeName() {
		err = ecode.Error(ecode.InvalidParam, "tree_name 为空")
		return
	}
	resp.List, _ = s.dao.GetTreeIds(ctx, req.GetTreeName())
	return
}

// GetByTreeId implementation
// grpc获取tree_id对应的全部配置 `internal:"true"`
func (s *TitansService) GetByTreeId(ctx context.Context, req *v1pb.TreeIdReq) (resp *v1pb.TreeIdResp, err error) {
	resp = &v1pb.TreeIdResp{}
	treeId := req.GetTreeId()
	if 0 == treeId {
		err = ecode.Error(ecode.InvalidParam, "tree_id 为空")
		return
	}
	value, err := s.dao.GetServiceConfig(ctx, treeId)
	if nil != err {
		err = ecode.Error(ecode.ServerErr, "内部错误")
		return
	}
	resp.List = value
	return
}

// GetMyTreeApps implementation
// 获取用户的应用树
func (s *TitansService) GetMyTreeApps(ctx context.Context, req *v1hpb.TreeAppsReq, cookie string, user string) (resp *v1hpb.TreeAppsResp, err error) {
	resp = &v1hpb.TreeAppsResp{}
	resp.List, err = s.Nodes(ctx, user, req.Node, req.Team, cookie)
	return
}

// SyncConfig job同步配置
func (s *TitansService) SyncConfig() {
	nodesList := []string{"live"}
	ctx := context.Background()
	for {
		// 获取tree_id列表
		for _, node := range nodesList {
			list, err := s.dao.GetTreeIds(ctx, node)
			if err != nil {
				log.Error("[Titans][syncConfig error][get tree_id list], err: %+v", err)
				continue
			}
			//通过tree_id获取对应配置
			if len(list) > 0 {
				for _, treeId := range list {
					value, err := s.dao.GetServiceConfig(ctx, treeId)
					if err != nil {
						log.Error("[Titans][syncConfig error][get tree_id config], tree_id"+strconv.Itoa(int(treeId))+"err: %+v", err)
						continue
					}
					s.titansCache.Put(fmt.Sprintf(_titansLruCache, treeId), value)
				}
			}
			log.Info("[Titans][syncConfig][sync success], node:" + node + ", tree_id num :" + strconv.Itoa(len(list)))
		}
		time.Sleep(30 * time.Second)
	}
}

// GetEasyList implementation
// 获取运营数据列表 `internal:"true"`
func (s *TitansService) GetEasyList(ctx context.Context, req *v1hpb.EasyGetReq) (resp *v1hpb.EasyGetResp, err error) {
	resp = &v1hpb.EasyGetResp{}
	resp.List = []*v1hpb.EasyList{}
	list := s.dao.GetEasyRecord(ctx, "live")

	if list.Value == "" {
		return
	}
	dbValue := map[string][]int64{}
	err = json.Unmarshal([]byte(list.Value), &dbValue)
	ids := dbValue["list"]
	if err != nil || len(ids) == 0 {
		return
	}

	pageSize := 100
	dbIds := make([]int64, 0)
	begin := (int(req.GetPage()) - 1) * pageSize
	end := int(req.GetPage()) * pageSize
	for k, v := range ids {
		if req.GetId() == v {
			dbIds = append(dbIds, v)
			break
		}
		if k >= begin && k < end {
			ids = append(ids, v)
		}
	}
	if req.GetId() != 0 && len(dbIds) == 0 {
		return
	}
	dbRes := []*model.ServiceModel{}
	if req.GetId() != 0 && len(dbIds) != 0 {
		dbRes, err = s.dao.GetListByIds(ctx, dbIds)
	} else {
		dbRes, err = s.dao.GetListByIds(ctx, ids)
	}
	if err != nil {
		return
	}
	for _, v := range dbRes {
		item := &v1hpb.EasyList{
			TreeName: v.TreeName,
			TreePath: v.TreePath,
			TreeId:   v.TreeId,
			Keyword:  v.Keyword,
			Name:     v.Name,
		}
		resp.List = append(resp.List, item)
	}
	return
}

// SetEasyList implementation
// 设置运营列表 `internal:"true"`
func (s *TitansService) SetEasyList(ctx context.Context, req *v1hpb.EasySetReq) (resp *v1hpb.EasySetResp, err error) {
	resp = &v1hpb.EasySetResp{}
	if req.GetId() == 0 {
		err = ecode.Error(ecode.InvalidParam, "记录id不能为空")
		return
	}
	ids := []int64{req.GetId()}
	list, err := s.dao.GetListByIds(ctx, ids)
	if err != nil {
		return
	}
	if len(list) == 0 || list[0] == nil {
		err = ecode.Error(ecode.InvalidParam, "记录信息为空")
		return
	}
	template := list[0].Template
	if template == 0 {
		err = ecode.Error(ecode.InvalidParam, "同步到运营操作的记录必须选择个模型")
		return
	}
	if list[0].Name == "" {
		err = ecode.Error(ecode.InvalidParam, "同步到运营操作的记录必须加个描述")
		return
	}
	treeName := list[0].TreeName
	//获取easyList
	easy := s.dao.GetEasyRecord(ctx, treeName)
	easyValue := map[string][]int64{}
	easyValue["list"] = make([]int64, 0)
	if easy == nil || easy.Id == 0 {
		//新增
		easyValue["list"] = append(easyValue["list"], req.GetId())
	} else {
		err = json.Unmarshal([]byte(easy.Value), &easyValue)
		if err != nil {
			return
		}
		if easyValue["list"] == nil {
			easyValue["list"] = append(easyValue["list"], req.GetId())
		} else {
			//编辑的时候，进行更新
			for _, v := range easyValue["list"] {
				if v == req.GetId() {
					err = ecode.Error(ecode.InvalidParam, "已经在运营列表里了哟~")
					return
				}
			}
			easyValue["list"] = append(easyValue["list"], req.GetId())
		}
	}

	value, err := json.Marshal(easyValue)
	if err != nil {
		return
	}
	valueDb := string(value)

	//插入数据
	resp.EId, err = s.dao.SetEasyRecord(ctx, treeName, valueDb, easy.Id)
	return
}
