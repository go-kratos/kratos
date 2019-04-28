package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/admin/main/apm/conf"
	cml "go-common/app/admin/main/apm/model/canal"
	"go-common/app/admin/main/apm/model/user"
	cgm "go-common/app/admin/main/config/model"
	"go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/BurntSushi/toml"
)

const (
	_dialTimeout     = "500ms"
	_readTimeout     = "1s"
	_writeTimeout    = "1s"
	_idleTimeout     = "60s"
	_flavor          = "mysql"
	_heartbeatPeriod = 60
	_canreadTimeout  = 90
)

var (
	getBuildIDAPI     = "%s/x/admin/config/build/builds"
	getConfigValueAPI = "%s/x/admin/config/config/value"
	getConfigIDAPI    = "%s/x/admin/config/config/configs"
	getAllErrorsAPI   = "%s/x/internal/canal/errors"
	createConfigAPI   = "%s/x/admin/config/canal/config/create"
	configByNameAPI   = "%s/x/admin/config/canal/name/configs"
	updateConfigAPI   = "%s/x/admin/config/home/config/update"
	checkMasterAPI    = "%s/x/internal/canal/master/check"
	ok                = 0
)

type result struct {
	Data json.RawMessage `json:"data"`
	Code int             `json:"code"`
}

type list []struct {
	ID int `json:"id"`
}

type configs struct {
	BuildFiles fileList `json:"build_files"`
}

type fileList []struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type groupInfo struct {
	Group string `json:"group"`
	Topic string `json:"topic"`
	AppID int    `json:"app_id"`
}

type appInfo struct {
	ID        int    `json:"id"`
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

// ConfigProxy config proxy
func (s *Service) ConfigProxy(c context.Context, method string, uri string, params url.Values, cookie string, args ...interface{}) (data json.RawMessage, err error) {
	//common params
	if params != nil && args == nil {
		params.Set("app_name", "main.common-arch.canal")
		params.Set("tree_id", "3766")
		params.Set("zone", env.Zone)
		params.Set("env", env.DeployEnv)
	}
	res := result{}
	fmt.Println("ConfigProxy uri=", uri, "params=", params.Encode())
	req, err := s.client.NewRequest(method, uri, "", params)
	if err != nil {
		log.Error("s.client.NewRequest() error(%v)", err)
	}
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if err = s.client.Do(c, req, &res); err != nil {
		log.Error("canal.request get url:"+uri+" params:(%v) error(%v)", params.Encode(), err)
		return
	}
	if res.Code != ok {
		log.Error("canal.request get url:"+uri+" params:(%v) returnCode:(%v)", params.Encode(), res.Code)
		return
	}
	data = res.Data
	return
}

func (s *Service) getBuildID(c context.Context, cookie string) (dat list, err error) {
	var (
		params = url.Values{}
		data   = json.RawMessage{}
	)
	uri := fmt.Sprintf(getBuildIDAPI, conf.Conf.Host.SVENCo)
	if data, err = s.ConfigProxy(c, "GET", uri, params, cookie); err != nil {
		log.Error("getBuildID() response errors: %d", err)
		return
	}
	if data == nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &dat); err != nil {
		log.Error("getBuildID() json.Unmarshal errors: %s", err)
		return
	}
	return
}

func (s *Service) getConfigID(c context.Context, id string, cookie string) (dat fileList, err error) {
	var (
		params = url.Values{}
		data   = json.RawMessage{}
	)

	params.Set("build_id", id)
	uri := fmt.Sprintf(getConfigIDAPI, conf.Conf.Host.SVENCo)
	if data, err = s.ConfigProxy(c, "GET", uri, params, cookie); err != nil {
		log.Error("getConfigID() response errors: %d", err)
		return
	}
	if data == nil {
		return
	}
	res := &configs{}
	if err = json.Unmarshal([]byte(data), res); err != nil {
		log.Error("getConfigID() json.Unmarshal errors: %s", err)
		return
	}
	dat = res.BuildFiles
	return
}

