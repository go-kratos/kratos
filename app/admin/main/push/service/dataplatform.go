package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/push/model"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

const (
	_countryForeign = 1
	_countryChina   = 2

	_sexMale     = 1
	_sexFemale   = 2
	_platformAll = "4" // 全部平台,iOS、小米、华为...

	_attentionTypeUnion = 1 // 多个自选关注条件用union查询 (另一种是join查询)

	_timeLayout = "2006-01-02 15:04:05"

	_retry = 3

	// 检查数据平台当天的计算任务有没有成功
	checkDpDataURL = "http://berserker.bilibili.co/api/archer/project/%d/status"
)

var (
	areas = map[int]string{1: "海外", 2: "中国", 3: "青海", 4: "上海", 5: "安徽", 6: "广西", 7: "贵州", 8: "吉林", 9: "福建", 10: "黑龙江",
		11: "江西", 12: "甘肃", 13: "云南", 14: "湖南", 15: "河北", 16: "山东", 17: "湖北", 18: "广东", 19: "宁夏", 20: "重庆",
		21: "辽宁", 22: "内蒙古", 23: "山西", 24: "澳门", 25: "陕西", 26: "江苏", 27: "四川", 28: "浙江", 29: "海南", 30: "河南",
		31: "北京", 32: "香港", 33: "台湾", 34: "天津", 35: "西藏", 36: "新疆", 37: "其他",
	}
	dpProjects = map[int]string{
		355170: "mid维度数据",
		357214: "buvid维度数据",
		355628: "自选关注数据",
	}
)

// DpTaskInfo .
func (s *Service) DpTaskInfo(ctx context.Context, id int64, job string) (res *model.DPTask, err error) {
	var (
		t     *model.Task
		cond  *model.DPCondition
		group = errgroup.Group{}
	)
	group.Go(func() error {
		t, err = s.dao.TaskInfo(ctx, id)
		return err
	})
	group.Go(func() error {
		cond, err = s.dao.DPCondition(ctx, job)
		return err
	})
	if err = group.Wait(); err != nil {
		return
	}
	if t == nil {
		return
	}
	res = new(model.DPTask)
	res.Task = *t
	p := new(model.DPParams)
	if cond != nil {
		if err = json.Unmarshal([]byte(cond.Condition), &p); err != nil {
			return
		}
	}
	if len(p.VipExpires) == 0 {
		p.VipExpires = make([]*model.VipExpire, 0)
	}
	if len(p.Attentions) == 0 {
		p.Attentions = make([]*model.SelfAttention, 0)
	}
	if len(p.ActivePeriods) == 0 {
		p.ActivePeriods = make([]*model.ActivePeriod, 0)
	}
	res.DPParams = *p
	return
}