func (s *Service) getConfigValue(c context.Context, id string, cookie string) (res *cml.Conf, err error) {
	var (
		params = url.Values{}
		data   = json.RawMessage{}
	)
	params.Set("config_id", id)
	uri := fmt.Sprintf(getConfigValueAPI, conf.Conf.Host.SVENCo)
	if data, err = s.ConfigProxy(c, "GET", uri, params, cookie); err != nil {
		log.Error("getConfigValue() response errors: %d", err)
		return
	}
	if data == nil {
		return
	}
	res = new(cml.Conf)
	if err = json.Unmarshal([]byte(data), &res); err != nil {
		log.Error("getConfigValue() json.Unmarshal errors: %s", err)
		return
	}
	return

}

//ScanByAddrFromConfig 根据addr查询配置信息
func (s *Service) ScanByAddrFromConfig(c context.Context, addr string, cookie string) (res *cml.Conf, err error) {
	var (
		buildIDList  list
		configIDList fileList
	)
	//find build_id
	if buildIDList, err = s.getBuildID(c, cookie); err != nil {
		return
	}
	fmt.Println("BuildID=", buildIDList)

	for _, v := range buildIDList {
		if v.ID != 0 {
			//find config_id
			if configIDList, err = s.getConfigID(c, strconv.Itoa(v.ID), cookie); err != nil {
				return
			}
			fmt.Println("configID=", configIDList)
			for _, val := range configIDList {
				if val.ID != 0 && val.Name == addr+".toml" {
					//find config
					res, err = s.getConfigValue(c, strconv.Itoa(val.ID), cookie)
					return
				}
				continue
			}
		}
		continue
	}
	return
}

//ApplyAdd canal apply
func (s *Service) ApplyAdd(c *bm.Context, v *cml.Canal, username string) (err error) {
	cnt := 0
	f := strings.Contains(v.Addr, ":")

	if !f {
		err = ecode.CanalAddrFmtErr
		return
	}
	if err = s.DBCanal.Model(&cml.Canal{}).Where("addr=?", v.Addr).Count(&cnt).Error; err != nil {
		log.Error("apmSvc.CanalAdd count error(%v)", err)
		err = ecode.RequestErr
		return
	}
	if cnt > 0 {
		err = ecode.CanalAddrExist
		return
	}
	canal := &cml.Canal{
		Addr:    v.Addr,
		Cluster: v.Cluster,
		Leader:  v.Leader,
		BinName: v.BinName,
		BinPos:  v.BinPos,
		Remark:  v.Remark,
	}
	if err = s.DBCanal.Create(canal).Error; err != nil {
		log.Error("apmSvc.CanalAdd create error(%v)", err)
		return
	}
	s.SendLog(*c, username, 0, 1, canal.ID, "apmSvc.CanalAdd", canal)
	return
}

//ApplyDelete canal delete
func (s *Service) ApplyDelete(c *bm.Context, v *cml.ScanReq, username string) (err error) {
	cc := &cml.Canal{}
	if err = s.DBCanal.Model(&cml.Canal{}).Where("addr=?", v.Addr).Find(cc).Error; err != nil {
		log.Error("apmSvc.ApplyDelete count error(%v)", err)
		err = ecode.RequestErr
		return
	}
	id := cc.ID

	if err = s.DBCanal.Model(&cml.Canal{}).Where("id = ?", id).Update("is_delete", 1).Error; err != nil {
		log.Error("apmSvc.canalDelete canalDelete error(%v)", err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "delete",
		"Value":   v.Addr,
	}
	s.SendLog(*c, username, 0, 3, id, "apmSvc.canalDelete", sqlLog)
	return
}

//ApplyEdit canal edit
func (s *Service) ApplyEdit(c *bm.Context, v *cml.EditReq, username string) (err error) {
	cc := &cml.Canal{}
	if err = s.DBCanal.Where("id = ?", v.ID).Find(cc).Error; err != nil {
		log.Error("apmSvc.CanalEdit find(%d) error(%v)", v.ID, err)
		return
	}
	ups := map[string]interface{}{}
	if _, ok := c.Request.Form["bin_name"]; ok {
		ups["bin_name"] = v.BinName
	}
	if _, ok := c.Request.Form["bin_pos"]; ok {
		ups["bin_pos"] = v.BinPos
	}
	if _, ok := c.Request.Form["remark"]; ok {
		ups["remark"] = v.Remark
	}
	if _, ok := c.Request.Form["project"]; ok {
		ups["cluster"] = v.Project
	}
	if _, ok := c.Request.Form["leader"]; ok {
		ups["leader"] = v.Leader
	}
	if err = s.DBCanal.Model(&cml.Canal{}).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
		log.Error("apmSvc.CanalEdit updates error(%v)", err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "id = ?",
		"Value1":  v.ID,
		"Update":  ups,
		"Old":     cc,
	}
	s.SendLog(*c, username, 0, 2, v.ID, "apmSvc.CanalEdit", sqlLog)
	return
}

//GetScanInfo is
func (s *Service) GetScanInfo(c context.Context, v *cml.ScanReq, username string, cookie string) (confData *cml.Results, err error) {

	if confData, err = s.getDocFromConf(c, v.Addr, cookie); err != nil {
		return
	}
	if confData == nil {
		return
	}
	if err = s.Permit(c, username, user.CanalEdit); err != nil {
		confData.Document.Instance.User = ""
		confData.Document.Instance.Password = ""
	}

	return confData, nil
}

//GetAllErrors 调用x/internal/canal/errors 查询错误信息
func (s *Service) GetAllErrors(c context.Context) (errs map[string]string, err error) {
	var (
		data json.RawMessage
		host string
	)

	type v struct {
		Error         string            `json:"error"`
		InstanceError map[string]string `json:"instance_error"`
	}

	if host, err = s.getCanalInstance(c); err != nil {
		return
	}
	uri := fmt.Sprintf(getAllErrorsAPI, host)
	if data, err = s.ConfigProxy(c, "GET", uri, nil, ""); err != nil {
		return
	}
	res := new(v)
	if err = json.Unmarshal([]byte(data), &res); err != nil {
		return
	}
	errs = res.InstanceError
	return
}

func (s *Service) getCanalInstance(c context.Context) (host string, err error) {
	params := url.Values{}
	params.Set("appid", "main.common-arch.canal")
	params.Set("env", env.DeployEnv)
	params.Set("hostname", env.Hostname)
	params.Set("status", "3")
	var ins struct {
		ZoneInstances map[string][]struct {
			Addrs []string `json:"addrs"`
		} `json:"zone_instances"`
	}
	resp, err := s.DiscoveryProxy(c, "GET", "fetch", params)
	if err != nil {
		return
	}
	rb, err := json.Marshal(resp)
	if err != nil {
		return
	}
	json.Unmarshal(rb, &ins)
	inss := ins.ZoneInstances[env.Zone]
	for _, zone := range inss {
		for _, addr := range zone.Addrs {
			if strings.Contains(addr, "http://") {
				host = addr
				break
			}
		}
	}
	return
}

//GetConfigsByName  obtain configs from  configByNameAPI
func (s *Service) getConfigsByName(c context.Context, name string, cookie string) (configs *cgm.Config, err error) {
	var (
		params = url.Values{}
		data   = json.RawMessage{}
		result []*cgm.Config
	)
	params.Set("token", conf.Conf.AppToken)
	params.Set("name", name+".toml")
	uri := fmt.Sprintf(configByNameAPI, conf.Conf.Canal.CANALSVENCo)
	if data, err = s.ConfigProxy(c, "POST", uri, params, cookie); err != nil {
		err = ecode.GetConfigByNameErr
		return
	}
	if data == nil {
		return
	}

	if err = json.Unmarshal([]byte(data), &result); err != nil {
		log.Error("configByNameAPI() json.Unmarshal errors: %s", err)
		return
	}
	if len(result) == 0 {
		return
	}
	configs = result[0]

	return
}

//ProcessCanalList get canal list
func (s *Service) ProcessCanalList(c context.Context, v *cml.ListReq) (listdata *cml.Paper, err error) {
	type errorCanal struct {
		cml.Canal
		Error string `json:"error"`
	}
	var (
		cc     []*errorCanal
		count  int
		errMap map[string]string
	)
	query := " is_delete = 0 "
	if v.Addr != "" {
		query += fmt.Sprintf("and addr = '%s' ", v.Addr)
	}
	if v.Project != "" {
		query += fmt.Sprintf("and cluster = '%s' ", v.Project)
	}
	err = s.DBCanal.Where(query).Order("id DESC").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cc).Error
	if err != nil {
		log.Error("apmSvc.CanalList error(%v)", err)
		return
	}
	err = s.DBCanal.Model(&cml.Canal{}).Where(query).Count(&count).Error
	if err != nil {
		log.Error("apmSvc.CanalList count error(%v)", err)
		return
	}
	// add error info
	if count > 0 {
		if errMap, err = s.GetAllErrors(c); err != nil {
			log.Error("apmSvc.DBCanalApply GetAllErrors error(%v)", err)
		}
		for _, va := range cc {
			va.Error = errMap[va.Addr]
		}
	}
	listdata = &cml.Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: cc,
		Total: count,
	}

	return
}