// AddDPTask add data platform task
func (s *Service) AddDPTask(ctx context.Context, task *model.DPTask) (err error) {
	var (
		tasks   []model.DPTask
		condTyp = task.Type
	)
	if len(task.DPParams.ActivePeriods) == 0 {
		tasks = append(tasks, *task)
	} else {
		// 多个推送需要添加多条任务
		for _, p := range task.DPParams.ActivePeriods {
			var (
				ptime time.Time
				etime time.Time
				t     = *task
			)
			if ptime, err = time.ParseInLocation(_timeLayout, p.PushTime, time.Local); err != nil {
				return
			}
			if etime, err = time.ParseInLocation(_timeLayout, p.ExpireTime, time.Local); err != nil {
				return
			}
			t.Job = strconv.FormatInt(pushmdl.JobName(ptime.Unix(), t.Summary, t.LinkValue, t.Group), 10)
			t.PushTime = ptime
			t.ExpireTime = etime
			t.ActivePeriod = p.Period
			tasks = append(tasks, t)
		}
	}
	for _, t := range tasks {
		go func(t model.DPTask) {
			var (
				id int64
				e  error
			)
			for i := 0; i < _retry; i++ {
				if id, e = s.dao.AddTask(context.Background(), &t.Task); err == nil {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			if e != nil {
				log.Error("s.AddDPTask(%+v) add task error(%v)", t.Task, e)
				return
			}
			t.DPParams.LevelStr = pushmdl.JoinInts(t.DPParams.Level)
			t.DPParams.AreaStr = pushmdl.JoinInts(t.DPParams.Area)
			t.DPParams.PlatformStr = pushmdl.JoinInts(t.DPParams.Platforms)
			t.DPParams.LikeStr = pushmdl.JoinInts(t.DPParams.Like)
			t.DPParams.ChannelStr = strings.Join(t.DPParams.Channel, ",")
			params, _ := json.Marshal(t.DPParams)
			cond := &model.DPCondition{
				Task:      id,
				Job:       t.Job,
				Type:      condTyp,
				Condition: string(params),
				SQL:       s.parseQuery(task.Type, &t.DPParams),
				Status:    pushmdl.DpCondStatusPending,
			}
			for i := 0; i < _retry; i++ {
				if _, e = s.dao.AddDPCondition(context.Background(), cond); e == nil {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			if e != nil {
				log.Error("s.AddDPTask(%+v) add condition error(%v)", cond, e)
			}
		}(t)
	}
	return
}

// 把查询结构体解析成sql
func (s *Service) parseQuery(typ int, p *model.DPParams) (sql string) {
	log.Info("data platform parse query start(%+v)", p)
	if typ != pushmdl.TaskTypeDataPlatformToken && typ != pushmdl.TaskTypeDataPlatformMid {
		log.Error("data platform parse query task type(%d) error", typ)
		return
	}
	logDate := time.Now().Add(-24 * time.Hour).Format("20060102")
	switch typ {
	case pushmdl.TaskTypeDataPlatformToken:
		sql = fmt.Sprintf("select platform_id,device_token from basic.dws_push_buvid where log_date='%s' ", logDate)
		if p.UserActiveDay > 0 {
			sql += fmt.Sprintf(" and is_visit_%d=1 ", p.UserActiveDay)
		}
		if p.UserNewDay > 0 {
			sql += fmt.Sprintf(" and is_new_%d=1 ", p.UserNewDay)
		}
		if p.UserSilentDay > 0 {
			sql += fmt.Sprintf(" and is_visit_%d=0 ", p.UserSilentDay)
		}
	case pushmdl.TaskTypeDataPlatformMid:
		if len(p.Attentions) > 0 {
			sql = fmt.Sprintf("select t1.platform_id,t1.device_token from (select mid,platform_id,device_token from basic.dws_push_mid where log_date='%s' ", logDate)
		} else {
			sql = fmt.Sprintf("select platform_id,device_token from basic.dws_push_mid where log_date='%s' ", logDate)
		}
		if p.Sex == _sexMale || p.Sex == _sexFemale {
			sql += fmt.Sprintf(" and final_sex=%d ", p.Sex)
		}
		// 年龄段。0:0-17, 1:18-24, 2:25-30, 3:31+, 9:全部
		if p.Age >= 0 && p.Age <= 3 {
			sql += fmt.Sprintf(" and final_age_range=%d ", p.Age)
		}
		if len(p.Level) > 0 {
			sql += fmt.Sprintf(" and level in (%s) ", pushmdl.JoinInts(p.Level))
		}
		// up主 0: 不是 1:是 2:全部
		if p.IsUp == 0 || p.IsUp == 1 {
			sql += fmt.Sprintf(" and is_up=%d ", p.IsUp)
		}
		// 正式会员 0:不是   1:是  2：全部
		if p.IsFormalMember == 0 || p.IsFormalMember == 1 {
			sql += fmt.Sprintf(" and is_formal_member=%d ", p.IsFormalMember)
		}
		if len(p.VipExpires) > 0 {
			var vips []string
			for _, v := range p.VipExpires {
				if v.Begin == "" && v.End == "" {
					continue
				}
				if v.Begin == "" {
					vips = append(vips, fmt.Sprintf(" vip_overdue_time <= '%s'", v.End))
					continue
				}
				if v.End == "" {
					vips = append(vips, fmt.Sprintf(" vip_overdue_time >= '%s'", v.Begin))
					continue
				}
				vips = append(vips, fmt.Sprintf(" vip_overdue_time between '%s' and '%s' ", v.Begin, v.End))
			}
			if len(vips) > 0 {
				sql += fmt.Sprintf(" and (%s) ", strings.Join(vips, "or"))
			}
		}
	}
	if p.PlatformStr != "" && p.PlatformStr != _platformAll {
		for _, v := range p.Platforms {
			if v == pushmdl.PlatformAndroid {
				for _, i := range pushmdl.Platforms {
					if i == pushmdl.PlatformIPhone || i == pushmdl.PlatformIPad {
						continue
					}
					p.Platforms = append(p.Platforms, i)
				}
				break
			}
		}
		sql += fmt.Sprintf(" and platform_id in (%s) ", pushmdl.JoinInts(p.Platforms))
	}
	if len(p.Channel) > 0 {
		var brands []string
		for _, v := range p.Channel {
			if v == "" {
				continue
			}
			brands = append(brands, "'"+v+"'")
		}
		if len(brands) > 0 {
			sql += fmt.Sprintf(" and device_brand in (%s) ", strings.Join(brands, ","))
		}
	}
	if len(p.Area) > 0 {
		var countries, provinces []string
		for _, v := range p.Area {
			if areas[v] == "" {
				continue
			}
			if v == _countryForeign || v == _countryChina {
				countries = append(countries, "'"+areas[v]+"'")
			} else {
				provinces = append(provinces, "'"+areas[v]+"'")
			}
		}
		if len(countries) > 0 {
			sql += fmt.Sprintf(" and country in(%s) ", strings.Join(countries, ","))
		}
		if len(provinces) > 0 {
			sql += fmt.Sprintf(" and province in(%s) ", strings.Join(provinces, ","))
		}
	}
	if len(p.Like) > 0 {
		var likes []string
		for _, id := range p.Like {
			if s.partitions[id] == "" {
				continue
			}
			likes = append(likes, "'"+s.partitions[id]+"'")
		}
		if len(likes) > 0 {
			sql += fmt.Sprintf(" and like_tid in(%s) ", strings.Join(likes, ","))
		}
	}
	if p.ActivePeriod > 0 {
		sql += fmt.Sprintf(" and active_hour_period=%d ", p.ActivePeriod)
	}
	if typ == pushmdl.TaskTypeDataPlatformMid && len(p.Attentions) > 0 {
		var attentions []string
		if p.AttentionsType == _attentionTypeUnion {
			// 并集
			sql += ") t1 join ("
			for _, v := range p.Attentions {
				s := fmt.Sprintf("select mid from basic.dwd_oid_mid_info where log_date='%s' and follow_type=%d ", logDate, v.Type)
				if v.Include != "" {
					s += fmt.Sprintf(" and name='%s' ", strings.Replace(v.Include, "'", "\\'", -1))
				}
				if v.Exclude != "" {
					s += fmt.Sprintf("and name!='%s'", strings.Replace(v.Exclude, "'", "\\'", -1))
				}
				attentions = append(attentions, s)
			}
			sql += strings.Join(attentions, " union ") + ") t2 on t1.mid=t2.mid"
		} else {
			// 交集
			sql += ") t1 join "
			for i, v := range p.Attentions {
				s := fmt.Sprintf("(select mid from basic.dwd_oid_mid_info where log_date='%s' and follow_type=%d ", logDate, v.Type)
				if v.Include != "" {
					s += fmt.Sprintf(" and name='%s' ", strings.Replace(v.Include, "'", "\\'", -1))
				}
				if v.Exclude != "" {
					s += fmt.Sprintf("and name!='%s'", strings.Replace(v.Exclude, "'", "\\'", -1))
				}
				s += fmt.Sprintf(") t%d on t1.mid=t%d.mid ", i+5, i+5)
				attentions = append(attentions, s)
			}
			sql += strings.Join(attentions, " join ")
		}
	}
	log.Info("data platform parse query end(%s)", sql)
	return
}

// UpdateDpCondtionStatus .
func (s *Service) UpdateDpCondtionStatus(ctx context.Context, job string, status int) (err error) {
	s.dao.UpdateDpCondtionStatus(ctx, job, status)
	return
}

type checkDpDataRes struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	ProjectID   int    `json:"projectId"`
	ProjectName string `json:"projectName"`
}

// CheckDpData check whether data platform data is ready
func (s *Service) CheckDpData(ctx context.Context) (err error) {
	group := errgroup.Group{}
	for id := range dpProjects {
		id := id
		now := time.Now()
		bts := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix() * 1000 // 13位时间戳
		ets := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).Unix() * 1000
		params := dpparams(s.c.DPClient.Key)
		params.Set("runStartTime", strconv.FormatInt(bts, 10))
		params.Set("runEndTime", strconv.FormatInt(ets, 10))
		group.Go(func() (err error) {
			u := fmt.Sprintf(checkDpDataURL, id)
			if enc := dpsign(s.c.DPClient.Secret, params); enc != "" {
				u = u + "?" + enc
			}
			req, err := http.NewRequest(http.MethodGet, u, nil)
			if err != nil {
				log.Error("CheckDpData url(%s) error(%v)", u, err)
				return
			}
			res := new(checkDpDataRes)
			if err = s.dpClient.Do(ctx, req, res); err != nil {
				log.Error("CheckDpData url(%s) error(%v)", u+"?"+params.Encode(), err)
				return
			}
			if res.Code != http.StatusOK {
				err = ecode.PushAdminDPNoDataErr
			}
			log.Info("check data platform url(%s) param(%s) res(%+v)", u, params.Encode(), res)
			return
		})
	}
	err = group.Wait()
	return
}

// dpsign calc appkey and appsecret sign.
func dpsign(secret string, params url.Values) (query string) {
	tmp := params.Encode()
	signTmp := dpencode(params)
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(secret)
	b.WriteString(signTmp)
	b.WriteString(secret)
	mh := md5.Sum(b.Bytes())
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(strings.ToUpper(hex.EncodeToString(mh[:])))
	query = qb.String()
	return
}

func dpparams(appkey string) url.Values {
	params := url.Values{}
	params.Set("appKey", appkey)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", "1.0")
	params.Set("signMethod", "md5")
	return params
}

var dpSignParams = []string{"appKey", "timestamp", "version"}

// encode data platform encodes the values into ``URL encoded'' form
// ("bar=baz&foo=quux") sorted by key.
func dpencode(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		found := false
		for _, p := range dpSignParams {
			if p == k {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		vs := v[k]
		prefix := k
		for _, v := range vs {
			buf.WriteString(prefix)
			buf.WriteString(v)
		}
	}
	return buf.String()
}