//ProcessApplyList get apply list
func (s *Service) ProcessApplyList(c context.Context, v *cml.ListReq) (listdata *cml.Paper, err error) {
	var (
		cc    []*cml.Apply
		count int
	)
	query := " state !=3 "
	if v.Addr != "" {
		query += fmt.Sprintf("and addr = '%s' ", v.Addr)
	}
	if v.Project != "" {
		query += fmt.Sprintf("and cluster = '%s' ", v.Project)
	}
	if v.Status > 0 {
		query += fmt.Sprintf("and state = '%d' ", v.Status)

	}
	err = s.DBCanal.Model(&cml.Apply{}).Where(query).Count(&count).Error
	if err != nil {
		log.Error("apmSvc.ApplyList count error(%v)", err)
		return
	}
	err = s.DBCanal.Where(query).Order("id DESC").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cc).Error

	if err != nil {
		log.Error("apmSvc.ApplyList error(%v)", err)
		return
	}
	listdata = &cml.Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: cc,
		Total: count,
	}
	return
}

//ProcessConfigInfo process info to config center
func (s *Service) ProcessConfigInfo(c context.Context, v *cml.ConfigReq, cookie string, username string) (err error) {
	var (
		confs   *cgm.Config
		comment string
		query   []map[string]string
		data    []byte
	)
	if comment, err = s.jointConfigInfo(c, v, cookie); err != nil {
		return
	}
	if confs, err = s.getConfigsByName(c, v.Addr, cookie); err != nil {
		return
	}
	if confs != nil {
		query = []map[string]string{{
			"name":    v.Addr + ".toml",
			"comment": comment,
			"mark":    v.Mark,
		}}
		data, err = json.Marshal(query)
		if err != nil {
			return
		}
		va := url.Values{}
		va.Set("data", string(data))
		va.Set("user", username)

		if _, err = s.updateConfig(c, va, cookie); err != nil {
			return
		}
	} else {
		//never register config's
		params := url.Values{}
		params.Set("comment", comment)
		params.Set("mark", v.Mark)
		params.Set("state", "2")
		params.Set("from", "0")
		params.Set("name", v.Addr+".toml")
		params.Set("user", username)

		if _, err = s.createConfig(c, params, cookie); err != nil {
			return
		}
	}
	if confs, err = s.getConfigsByName(c, v.Addr, cookie); confs == nil {
		return
	}
	if err = s.dao.SetConfigID(confs.ID, v.Addr); err != nil {
		return
	}
	return
}

//CheckMaster canal check
func (s *Service) CheckMaster(c context.Context, v *cml.ConfigReq) (err error) {
	res := result{}
	host, _ := s.getCanalInstance(c)
	params := url.Values{}
	params.Set("user", v.User)
	params.Set("password", v.Password)
	params.Set("addr", v.Addr)

	uri := fmt.Sprintf(checkMasterAPI, host)
	req, err := s.client.NewRequest("POST", uri, "", params)

	if err != nil {
		err = ecode.CheckMasterErr
		return
	}
	if err = s.client.Do(c, req, &res); err != nil {
		err = ecode.CheckMasterErr
		return
	}
	if res.Code != 0 {
		err = ecode.CheckMasterErr
		return
	}
	return
}

//ProcessCanalInfo is
func (s *Service) ProcessCanalInfo(c context.Context, v *cml.ConfigReq, username string) (err error) {
	var (
		cnt  = 0
		info = &cml.Canal{}
	)
	if err = s.DBCanal.Select("addr,cluster,leader").Where("addr=?", v.Addr).Find(info).Error; err == nil {
		if v.Project == "" {
			v.Project = info.Cluster
		}
		if v.Leader == "" {
			v.Leader = info.Leader
		}
	}
	if cnt, err = s.dao.CanalApplyCounts(v); err != nil {
		return
	}
	if cnt > 0 {
		if err = s.dao.CanalApplyEdit(v, username); err != nil {
			return
		}
	} else {
		if err = s.dao.CanalApplyCreate(v, username); err != nil {
			return
		}
	}
	return
}

//getCommentFromConf get document info from configbyname
func (s *Service) getDocFromConf(c context.Context, addr string, cookie string) (confData *cml.Results, err error) {
	var conf *cgm.Config
	if conf, err = s.getConfigsByName(c, addr, cookie); err != nil {
		return
	}
	if conf == nil {
		return
	}
	confData = new(cml.Results)

	if _, err = toml.Decode(conf.Comment, &confData.Document); err != nil {
		err = ecode.ConfigParseErr
		log.Error("comment toml.decode error(%v)", err)
		return
	}
	row := &cml.Apply{}
	err = s.DBCanal.Model(&cml.Apply{}).Select("`cluster`,`leader`").Where("addr=?", addr).Scan(row).Error
	if err == nil {
		confData.Cluster = row.Cluster
		confData.Leader = row.Leader
	} else {
		res := &cml.Canal{}
		err = s.DBCanal.Model(&cml.Canal{}).Select("`cluster`,`leader`").Where("addr=?", addr).Scan(res).Error
		if err != nil {
			log.Error("canalinfo get error(%v)", err)
			return
		}
		confData.Cluster = res.Cluster
		confData.Leader = res.Leader
	}
	confData.ID = conf.ID
	confData.Addr = addr

	return
}

//ProcessConfigInfo is
func (s *Service) jointConfigInfo(c context.Context, v *cml.ConfigReq, cookie string) (comment string, err error) {
	var (
		buf      *bytes.Buffer
		cfg      *cml.Config
		dbs      []*cml.DB
		di       *cml.Databus
		confData *cml.Results
		sid      int64
		schemas  = make(map[string]bool)
	)
	//analysis request params
	if err = json.Unmarshal([]byte(v.Databases), &dbs); err != nil {
		log.Error("apmSvc.jointConfigInfo Unmarshal error(%v)", err)
		err = ecode.DatabasesUnmarshalErr
		return
	}
	for _, db := range dbs {
		if db.Databus == nil {
			continue
		}
		//Find duplicate
		if schemas[db.Schema+db.Databus.Group] {
			log.Error("jointConfigInfo find duplicate databus group (%v)", db.Databus.Group)
			err = ecode.DatabusDuplErr
			return
		}
		schemas[db.Schema+db.Databus.Group] = true
		//get databusinfo
		if di, err = s.databusInfo(db.Databus.Group, db.Databus.Addr, db.Schema, db.Table); err != nil {
			return
		}
		db.Databus = di
	}
	if sid, err = s.getServerID(v.Addr); err != nil {
		err = ecode.CanalAddrFmtErr
		return
	}
	ist := &cml.Instance{
		Caddr:           v.Addr,
		MonitorPeriod:   v.MonitorPeriod,
		ServerID:        sid,
		Flavor:          _flavor,
		HeartbeatPeriod: _heartbeatPeriod,
		ReadTimeout:     _canreadTimeout,
		DB:              dbs,
	}
	if confData, err = s.getDocFromConf(c, v.Addr, cookie); err != nil {
		return
	}
	if confData != nil && v.User == "" {
		ist.User = confData.Document.Instance.User
	} else {
		ist.User = v.User
	}
	if confData != nil && v.Password == "" {
		ist.Password = confData.Document.Instance.Password
	} else {
		ist.Password = v.Password
	}

	cfg = &cml.Config{
		Instance: ist,
	}
	buf = new(bytes.Buffer)
	if err = toml.NewEncoder(buf).Encode(cfg); err != nil {
		return
	}
	comment = buf.String()
	return
}

//CreateConfig send requests to createConfigAPI
func (s *Service) createConfig(c context.Context, params url.Values, cookie string) (res map[string]interface{}, err error) {
	var (
		data = json.RawMessage{}
	)
	uri := fmt.Sprintf(createConfigAPI, conf.Conf.Canal.CANALSVENCo)
	params.Set("token", conf.Conf.AppToken)
	if data, err = s.ConfigProxy(c, "POST", uri, params, cookie); err != nil {
		err = ecode.ConfigCreateErr
		return
	}
	if data == nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &res); err != nil {
		log.Error("updateConfigAPI() json.Unmarshal errors: %s", err)
		return
	}
	return

}

//UpdateConfig send requests to updateConfigAPI
func (s *Service) updateConfig(c context.Context, params url.Values, cookie string) (res map[string]interface{}, err error) {
	var (
		data = json.RawMessage{}
	)
	uri := fmt.Sprintf(updateConfigAPI, conf.Conf.Canal.CANALSVENCo)
	params.Set("token", conf.Conf.AppToken)
	if data, err = s.ConfigProxy(c, "POST", uri, params, cookie); err != nil {
		err = ecode.ConfigUpdateErr
		return
	}
	if data == nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &res); err != nil {
		log.Error("updateConfigAPI() json.Unmarshal errors: %s", err)
		return
	}
	return

}

//DatabusInfo joint databusinfo
func (s *Service) databusInfo(group string, addr string, schema string, table []*cml.Table) (d *cml.Databus, err error) {

	var (
		ai  appInfo
		gi  groupInfo
		act string
	)
	if ai, err = s.getAppInfo(group); err != nil {
		err = ecode.DatabusAppErr
		return
	}
	if gi, _, err = s.getGroupInfo(group); err != nil {
		err = ecode.DatabusGroupErr
		return
	}
	act = s.getAction(group)
	name := "canal/" + schema

	d = &cml.Databus{
		Key:          ai.AppKey,
		Secret:       ai.AppSecret,
		Group:        group,
		Topic:        gi.Topic,
		Action:       act,
		Name:         name,
		Proto:        "tcp",
		Addr:         addr,
		Idle:         1,
		Active:       len(table),
		DialTimeout:  _dialTimeout,
		ReadTimeout:  _readTimeout,
		WriteTimeout: _writeTimeout,
		IdleTimeout:  _idleTimeout,
	}
	return
}

//getAppInfo according group get appinfo
func (s *Service) getAppInfo(group string) (ai appInfo, err error) {
	var table string
	g, new, _ := s.getGroupInfo(group)
	if !new {
		table = "app"
	} else {
		table = "app2"
	}
	err = s.DBDatabus.Table(table).Select("`id`,`app_key`,`app_secret`").Where("`id`= ?", g.AppID).Find(&ai).Error
	if err != nil {
		log.Error("apmSvc.getAppInfo error(%v)", err)
		return
	}
	return
}

//getGroupInfo according group get groupinfo
func (s *Service) getGroupInfo(group string) (gi groupInfo, new bool, err error) {

	err = s.DBDatabus.Table("auth2").Select("auth2.group,topic.topic,auth2.app_id").Joins("join topic on topic.id=auth2.topic_id").Where("auth2.group= ?", group).Scan(&gi).Error
	if err == nil {
		new = true
		return
	}
	err = s.DBDatabus.Table("auth").Select("group_name as `group`,topic,app_id").Where("group_name = ?", group).Find(&gi).Error
	if err != nil {
		log.Error("apmSvc.getGroupInfo  error(%v", err)
		return
	}
	return
}

//getAction according group get action
func (s *Service) getAction(group string) (action string) {

	if strings.HasSuffix(group, "-P") {
		action = "pub"
	} else if strings.HasSuffix(group, "-S") {
		action = "sub"
	} else {
		action = "notify"
	}
	return
}

//TableInfo get array table info
func (s *Service) TableInfo(table string) (infos []*cml.Table, err error) {

	info := strings.Split(table, ",")
	tab := make([]*cml.Table, len(info))
	for i := range info {
		tab[i] = &cml.Table{Name: info[i]}
	}
	infos = tab
	return
}

//getServerID get server id from addr
func (s *Service) getServerID(addr string) (sid int64, err error) {
	ip := strings.Split(addr, ".")
	last := ip[len(ip)-1]
	port := strings.Split(last, ":")
	joint := fmt.Sprintf("%s%s%s", port[len(port)-1], ip[len(ip)-2], port[len(port)-2])
	sid, err = strconv.ParseInt(joint, 10, 64)

	return
}

//SendWechatMessage send  wechat message
func (s *Service) SendWechatMessage(c context.Context, addr, aType, result, sender, note string, receiver []string) (err error) {
	var (
		detail string
		users  = []string{sender}
	)
	users = append(users, receiver...)
	if env.DeployEnv != "prod" {
		detail = fmt.Sprintf("http://%s-%s", env.DeployEnv, "%s")
	} else {
		detail = "http://%s"
	}
	switch aType {
	case cml.TypeMap[cml.TypeApply]:
		detail = fmt.Sprintf(detail, "sven.bilibili.co/#/canal/apply")
	case cml.TypeMap[cml.TypeReview]:
		detail = fmt.Sprintf(detail, "ops-log.bilibili.co/app/kibana 确认canal订阅消息")
	}
	msg := fmt.Sprintf("[sven系统抄送提示]\n发送方:%s\n事件: %s环境 DB:%s canal%s%s\n接收方:%s\n备注:%s\n详情:%s（请复制到浏览器打开）\n", sender, env.DeployEnv, addr, aType, result, strings.Join(receiver, ","), note, detail)
	if err = s.dao.SendWechatToUsers(c, users, msg); err != nil {
		log.Error("apmSvc.SendWechatMessage error(%v)", err)
		return
	}
	return
}

// UpdateProcessTag canal审核通过之后,调用/x/admin/config/canal/tag/update,同步到配置中心发版
func (s *Service) UpdateProcessTag(c context.Context, configID int, cookie string) (err error) {
	client := &http.Client{}

	tokenURL := fmt.Sprintf("%s/x/admin/config/app/envs", conf.Conf.Canal.CANALSVENCo)
	tokenParams := url.Values{}
	tokenParams.Set("app_name", "main.common-arch.canal")
	tokenParams.Set("tree_id", "3766")
	tokenParams.Set("zone", env.Zone)

	tokenReq, err := http.NewRequest(http.MethodGet, tokenURL+"?"+tokenParams.Encode(), nil)
	if err != nil {
		err = ecode.RequestErr
		return
	}
	tokenReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	tokenReq.Header.Set("Cookie", cookie)

	tokenResp, err := client.Do(tokenReq)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(tokenResp.Body)
	if err != nil {
		return
	}
	defer tokenResp.Body.Close()

	var tokenRespObj struct {
		Code int `json:"code"`
		Data []struct {
			Name     string `json:"name"`
			NickName string `json:"nickname"`
			Token    string `json:"token"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &tokenRespObj)
	if err != nil {
		return fmt.Errorf("json unmarshal error: %s, get body: %s", err, body)
	}

	if tokenRespObj.Code != 0 {
		log.Error(" tokenRespObj.Code: %d", tokenRespObj.Code)
		return fmt.Errorf("tokenRespObj.Code: %d", tokenRespObj.Code)
	}

	var token string
	for _, v := range tokenRespObj.Data {
		if v.Name == env.DeployEnv {
			token = v.Token
			break
		}
	}
	log.Info("tokenURL(%s), env.DeployEnv(%v), token(%v)", tokenURL+"?"+tokenParams.Encode(), env.DeployEnv, token)

	updateURL := fmt.Sprintf("%s/x/admin/config/canal/tag/update", conf.Conf.Canal.CANALSVENCo)
	params := url.Values{}
	params.Set("app_name", "main.common-arch.canal")
	params.Set("env", env.DeployEnv)
	params.Set("zone", env.Zone)
	params.Set("config_ids", fmt.Sprintf("%d", configID))
	params.Set("tree_id", "3766")
	params.Set("mark", "canal发版")
	params.Set("user", "canalApprovalProcess")
	params.Set("force", "1")

	if conf.Conf.Canal.BUILD != "" {
		params.Set("build", conf.Conf.Canal.BUILD)
	} else {
		params.Set("build", "docker-1")
	}
	log.Info("env:(%v), zone:(%v), build:(%v)", params.Get("env"), params.Get("zone"), params.Get("build"))
	params.Set("token", token)

	req, err := http.NewRequest(http.MethodPost, updateURL, strings.NewReader(params.Encode()))
	if err != nil {
		err = ecode.RequestErr
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err = s.client.Do(c, req, &res); err != nil || res.Code != 0 {
		log.Error("ApplyApprovalUpdateTag url(%s) err(%v), code(%v), message(%s)", updateURL+params.Encode(), err, res.Code, res.Message)
		err = ecode.RequestErr
		return
	}
	return
}
