package http

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/databus"
	"go-common/app/admin/main/apm/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func databusProjects(c *bm.Context) {
	username := name(c)
	projects, err := apmSvc.Projects(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		log.Error("%v", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(projects, nil)
}

func databusClusters(c *bm.Context) {
	var clusters []string
	for k := range conf.Conf.Kafka {
		clusters = append(clusters, k)
	}
	c.JSON(clusters, nil)
}

func databusApps(c *bm.Context) {
	v := new(struct {
		Pn      int    `form:"pn" default:"1" validate:"min=1"`
		Ps      int    `form:"ps" default:"20" validate:"min=1"`
		Project string `form:"project"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	projects, err := apmSvc.Projects(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		log.Error("%v", err)
		c.JSON(nil, err)
		return
	}
	var (
		apps  []*databus.App
		total int
	)
	if v.Project == "" {
		err = apmSvc.DBDatabus.Order("id").Offset((v.Pn-1)*v.Ps).Limit(v.Ps).Where("project in (?)", projects).Find(&apps).Error
	} else {
		err = apmSvc.DBDatabus.Order("id").Offset((v.Pn-1)*v.Ps).Limit(v.Ps).Where("project in (?) AND project LIKE ?", projects, "%"+v.Project+"%").Find(&apps).Error
	}
	if err != nil {
		log.Error("apmSvr.Apps error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DBDatabus.Model(&databus.App{}).Where("project in (?)", projects).Count(&total).Error; err != nil {
		log.Error("apmSvr.Apps count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: apps,
		Total: total,
	}
	c.JSON(data, nil)
}

func databusAppAdd(c *bm.Context) {
	v := new(struct {
		Project string `form:"project" validate:"required"`
		Remark  string `form:"remark"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	m := &databus.App{
		AppKey:    genKey(v.Project),
		AppSecret: genSecret(v.Project + strconv.FormatInt(rand.Int63(), 10)),
		Project:   v.Project,
		Remark:    v.Remark,
	}
	db := apmSvc.DBDatabus.Model(&databus.App{}).Create(m)
	if err = db.Error; err != nil {
		log.Error("apmSvc.appAdd error(%v)", err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 1, int64(m.ID), "apmSvc.appAdd", m)
	data := map[string]int{
		"id": db.Value.(*databus.App).ID,
	}
	c.JSON(data, nil)
}

func databusAppEdit(c *bm.Context) {
	v := new(struct {
		ID     int64  `form:"id" validate:"required"`
		Remark string `form:"remark"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Find(app, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DBDatabus.Model(&databus.App{}).Where("id = ?", v.ID).Omit("id").UpdateColumns(v).Error; err != nil {
		log.Error("apmSvc.serviceMod error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "id = ?",
		"Value1":  v.ID,
		"Update":  v,
		"Old":     app,
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 2, v.ID, "apmSvc.serviceMod", sqlLog)
	c.JSON(nil, err)
}

func databusGroups(c *bm.Context) {
	v := new(struct {
		Topic string `form:"topic"`
		Group string `form:"group" default:""`
		Pn    int    `form:"pn" default:"1" validate:"min=1"`
		Ps    int    `form:"ps" default:"20" validate:"min=1"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	projects, err := apmSvc.Projects(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var (
		apps   []*databus.App
		topics []*databus.Topic
		groups []*databus.Group
		total  int
	)
	if err = apmSvc.DBDatabus.Model(&databus.App{}).Where("project in (?)", projects).Find(&apps).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	var appIDs []int
	for _, val := range apps {
		appIDs = append(appIDs, val.ID)
	}
	if v.Topic != "" {
		err = apmSvc.DBDatabus.Where("topic = ?", v.Topic).Find(&topics).Error
	} else {
		err = apmSvc.DBDatabus.Find(&topics).Error
	}
	if err != nil {
		log.Error("apmSvc.groups error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var topicIDS []int
	for _, value := range topics {
		topicIDS = append(topicIDS, value.ID)
	}
	if err = apmSvc.DBDatabus.Raw(`SELECT auth2.id,auth2.group,auth2.app_id,app2.app_key,app2.project,auth2.topic_id,auth2.operation,topic.cluster,topic.topic,auth2.remark,auth2.ctime,auth2.mtime,auth2.percentage,auth2.alarm, if(auth2.number=0,100,auth2.number) as number
		FROM auth2 LEFT JOIN app2 ON app2.id=auth2.app_id LEFT JOIN topic ON topic.id=auth2.topic_id WHERE auth2.is_delete=0 AND auth2.topic_id in (?) and auth2.app_id in (?) and auth2.group LIKE '%`+v.Group+`%' ORDER BY auth2.id LIMIT ?,?`,
		topicIDS, appIDs, (v.Pn-1)*v.Ps, v.Ps).Scan(&groups).Error; err != nil {
		log.Error("apmSvc.Groups error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DBDatabus.Model(&databus.Group{}).Joins("LEFT JOIN app2 ON app2.id=auth2.app_id LEFT JOIN topic ON topic.id=auth2.topic_id").
		Where("auth2.is_delete=0 AND auth2.topic_id in (?) AND auth2.app_id in (?) and auth2.group LIKE '%"+v.Group+"%'", topicIDS, appIDs).Count(&total).Error; err != nil {
		log.Error("apmSvc.Groups count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: groups,
		Total: total,
	}
	c.JSON(data, nil)

}

func databusGroupProjects(c *bm.Context) {
	var apps []*databus.App
	var err error
	if err = apmSvc.DBDatabus.Find(&apps).Error; err != nil {
		log.Error("apmSvc.GroupProjects error(%v)", err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	allPs, err := apmSvc.Projects(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		log.Error("apmSvc.Projects(%s) error(%v)", username, err)
		c.JSON(nil, err)
		return
	}
	ps := make([]string, 0, len(apps))
	for _, app := range apps {
		for _, pj := range allPs {
			if pj == app.Project {
				ps = append(ps, app.Project)
				break
			}
		}
	}
	c.JSON(ps, nil)
}

func databusGroupSubAdd(c *bm.Context) {
	v := new(struct {
		Project string `form:"project" validate:"required"`
		Topic   string `form:"topic" validate:"required"`
		// Operation int8   `form:"operation" validate:"required"`
		Remark string `form:"remark"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	topic := &databus.Topic{}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("topic = ?", v.Topic).First(topic).Error; err != nil {
		log.Error("apmSvc.databusGroupSubAdd topic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DBDatabus.Where("project = ? ", v.Project).First(app).Error; err != nil {
		log.Error("apmSvc.databusGroupSubAdd project error(%v)", err)
		c.JSON(nil, err)
		return
	}
	m := &databus.Group{
		Group:      genGroup(v.Topic, v.Project, 1),
		TopicID:    topic.ID,
		AppID:      app.ID,
		Operation:  1, // NOTE: sub
		Remark:     v.Remark,
		Number:     100,
		Percentage: "80",
	}
	exist := 0
	if err = apmSvc.DBDatabus.Model(&databus.Group{}).Where("`group`=? AND topic_id=?", m.Group, m.TopicID).Count(&exist).Error; err != nil {
		log.Error("apmSvc.databusGroupSubAdd exist group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if exist > 0 {
		m.Group = genGroup(v.Topic, v.Project, 11)
	}
	db := apmSvc.DBDatabus.Create(m)
	if err = db.Error; err != nil {
		log.Error("apmSvc.databusGroupSubAdd error(%v)", err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 1, int64(m.ID), "apmSvc.databusGroupSubAdd", m)
	data := map[string]int{
		"id": db.Value.(*databus.Group).ID,
	}
	c.JSON(data, nil)
}

func databusGroupPubAdd(c *bm.Context) {
	v := new(struct {
		Project string `form:"project" validate:"required"`
		Topic   string `form:"topic" validate:"required"`
		// Operation int8   `form:"operation" validate:"required"`
		Remark string `form:"remark"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	topic := &databus.Topic{}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("topic = ?", v.Topic).First(topic).Error; err != nil {
		log.Error("apmSvc.databusGroupPubAdd topic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DBDatabus.Where("project = ? ", v.Project).First(app).Error; err != nil {
		log.Error("apmSvc.databusGroupPubAdd project error(%v)", err)
		c.JSON(nil, err)
		return
	}
	m := &databus.Group{
		Group:     genGroup(v.Topic, v.Project, 2),
		TopicID:   topic.ID,
		AppID:     app.ID,
		Operation: 2, // NOTE: pub
		Remark:    v.Remark,
	}
	db := apmSvc.DBDatabus.Create(m)
	if err = db.Error; err != nil {
		log.Error("apmSvc.databusGroupPubAdd error(%v)", err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 1, int64(m.ID), "apmSvc.databusGroupPubAdd", m)
	data := map[string]int{
		"id": db.Value.(*databus.Group).ID,
	}
	c.JSON(data, nil)
}

func databusTopics(c *bm.Context) {
	v := new(struct {
		Pn      int    `form:"pn" default:"1" validate:"min=1"`
		Ps      int    `form:"ps" default:"20" validate:"min=1"`
		Cluster string `form:"cluster"`
		Topic   string `form:"topic"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		tps   []*databus.Topic
		total int
	)
	if v.Cluster != "" && v.Topic == "" {
		err = apmSvc.DBDatabus.Where("cluster = ?", v.Cluster).Order("id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&tps).Error
	} else if v.Cluster == "" && v.Topic != "" {
		err = apmSvc.DBDatabus.Where("topic LIKE ?", "%"+v.Topic+"%").Order("id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&tps).Error
	} else if v.Cluster != "" && v.Topic != "" {
		err = apmSvc.DBDatabus.Where("cluster = ? AND topic LIKE ?", v.Cluster, "%"+v.Topic+"%").Order("id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&tps).Error
	} else {
		err = apmSvc.DBDatabus.Order("id").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&tps).Error
	}
	if err != nil {
		log.Error("apmSvr.Topics error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.Cluster != "" && v.Topic == "" {
		err = apmSvc.DBDatabus.Model(&databus.Topic{}).Where("cluster = ?", v.Cluster).Count(&total).Error
	} else if v.Cluster == "" && v.Topic != "" {
		err = apmSvc.DBDatabus.Model(&databus.Topic{}).Where("topic LIKE ?", "%"+v.Topic+"%").Count(&total).Error
	} else if v.Cluster != "" && v.Topic != "" {
		err = apmSvc.DBDatabus.Model(&databus.Topic{}).Where("cluster = ? AND topic LIKE ?", v.Cluster, "%"+v.Topic+"%").Count(&total).Error
	} else {
		err = apmSvc.DBDatabus.Model(&databus.Topic{}).Count(&total).Error
	}
	if err != nil {
		log.Error("apmSvr.Topics cluster(%s) count error(%v)", v.Cluster, err)
		c.JSON(nil, err)
		return
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: tps,
		Total: total,
	}
	c.JSON(data, nil)
}

func databusTopicNames(c *bm.Context) {
	v := new(struct {
		Cluster string `form:"cluster" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var tps []*databus.Topic
	if err = apmSvc.DBDatabus.Where("cluster = ?", v.Cluster).Find(&tps).Error; err != nil {
		log.Error("apmSvr.Topics error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var ns []string
	for _, t := range tps {
		ns = append(ns, t.Topic)
	}
	c.JSON(ns, nil)
}

func databusTopicAdd(c *bm.Context) {
	res := make(map[string]interface{}, 1)
	v := new(struct {
		Topic   string `form:"topic" validate:"required"`
		Cluster string `form:"cluster" validate:"required"`
		Remark  string `form:"remark"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	f := strings.HasSuffix(v.Topic, "-T")
	if !f {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	topic := &databus.Topic{
		Topic:   v.Topic,
		Cluster: v.Cluster,
		Remark:  v.Remark,
	}
	db := apmSvc.DBDatabus.Create(topic)
	if err = db.Error; err != nil {
		log.Error("apmSvc.databusTopicAdd error(%v)", err)
		res["message"] = "DB创建topic失败"
		c.JSONMap(res, err)
		return
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 1, int64(topic.ID), "apmSvc.databusTopicAdd", topic)
	if conf.Conf.Kafka[topic.Cluster] == nil || !(len(conf.Conf.Kafka[v.Cluster].Brokers) > 0) {
		log.Error("apmSvc.topicAdd  CreateTopic kafka cluster error(%v)", v.Cluster)
		res["message"] = "kafka集群(" + v.Cluster + ")未配置，请手动创建"
		c.JSONMap(res, err)
		return
	}
	if err = service.CreateTopic(conf.Conf.Kafka[v.Cluster].Brokers, v.Topic, conf.Conf.DatabusConfig.Partitions, conf.Conf.DatabusConfig.Factor); err != nil {
		log.Error("apmSvc.topicAdd  CreateTopic kafka error(%v)", err)
		res["message"] = "kafka创建topic失败，请手动创建"
		c.JSONMap(res, err)
		return
	}
	data := map[string]int{
		"id": db.Value.(*databus.Topic).ID,
	}
	c.JSON(data, nil)
}

func databusTopicEdit(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		ID      int64  `form:"id" validate:"required"`
		Cluster string `form:"cluster" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	topic := &databus.Topic{}
	if err = apmSvc.DBDatabus.Find(topic, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	username := name(c)
	ups := map[string]interface{}{
		"cluster":  v.Cluster,
		"operator": username,
	}
	if v.Cluster != topic.Cluster {
		if err = apmSvc.DBDatabus.Model(&databus.Topic{}).Where("`id` = ?", v.ID).Updates(ups).Error; err != nil {
			log.Error("apmSvc.databusTopicEdit error(%v)", err)
			c.JSON(nil, err)
			return
		}
		sqlLog := &map[string]interface{}{
			"SQLType": "update",
			"Where":   "id = ?",
			"Value1":  v.ID,
			"Update":  ups,
			"Old":     topic,
		}
		apmSvc.SendLog(*c, username, 0, 2, v.ID, "apmSvc.databusTopicEdit", sqlLog)
		auth := &databus.OldAuth{}
		if err = apmSvc.DBDatabus.Where("`topic` = ?", topic.Topic).First(auth).Error; err == nil {
			app := &databus.OldApp{}
			if err = apmSvc.DBDatabus.Where("`id` = ?", auth.AppID).First(app).Error; err != nil {
				log.Error("apmSvc.databusTopicEdit first app id(%v) error(%v)", auth.AppID, err)
				c.JSON(nil, err)
				return
			}
			update := map[string]interface{}{
				"cluster": v.Cluster,
			}
			if err = apmSvc.DBDatabus.Model(&databus.OldApp{}).Where("`id` = ?", auth.AppID).Updates(update).Error; err != nil {
				log.Error("apmSvc.databusTopicEdit error(%v)", err)
				c.JSON(nil, err)
				return
			}
			sqlLog := &map[string]interface{}{
				"SQLType": "update",
				"Where":   "id = ?",
				"Value1":  auth.AppID,
				"Update":  update,
				"Old":     app,
			}
			apmSvc.SendLog(*c, username, 0, 2, v.ID, "apmSvc.databusTopicEdit", sqlLog)
		}
		if conf.Conf.Kafka[topic.Cluster] == nil || !(len(conf.Conf.Kafka[v.Cluster].Brokers) > 0) {
			log.Error("apmSvc.topicAdd  CreateTopic kafka cluster error(%v)", v.Cluster)
			res["message"] = "kafka集群(" + v.Cluster + ")未配置，请手动创建"
			c.JSONMap(res, err)
			return
		}
		if err = service.CreateTopic(conf.Conf.Kafka[v.Cluster].Brokers, topic.Topic, conf.Conf.DatabusConfig.Partitions, conf.Conf.DatabusConfig.Factor); err != nil {
			log.Error("apmSvc.topicAdd  CreateTopic kafka error(%v)", err)
			res["message"] = "kafka创建topic失败，请手动创建"
			c.JSONMap(res, err)
			return
		}
	}
	c.JSON(nil, err)
}

func genKey(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	cipherByte := m.Sum(nil)
	return hex.EncodeToString(cipherByte[:8])
}

func genSecret(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	cipherByte := m.Sum(nil)
	return hex.EncodeToString(cipherByte)
}

func genGroup(t, p string, o int8) string {
	group := ""
	ts := []rune(t)
	topic := string(ts[:len(ts)-2])
	group += topic + "-"
	p = strings.Replace(p, "-", ".", -1)
	p = strings.Replace(p, "_", ".", -1)
	ps := strings.Split(p, ".")
	path := ""
	for _, v := range ps {
		path += strings.Title(v)
	}
	group += path
	switch o {
	case 1:
		group += "-S"
	case 2:
		group += "-P"
	case 3:
		group += "-PS"
	case 11:
		group += "-2-S"
	case 4:
		group += "-N"
	}
	return group
}

func renameGroup(t, p, r string, o int8) (group string) {
	group = ""
	ts := []rune(t)
	topic := string(ts[:len(ts)-2])
	group += topic + "-"
	p = strings.Replace(p, "-", ".", -1)
	p = strings.Replace(p, "_", ".", -1)
	ps := strings.Split(p, ".")
	path := ""
	for _, v := range ps {
		path += strings.Title(v)
	}
	group += path
	switch o {
	case 1:
		group += "-" + r + "-S"
	case 2:
		group += "-" + r + "-P"
	case 3:
		group += "-" + r + "-N"
	case 4:
		group += "-" + r + "-N"
	}
	return
}

func databusAlarm(c *bm.Context) {
	var (
		gps    []*databus.Group
		gpMaps []*databus.Alarm
	)
	if err := apmSvc.DBDatabus.Select("auth2.group, app2.project, auth2.alarm, auth2.percentage").Joins("left join app2 on auth2.app_id = app2.id").Where("auth2.is_delete = 0").Find(&gps).Error; err != nil {
		log.Error("apmSvc.databusProject error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, v := range gps {
		a := &databus.Alarm{}
		a.Group = v.Group
		a.Project = v.Project
		a.Alarm = v.Alarm
		if len(v.Percentage) == 0 {
			v.Percentage = "0"
		}
		a.Percentage = v.Percentage
		gpMaps = append(gpMaps, a)
	}
	c.JSON(gpMaps, nil)
}

func databusAlarms(c *bm.Context) {
	var (
		gps     []*databus.Group
		gpMaps  []*databus.Alarms
		records []*databus.Record
		count   int
		err     error
	)
	v := new(struct {
		Pn int `form:"pn" default:"1"`
		Ps int `form:"ps" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = apmSvc.DBDatabus.Select("auth2.group, app2.project, auth2.alarm, auth2.percentage, topic.topic, topic.cluster").Joins("left join app2 on auth2.app_id = app2.id").Joins("left join topic on topic.id = auth2.topic_id").Where("auth2.is_delete = 0").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&gps).Error; err != nil {
		log.Error("apmSvc.databusAlarms error(%v)", err)
		c.JSON(nil, err)
	}
	if err = apmSvc.DBDatabus.Select("auth2.group, app2.project, auth2.alarm, auth2.percentage, topic.topic, topic.cluster").Joins("left join app2 on auth2.app_id = app2.id").Joins("left join topic on topic.id = auth2.topic_id").Where("auth2.is_delete = 0").Model(&databus.Group{}).Count(&count).Error; err != nil {
		log.Error("apmSvc.databusAlarms count error(%v)", err)
		c.JSON(nil, err)
	}
	for _, v := range gps {
		a := &databus.Alarms{}
		a.Group = v.Group
		a.Project = v.Project
		a.Alarm = v.Alarm
		if len(v.Percentage) == 0 {
			v.Percentage = "0"
		}
		a.Percentage = v.Percentage
		a.Topic = v.Topic
		a.Cluster = v.Cluster
		records, err = service.Diff(v.Cluster, v.Topic, v.Group)
		if err != nil {
			log.Error("apmSvc.databusAlarms Diff error(%v)", err)
		}
		a.Diff = records
		gpMaps = append(gpMaps, a)
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: gpMaps,
		Total: count,
	}
	c.JSON(data, nil)
}

func databusGroupDelete(c *bm.Context) {
	v := new(struct {
		Group string `form:"group" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", v.Group).Find(group).Error; err != nil {
		log.Error("apmSvc.databusGroupDelete select error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if group.IsDelete == 1 {
		return
	}
	username := name(c)
	ups := map[string]interface{}{
		"is_delete": 1,
		"operator":  username,
	}
	if err = apmSvc.DBDatabus.Model(group).Where("`group` = ?", v.Group).Updates(ups).Error; err != nil {
		log.Error("apmSvc.databusGroupDelete delete error(%v)", err)
		c.JSON(nil, err)
		return
	}
	apmSvc.SendLog(*c, username, 0, 4, int64(group.ID), "apmSvc.databusGroupDelete", group)
	c.JSON(nil, err)
}

func databusGroupRename(c *bm.Context) {
	v := new(struct {
		Group   string `form:"group" validate:"required"`
		NewName string `form:"newname" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	reg := regexp.MustCompile(`.+?-.+?-2-.+?`)
	match := reg.FindAllString(v.Group, -1)
	if len(match) == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	newGroup := strings.Replace(v.Group, "-2-", "-"+v.NewName+"-", -1)
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", v.Group).Find(group).Error; err != nil {
		log.Error("apmSvc.databusGroupRename select error(%v)", err)
		c.JSON(nil, err)
		return
	}
	username := name(c)
	ups := map[string]interface{}{
		"group":    newGroup,
		"operator": username,
	}
	if err = apmSvc.DBDatabus.Model(&databus.Group{}).Where("`group` = ?", v.Group).Updates(ups).Error; err != nil {
		log.Error("apmSvc.databusGroupRename update error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "group = ?",
		"Value1":  v.Group,
		"Update":  ups,
		"Old":     group,
	}
	apmSvc.SendLog(*c, username, 0, 2, 0, "apmSvc.databusGroupRename", sqlLog)
	c.JSON(nil, err)
}

func databusGroupOffset(c *bm.Context) {
	v := new(struct {
		Group string `form:"group" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", v.Group).Find(group).Error; err != nil {
		log.Error("databusGroupOffset select group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	topic := &databus.Topic{}
	if err = apmSvc.DBDatabus.Where("id = ?", group.TopicID).Find(topic).Error; err != nil {
		log.Error("databusGroupOffset select topic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	client, err := service.NewClient(conf.Conf.Kafka[topic.Cluster].Brokers, topic.Topic, group.Group)
	if err != nil {
		log.Error("service.NewClient() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	defer client.Close()
	marked, err := client.OffsetMarked()
	if err != nil {
		log.Error("client.OffsetMarked() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	new, err := client.OffsetNew()
	if err != nil {
		log.Error("client.OffsetNew() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	type record struct {
		Partition int32 `json:"partition"`
		Diff      int64 `json:"diff"`
		New       int64 `json:"new"`
	}
	records := make([]*record, len(new))
	for partition, offset := range new {
		r := &record{
			Partition: partition,
			New:       offset,
		}
		if tmp, ok := marked[partition]; ok {
			if tmp == -1 {
				r.Diff = -1
			} else {
				r.Diff = offset - tmp
			}
		}
		records[partition] = r
	}
	c.JSON(records, nil)
}

func databusGroupMarked(c *bm.Context) {
	v := new(struct {
		Group string `form:"group" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", v.Group).Find(group).Error; err != nil {
		log.Error("databusGroupMarked select group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	topic := &databus.Topic{}
	if err = apmSvc.DBDatabus.Where("id = ?", group.TopicID).Find(topic).Error; err != nil {
		log.Error("databusGroupMarked select topic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	client, err := service.NewClient(conf.Conf.Kafka[topic.Cluster].Brokers, topic.Topic, group.Group)
	if err != nil {
		log.Error("service.NewClient() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	defer client.Close()
	err = client.SeekEnd()
	sqlLog := &map[string]interface{}{
		"SQLType": "kafka",
		"Cluster": conf.Conf.Kafka[topic.Cluster].Brokers,
		"Topic":   topic,
		"Group":   group,
		"action":  "SeekEnd",
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 6, 0, "apmSvc.databusGroupMarked", sqlLog)
	if err != nil {
		log.Error("client.SeekEnd() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func databusAlarmEdit(c *bm.Context) {
	v := new(struct {
		ID         int64  `form:"id" validate:"required"`
		Alarm      int8   `form:"alarm" validate:"required"`
		Percentage string `form:"percentage" validate:"required"`
		Remark     string `form:"remark"`
		Number     int    `form:"number"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Find(group, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	ups := map[string]interface{}{
		"alarm":      v.Alarm,
		"percentage": v.Percentage,
		"remark":     v.Remark,
	}
	if v.Number > 0 {
		ups["number"] = v.Number
	}
	if err = apmSvc.DBDatabus.Model(&databus.Group{}).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
		log.Error("apmSvc.databusAlarmEdit error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "id = ?",
		"Value1":  v.ID,
		"Update":  ups,
		"Old":     group,
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 2, v.ID, "apmSvc.databusAlarmEdit", sqlLog)
	c.JSON(nil, err)
}

func databusGroupBegin(c *bm.Context) {
	v := new(struct {
		Group string `form:"group" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", v.Group).Find(group).Error; err != nil {
		log.Error("databusGroupMarked select group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	topic := &databus.Topic{}
	if err = apmSvc.DBDatabus.Where("id = ?", group.TopicID).Find(topic).Error; err != nil {
		log.Error("databusGroupMarked select topic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	client, err := service.NewClient(conf.Conf.Kafka[topic.Cluster].Brokers, topic.Topic, group.Group)
	if err != nil {
		log.Error("service.NewClient() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	defer client.Close()
	err = client.SeekBegin()
	if err != nil {
		log.Error("client.SeekBegin() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func databusApplyPubAdd(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		Project     string `form:"project" validate:"required"`
		Remark      string `form:"remark" validate:"required"`
		TopicName   string `form:"topic" validate:"required"`
		TopicRemark string `form:"topic_remark"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	clstr := ""
	topic := &databus.Topic{}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("project = ? ", v.Project).First(app).Error; err != nil {
		log.Error("databusApplyPubAdd project error(%v)", err)
		res["message"] = "project不存在"
		c.JSONMap(res, err)
		return
	}
	f := strings.HasSuffix(v.TopicName, "-T")
	if !f {
		log.Error("databusApplyPubAdd topic_name not standard error(%v)", err)
		res["message"] = "topic名称不合法"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if err = apmSvc.DBDatabus.Where("topic = ?", v.TopicName).First(topic).Error; err == nil {
		clstr = topic.Cluster
		// log.Error("databusApplyPubAdd cluster error(%v)", err)
		// res["message"] = "该topic不存在"
		// c.JSONMap(res, err)
		// return
	}
	username := name(c)
	m := &databus.Apply{
		Group:       genGroup(v.TopicName, v.Project, 2),
		Cluster:     clstr,
		TopicRemark: v.TopicRemark,
		TopicName:   v.TopicName,
		AppID:       app.ID,
		Project:     app.Project,
		Operation:   2, // NOTE: pub
		State:       1,
		Operator:    username,
		Remark:      v.Remark,
	}
	//db := apmSvc.DBDatabus.Create(m)
	tmp := &databus.Apply{}
	db := apmSvc.DBDatabus.Where(databus.Apply{Group: m.Group, State: 4}).Assign(m).FirstOrCreate(tmp)
	if err = db.Error; err != nil {
		log.Error("apmSvc.databusApplyPubAdd create error(%v)", err)
		res["message"] = "创建失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, err)
}

func databusApplySubAdd(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		Project   string `form:"project" validate:"required"`
		Remark    string `form:"remark" validate:"required"`
		TopicName string `form:"topic" validate:"required"`
		Rename    string `form:"rename"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	clstr := ""
	topic := &databus.Topic{}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("project = ? ", v.Project).First(app).Error; err != nil {
		log.Error("databusApplySubAdd project error(%v)", err)
		res["message"] = "project不存在"
		c.JSONMap(res, err)
		return
	}
	f := strings.HasSuffix(v.TopicName, "-T")
	if !f {
		log.Error("databusApplySubAdd topic_name not standard error(%v)", err)
		res["message"] = "topic不合法"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if err = apmSvc.DBDatabus.Where("topic = ?", v.TopicName).First(topic).Error; err == nil {
		clstr = topic.Cluster
		// log.Error("databusApplySubAdd topic error(%v)", err)
		// res["message"] = "不存在此topic"
		// c.JSONMap(res, err)
		// return
	}
	var groupName string
	var message string
	groupName, message, err = GroupName(v.TopicName, v.Project, v.Rename, 1)
	if err != nil {
		log.Error("databusApplySubAdd groupname error(%v)", err)
		res["message"] = message
		c.JSONMap(res, err)
		return
	}
	username := name(c)
	m := &databus.Apply{
		Group:     groupName,
		Cluster:   clstr,
		TopicName: v.TopicName,
		AppID:     app.ID,
		Project:   app.Project,
		Operation: 1, // NOTE: sub
		State:     1,
		Operator:  username,
		Remark:    v.Remark,
	}
	db := apmSvc.DBDatabus.Create(m)
	if err = db.Error; err != nil {
		log.Error("apmSvc.databusApplySubAdd create error(%v)", err)
		res["message"] = "创建失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, err)
}

func databusApplyEdit(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		ID          int    `form:"id" validate:"required"`
		TopicName   string `form:"topic" validate:"required"`
		Project     string `form:"project" validate:"required"`
		Remark      string `form:"remark" validate:"required"`
		TopicRemark string `form:"topic_remark"`
		Rename      string `form:"rename"`
		Cluster     string `form:"cluster"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	apply := &databus.Apply{}
	if err = apmSvc.DBDatabus.Where("id = ?", v.ID).First(apply).Error; err != nil {
		log.Error("databusApplyEdit id error(%v)", err)
		res["message"] = "id有误"
		c.JSONMap(res, err)
		return
	}
	if !(apply.State == 1 || apply.State == 2) {
		log.Error("databusApplyEdit apply state error(%v)", err)
		res["message"] = "只有申请中和打回才可修改"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	clstr := v.Cluster
	topic := &databus.Topic{}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("project = ? ", v.Project).First(app).Error; err != nil {
		log.Error("databusApplyEdit project error(%v)", err)
		res["message"] = "project不存在"
		c.JSONMap(res, err)
		return
	}
	f := strings.HasSuffix(v.TopicName, "-T")
	if !f {
		log.Error("databusApplyEdit topic_name not standard error(%v)", err)
		res["message"] = "topic不合法"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if err = apmSvc.DBDatabus.Where("topic = ?", v.TopicName).First(topic).Error; err == nil {
		clstr = topic.Cluster
		// log.Error("databusApplyEdit topic error(%v)", err)
		// res["message"] = "不存在此topic"
		// c.JSONMap(res, err)
		// return
	}
	ups := map[string]interface{}{
		"cluster":    clstr,
		"topic_name": v.TopicName,
		"app_id":     app.ID,
		"project":    v.Project,
		"remark":     v.Remark,
	}
	if apply.State == 2 {
		ups["state"] = 1
	}
	if apply.Operation == 2 && apply.State != 3 {
		ups["topic_remark"] = v.TopicRemark
	}
	var groupName string
	var message string
	if (v.TopicName != apply.TopicName) || (v.Project != apply.Project) || len(v.Rename) > 0 {
		groupName, message, err = GroupName(v.TopicName, v.Project, v.Rename, apply.Operation)
		if err != nil {
			log.Error("databusApplySubAdd groupname error(%v)", err)
			res["message"] = message
			c.JSONMap(res, err)
			return
		}
		ups["group"] = groupName
	}
	if err = apmSvc.DBDatabus.Model(apply).Where("id = ?", apply.ID).Updates(ups).Error; err != nil {
		log.Error("databusApplyEdit update error(%v)", err)
		res["message"] = "修改失败"
		c.JSONMap(res, err)
		return
	}
	c.JSON(nil, err)
}

func databusApplyList(c *bm.Context) {
	v := new(struct {
		TopicName string `form:"topic"`
		Group     string `form:"group"`
		State     int8   `form:"state"`
		Pn        int    `form:"pn" default:"1"`
		Ps        int    `form:"ps" default:"20"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	var (
		cs    []*databus.Apply
		count int
		lk    = "%" + v.Group + "%"
		in    []int8
	)
	if v.State > 0 {
		in = append(in, v.State)
	} else {
		in = []int8{0, 1, 2, 3, 4}
	}
	if v.TopicName != "" {
		err = apmSvc.DBDatabus.Where("topic_name = ? and `group` like ? and state in (?)", v.TopicName, lk, in).Order("mtime desc").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cs).Error
	} else {
		err = apmSvc.DBDatabus.Select("group_apply.*,notify.id as nid,notify.gid,notify.offset,notify.state as nstate,notify.filter,notify.concurrent,notify.callback,notify.zone").Joins(
			"left join notify on  notify.gid=group_apply.id ").Where("group_apply.group like ? and group_apply.state in (?)",
			lk, in).Order("mtime desc").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&cs).Error
	}
	if err != nil {
		log.Error("apmSvc.databusApplyList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if v.TopicName != "" {
		err = apmSvc.DBDatabus.Where("topic_name = ? and `group` like ? and state in (?)", v.TopicName, lk, in).Model(&databus.Apply{}).Count(&count).Error
	} else {
		err = apmSvc.DBDatabus.Joins("left join notify on group_apply.id = notify.gid ").Where("group_apply.group like ? and group_apply.state in (?)", lk, in).Model(&databus.NotifyGroup{}).Count(&count).Error
	}
	if err != nil {
		log.Error("apmSvc.databusApplyList count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var ids []int64
	for _, val := range cs {
		ids = append(ids, val.Nid)
	}
	filter := []*databus.Filter{}
	if err = apmSvc.DBDatabus.Where("nid in (?)", ids).Find(&filter).Error; err != nil && err != gorm.ErrRecordNotFound {
		log.Error("apmSvc.databusApplyList filter error(%v)", err)
		c.JSON(nil, err)
		return
	}
	type result struct {
		*databus.Apply
	}
	var results []*result
	for _, val := range cs {
		rs := new(result)
		rs.Apply = val
		for _, vv := range filter {
			if val.Nid == int64(vv.Nid) {
				err = json.Unmarshal([]byte(vv.Filters), &rs.FilterList)
			}
		}
		results = append(results, rs)
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: results,
		Total: count,
	}
	c.JSON(data, nil)
}

func databusApplyApprovalProcess(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		ID    int  `form:"id" validate:"required"`
		State int8 `form:"state" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	apply := &databus.Apply{}
	if err = apmSvc.DBDatabus.Where("id = ?", v.ID).First(apply).Error; err != nil {
		log.Error("databusApplyApprovalProcess id error(%v)", err)
		res["message"] = "id有误"
		c.JSONMap(res, err)
		return
	}
	if !(apply.State == 1 || apply.State == 2) {
		log.Error("databusApplyApprovalProcess apply.state error(%v)", apply.State)
		res["message"] = "只有申请中和打回才可审核"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if !(v.State == 2 || v.State == 3 || v.State == 4) {
		log.Error("databusApplyApprovalProcess v.state error(%v)", v.State)
		res["message"] = "state值范围为2，3，4"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	ups := map[string]interface{}{
		"state": v.State,
	}
	tx := apmSvc.DBDatabus.Begin()
	if v.State == 3 && apply.State != 3 {
		if apply.Cluster == "" {
			log.Error("databusApplyApprovalProcess apply.cluster err")
			res["message"] = "新申请的Topic需填写对应Cluster才可审核"
			c.JSONMap(res, ecode.RequestErr)
			tx.Rollback()
			return
		}
		app := &databus.App{}
		if err = tx.Where("id = ?", apply.AppID).First(app).Error; err != nil {
			log.Error("databusApplyApprovalProcess app first error(%v)", err)
			res["message"] = "申请数据状态有误"
			c.JSONMap(res, err)
			tx.Rollback()
			return
		}
		topic := &databus.Topic{}
		//可能要创建topic，创建group,更新apply
		if err = tx.Where("topic=? and cluster=?", apply.TopicName, apply.Cluster).First(topic).Error; err == gorm.ErrRecordNotFound {
			//不存在就要创建
			topic = &databus.Topic{
				Topic:   apply.TopicName,
				Cluster: apply.Cluster,
				Remark:  apply.TopicRemark,
			}
			db := tx.Create(topic)
			if err = db.Error; err != nil {
				log.Error("apmSvc.databusApplyApprovalProcess create topic error(%v)", err)
				res["message"] = "创建topic失败"
				c.JSONMap(res, err)
				tx.Rollback()
				return
			}
			username := name(c)
			//create topic on kafka cluster
			apmSvc.SendLog(*c, username, 0, 1, int64(topic.ID), "apmSvc.databusTopicAdd", topic)
			if conf.Conf.Kafka[topic.Cluster] == nil || !(len(conf.Conf.Kafka[topic.Cluster].Brokers) > 0) {
				log.Error("apmSvc.topicAdd  CreateTopic kafka cluster error(%v)", topic.Cluster)
				res["message"] = "kafka集群(" + topic.Cluster + ")未配置，请手动创建"
				c.JSONMap(res, err)
				tx.Rollback()
				return
			}
			if err = service.CreateTopic(conf.Conf.Kafka[topic.Cluster].Brokers, topic.Topic, conf.Conf.DatabusConfig.Partitions, conf.Conf.DatabusConfig.Factor); err != nil {
				log.Error("apmSvc.topicAdd  CreateTopic kafka error(%v)", err)
				res["message"] = "kafka创建topic失败，请手动创建"
				c.JSONMap(res, err)
				tx.Rollback()
				return
			}
		}
		number := 0
		if apply.Operation == 1 {
			number = 100
		}
		//创建group
		group := &databus.Group{
			Group:      apply.Group,
			TopicID:    topic.ID,
			AppID:      apply.AppID,
			Operation:  apply.Operation,
			Remark:     apply.Remark,
			Number:     number,
			Percentage: "80",
		}
		db := tx.Create(group)
		if err = db.Error; err != nil {
			log.Error("apmSvc.databusApplyApprovalProcess create group error(%v)", err)
			res["message"] = "创建group失败"
			c.JSONMap(res, err)
			tx.Rollback()
			return
		}
		ups["group"] = group.Group
		if strings.HasSuffix(group.Group, "-N") {
			notify := &databus.Notify{}
			nus := map[string]interface{}{
				"gid":   group.ID,
				"state": 1,
			}
			if err = tx.Model(notify).Where("gid = ? ", apply.ID).Updates(nus).Error; err != nil {
				log.Error("databusApplyApprovalProcess updates error(%v)", err)
				res["message"] = "修改状态失败"
				c.JSONMap(res, err)
				tx.Rollback()
				return
			}
		}
	}
	//更新apply
	if err = tx.Model(apply).Where("id = ?", apply.ID).Updates(ups).Error; err != nil {
		log.Error("databusApplyApprovalProcess updates error(%v)", err)
		res["message"] = "修改状态失败"
		c.JSONMap(res, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	c.JSON(nil, err)
}

// GroupName get group name
func GroupName(TopicName, Project, Rename string, operation int8) (groupName, message string, err error) {
	auth := &databus.Group{}
	if len(Rename) != 0 {
		groupName = renameGroup(TopicName, Project, Rename, operation)
	} else {
		groupName = genGroup(TopicName, Project, operation)
	}
	if err = apmSvc.DBDatabus.Where("`group` = ?", groupName).First(auth).Error; err == nil {
		if err = apmSvc.DBDatabus.Where("`group` = ?", groupName).First(auth).Error; err == nil {
			log.Error("GroupName group2 error(%v)", err)
			err = ecode.SvenRepeat
			message = "该group(" + groupName + ")可能己存在，请重命名group或查看group列表"
			return
		}
	}
	apply := &databus.Apply{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", groupName).First(apply).Error; err == nil {
		if apply.State == 4 {
			ups := map[string]interface{}{
				"state": 1,
			}
			if err = apmSvc.DBDatabus.Model(&databus.Apply{}).Where("id = ?", apply.ID).Updates(ups).Error; err != nil {
				log.Error("GroupName update error(%v)", err)
				message = "状态变更失败"
			}
			return
		}
		log.Error("GroupName group3 error(%v)", err)
		err = ecode.SvenRepeat
		message = "该group(" + groupName + ")己在申请中，请重命名group或查看group申请列表"
		return
	}
	err = nil
	return
}

func databusNotifyApplyAdd(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		Project     string `form:"project" validate:"required"`
		Remark      string `form:"remark" validate:"required"`
		TopicName   string `form:"topic" validate:"required"`
		TopicRemark string `form:"topic_remark"`
		Filter      int8   `form:"filter"`
		Concurrent  int8   `form:"concurrent" validate:"required"`
		Callback    string `form:"callback" validate:"required"`
		Filters     string `form:"filters"`
		Rename      string `form:"rename"`
		Zone        string `form:"notify_zone"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	if len(v.Zone) == 0 {
		v.Zone = "sh001"
	}
	topic := &databus.Topic{}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("project = ? ", v.Project).First(app).Error; err != nil {
		log.Error("databusNotifyApplyAdd project error(%v)", err)
		res["message"] = "project不存在"
		c.JSONMap(res, err)
		return
	}
	f := strings.HasSuffix(v.TopicName, "-T")
	if !f {
		log.Error("databusNotifyApplyAdd topic_name not standard error(%v)", err)
		res["message"] = "topic不合法"
		c.JSONMap(res, ecode.RequestErr)
		return
	}
	if err = apmSvc.DBDatabus.Where("topic = ?", v.TopicName).First(topic).Error; err != nil {
		log.Error("databusNotifyApplyAdd topic error(%v)", err)
		res["message"] = "不存在此topic"
		c.JSONMap(res, err)
		return
	}
	var groupName string
	var message string
	groupName, message, err = GroupName(v.TopicName, v.Project, v.Rename, 4)
	if err != nil {
		log.Error("databusNotifyApplyAdd groupname error(%v)", err)
		res["message"] = message
		c.JSONMap(res, err)
		return
	}
	username := name(c)
	m := &databus.Apply{
		Group:       groupName,
		Cluster:     topic.Cluster,
		TopicRemark: v.TopicRemark,
		TopicName:   v.TopicName,
		AppID:       app.ID,
		Operation:   4, // NOTE: sub
		State:       1,
		Operator:    username,
		Remark:      v.Remark,
		Project:     v.Project,
	}
	tx := apmSvc.DBDatabus.Begin()
	db := tx.Create(m)
	if err = db.Error; err != nil {
		tx.Rollback()
		log.Error("apmSvc.databusNotifyApplyAdd create group error(%v)", err)
		res["message"] = "创建group申请表失败"
		c.JSONMap(res, err)
		return
	}
	n := &databus.Notify{
		Gid:        m.ID,
		State:      0,
		Filter:     v.Filter,
		Concurrent: v.Concurrent,
		Callback:   v.Callback,
		Zone:       v.Zone,
	}
	db = tx.Create(n)
	if err = db.Error; err != nil {
		tx.Rollback()
		log.Error("apmSvc.databusNotifyApplyAdd create notify error(%v)", err)
		res["message"] = "notify表创建失败"
		c.JSONMap(res, err)
		return
	}
	filter := &databus.Filter{
		Nid:     n.ID,
		Filters: v.Filters,
	}
	db = tx.Create(filter)
	if err = db.Error; err != nil {
		tx.Rollback()
		log.Error("apmSvc.databusNotifyApplyAdd create filter error(%v)", err)
		res["message"] = "filter表创建失败"
		c.JSONMap(res, err)
		return
	}
	tx.Commit()
	c.JSON(nil, err)
}

func databusNotifyEdit(c *bm.Context) {
	res := map[string]interface{}{}
	v := new(struct {
		ID  int `form:"id" validate:"required"`
		NID int `form:"n_id" validate:"required"`
		// Cluster   string `form:"cluster" validate:"required"`
		TopicName  string `form:"topic" validate:"required"`
		Project    string `form:"project" validate:"required"`
		Rename     string `form:"rename"`
		Remark     string `form:"remark"`
		State      int8   `form:"state"`
		Filter     int8   `form:"filter"`
		Concurrent int8   `form:"concurrent"`
		Callback   string `form:"callback" validate:"required"`
		Filters    string `form:"filters"`
		Status     int8   `form:"status"` //1审核表，2auth2表
		Zone       string `form:"notify_zone" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("project = ? ", v.Project).First(app).Error; err != nil {
		log.Error("databusNotifyEdit project error(%v)", err)
		res["message"] = "project不存在"
		c.JSONMap(res, err)
		return
	}
	ups := map[string]interface{}{
		"remark": v.Remark,
	}
	if len(v.Rename) != 0 {
		var groupName, message string
		groupName, message, err = GroupName(v.TopicName, v.Project, v.Rename, 4)
		if err != nil {
			log.Error("databusNotifyEdit groupname error(%v)", err)
			res["message"] = message
			c.JSONMap(res, err)
			return
		}
		ups["group"] = groupName
	}
	if v.Status == 1 {
		ups["topic_name"] = v.TopicName
		ups["app_id"] = app.ID
		ups["project"] = v.Project
		apply := &databus.Apply{}
		if err = apmSvc.DBDatabus.Model(apply).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
			log.Error("databusNotifyEdit appply update error(%v)", err)
			res["message"] = "修改失败"
			c.JSONMap(res, err)
			return
		}
	} else {
		group := &databus.Group{}
		if err = apmSvc.DBDatabus.Model(group).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
			log.Error("databusNotifyEdit auth2 update error(%v)", err)
			res["message"] = "修改失败"
			c.JSONMap(res, err)
			return
		}
	}
	notify := &databus.Notify{}
	if err = apmSvc.DBDatabus.Where("id = ?", v.NID).First(&notify).Error; err != nil {
		log.Error("databusNotifyEdit notify first error(%v)", err)
		c.JSON(nil, err)
		return
	}
	nups := map[string]interface{}{
		"filter":     v.Filter,
		"concurrent": v.Concurrent,
		"callback":   v.Callback,
		"state":      v.State,
		"zone":       v.Zone,
	}
	tx := apmSvc.DBDatabus.Begin()
	if err = tx.Model(notify).Where("id = ?", v.NID).Updates(nups).Error; err != nil {
		log.Error("databusNotifyEdit updates error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	filter := &databus.Filter{}
	fups := map[string]interface{}{
		"filters": v.Filters,
	}
	if err = tx.Where("nid = ?", v.NID).First(&filter).Error; err != nil {
		filter.Nid = v.NID
		filter.Filters = v.Filters
		tx.Create(filter)
	} else {
		if err = tx.Model(filter).Where("nid = ?", v.NID).Updates(fups).Error; err != nil {
			c.JSON(nil, err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	c.JSON(nil, err)
}

func databusNotifyList(c *bm.Context) {
	v := new(struct {
		Topic string `form:"topic"`
		Group string `form:"group" default:""`
		Pn    int    `form:"pn" default:"1" validate:"min=1"`
		Ps    int    `form:"ps" default:"20" validate:"min=1"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	username := name(c)
	projects, err := apmSvc.Projects(c, username, c.Request.Header.Get("Cookie"))
	if err != nil {
		c.JSON(nil, err)
		return
	}
	var (
		apps   []*databus.App
		topics []*databus.Topic
		groups []*databus.Group
		total  int
	)
	if err = apmSvc.DBDatabus.Model(&databus.App{}).Where("project in (?)", projects).Find(&apps).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	var appIDs []int
	for _, val := range apps {
		appIDs = append(appIDs, val.ID)
	}
	if v.Topic != "" {
		err = apmSvc.DBDatabus.Where("topic = ?", v.Topic).Find(&topics).Error
	} else {
		err = apmSvc.DBDatabus.Find(&topics).Error
	}
	if err != nil {
		log.Error("apmSvc.groups error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var topicIDS []int
	for _, value := range topics {
		topicIDS = append(topicIDS, value.ID)
	}
	if err = apmSvc.DBDatabus.Raw(`SELECT auth2.id,auth2.group,auth2.app_id,app2.app_key,app2.project,auth2.topic_id,topic.cluster,topic.topic,auth2.remark,auth2.ctime,
		auth2.mtime,auth2.percentage,auth2.alarm, if(auth2.number=0,100,auth2.number) as number,
	    notify.id as nid, notify.gid,notify.callback,notify.concurrent,notify.filter,notify.state,filters.filters,notify.zone
		FROM auth2 LEFT JOIN app2 ON app2.id=auth2.app_id LEFT JOIN topic ON topic.id=auth2.topic_id
		INNER JOIN notify ON notify.gid=auth2.id
		LEFT JOIN filters ON filters.nid = notify.id 
		WHERE auth2.is_delete=0 AND auth2.topic_id in (?)  and auth2.app_id in (?) and auth2.group LIKE '%`+v.Group+`%' ORDER BY auth2.id LIMIT ?,?`,
		topicIDS, appIDs, (v.Pn-1)*v.Ps, v.Ps).Scan(&groups).Error; err != nil {
		log.Error("apmSvc.Groups error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if err = apmSvc.DBDatabus.Model(&databus.Group{}).Joins(`LEFT JOIN app2 ON app2.id=auth2.app_id LEFT JOIN topic ON topic.id=auth2.topic_id 
			INNER JOIN notify ON notify.gid=auth2.id LEFT JOIN filters ON filters.nid = notify.id`).Where("auth2.is_delete=0 AND auth2.topic_id in (?) AND auth2.app_id in (?) and auth2.group LIKE '%"+v.Group+"%'", topicIDS, appIDs).Count(&total).Error; err != nil {
		log.Error("apmSvc.Groups count error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, g := range groups {
		if g.Filters != "" {
			json.Unmarshal([]byte(g.Filters), &g.FilterList)
		}
	}
	data := &Paper{
		Pn:    v.Pn,
		Ps:    v.Ps,
		Items: groups,
		Total: total,
	}
	c.JSON(data, nil)
}

// func databusNotifyFilterAdd(c *bm.Context) {
//
// 	username, err := permit(c, user.DatabusNotifyEdit)
// 	if err != nil {
// 		log.Error("permit(%v, %s)", username, user.DatabusNotifyEdit)
// 		c.JSON(nil, err)
// 		return
// 	}
// 	type filt struct {
// 		Field     string `json:"field"`
// 		Condition int8   `json:"conidtion"`
// 		Value     string `json:"value"`
// 	}
// 	v := new(struct {
// 		ID      int    `form:"n_id" validate:"required"` //notify id
// 		Filters string `form:"filters" validate:"required"`
// 	})
// 	if err = c.Bind(v); err != nil {
// 		c.JSON(nil, ecode.RequestErr)
// 		res["message"] = "参数有误"
// 		return
// 	}
// 	filters := []filt{}
// 	if err = json.Unmarshal([]byte(v.Filters), &filters); err != nil {
// 		c.JSON(nil, ecode.RequestErr)
// 		res["message"] = "filters有误"
// 		return
// 	}
// 	notify := &databus.Notify{}
// 	if err = apmSvc.DBDatabus.Where("id = ?", v.ID).First(notify).Error; err != nil {
// 		c.JSON(nil, ecode.RequestErr)
// 		res["message"] = "未找到该notify"
// 		return
// 	}
// 	tx := apmSvc.DBDatabus.Begin()
// 	var db *gorm.DB
// 	var filter *databus.Filter
// 	for _, val := range filters {
// 		filter.Nid = v.ID
// 		filter.Field = val.Field
// 		filter.Value = val.Value
// 		filter.Condition = val.Condition
// 		db = tx.Create(filter)
// 		if err = db.Error; err != nil {
// 			tx.Rollback()
// 			log.Error("apmSvc.databusNotifyFilterAdd add filters error(%v)", err)
// 			c.JSON(nil, err)
// 			res["message"] = "filter表添加失败"
// 			return
// 		}
// 	}
// 	tx.Commit()
// }

// func databusNotifyFilterEdit(c *bm.Context) {
//
// 	username, err := permit(c, user.DatabusNotifyEdit)
// 	if err != nil {
// 		log.Error("permit(%v, %s)", username, user.DatabusNotifyEdit)
// 		c.JSON(nil, err)
// 		return
// 	}
// 	v := new(struct {
// 		ID        int    `form:"f_id" validate:"required"`
// 		Field     string `form:"field" validate:"required"`
// 		Condition int8   `form:"conidtion" validate:"required"`
// 		Value     string `form:"value" validate:"required"`
// 	})
// 	if err = c.Bind(v); err != nil {
// 		c.JSON(nil, ecode.RequestErr)
// 		res["message"] = "参数有误"
// 		return
// 	}
// 	filter := &databus.Filter{}
// 	if err = apmSvc.DBDatabus.Where("id = ?", v.ID).First(&filter).Error; err != nil {
// 		log.Error("apmSvc.databusNotifyFilterEdit filter first error(%v)", err)
// 		c.JSON(nil, err)
// 		return
// 	}
// 	ups := map[string]interface{}{
// 		"field":     v.Field,
// 		"conidtion": v.Condition,
// 		"value":     v.Value,
// 	}
// 	if err = apmSvc.DBDatabus.Model(filter).Where("id = ?", v.ID).Updates(ups).Error; err != nil {
// 		log.Error("databusNotifyFilterEdit updates error(%v)", err)
// 		c.JSON(nil, err)
// 	}
// }

func databusTopicAll(c *bm.Context) {
	var tps []*databus.Topic
	var err error
	if err = apmSvc.DBDatabus.Find(&tps).Error; err != nil {
		log.Error("apmSvr.Topics error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var ns []string
	for _, t := range tps {
		ns = append(ns, t.Topic)
	}
	c.JSON(ns, nil)
}

func databusGroupNewOffset(c *bm.Context) {
	v := new(struct {
		Group     string `form:"group" validate:"required"`
		Partition int32  `form:"partition" validate:"required"`
		Offset    int64  `form:"offset" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", v.Group).Find(group).Error; err != nil {
		log.Error("databusGroupNewOffset select group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	topic := &databus.Topic{}
	if err = apmSvc.DBDatabus.Where("id = ?", group.TopicID).Find(topic).Error; err != nil {
		log.Error("databusGroupNewOffset select topic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	client, err := service.NewClient(conf.Conf.Kafka[topic.Cluster].Brokers, topic.Topic, group.Group)
	if err != nil {
		log.Error("service.NewClient() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	defer client.Close()
	err = client.NewOffset(v.Partition, v.Offset)
	sqlLog := &map[string]interface{}{
		"SQLType": "kafka",
		"Cluster": conf.Conf.Kafka[topic.Cluster].Brokers,
		"Topic":   topic,
		"Group":   group,
		"action":  "NewOffset",
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 6, 0, "apmSvc.databusGroupNewOffset", sqlLog)
	if err != nil {
		log.Error("client.NewOffset() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func databusGroupTime(c *bm.Context) {
	v := new(struct {
		Group string `form:"group" validate:"required"`
		Time  int64  `form:"time" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group` = ?", v.Group).Find(group).Error; err != nil {
		log.Error("databusGroupTime select group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	topic := &databus.Topic{}
	if err = apmSvc.DBDatabus.Where("id = ?", group.TopicID).Find(topic).Error; err != nil {
		log.Error("databusGroupTime select topic error(%v)", err)
		c.JSON(nil, err)
		return
	}
	client, err := service.NewClient(conf.Conf.Kafka[topic.Cluster].Brokers, topic.Topic, group.Group)
	if err != nil {
		log.Error("service.NewClient() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	defer client.Close()
	err = client.NewTime(v.Time)
	sqlLog := &map[string]interface{}{
		"SQLType": "kafka",
		"Cluster": conf.Conf.Kafka[topic.Cluster].Brokers,
		"Topic":   topic,
		"Group":   group,
		"action":  "NewTime",
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 6, 0, "apmSvc.databusGroupTime", sqlLog)
	if err != nil {
		log.Error("client.NewTime() error(%v)\n", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, err)
}

func databusOpsmind(c *bm.Context) {
	v := new(struct {
		Group      string `form:"group" validate:"required"`
		Percentage int64  `form:"percentage" default:"0"`
		Owner      string `form:"owner"`
		ForTime    int64  `form:"for_time" default:"300"`
		Silence    bool   `form:"silence" default:"false"`
		AdjustID   string `form:"adjust_id"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := &databus.Group{}
	if err = apmSvc.DBDatabus.Where("`group`=?", v.Group).First(group).Error; err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	app := &databus.App{}
	if err = apmSvc.DBDatabus.Where("`id`=?", group.AppID).First(app).Error; err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := apmSvc.Opsmind(c, app.Project, v.Group, "OpsMindAdjustCreate", v.Owner, v.Percentage, v.ForTime, v.Silence)
	if res.RetCode == 290 {
		res, err = apmSvc.OpsmindRemove(c, res.Data.AdjustID, "OpsMindPolicyAdjustRemove")
		if err != nil || res.RetCode != 0 {
			c.JSON(nil, ecode.RequestErr)
			return
		}
		res, err = apmSvc.Opsmind(c, app.Project, v.Group, "OpsMindAdjustCreate", v.Owner, v.Percentage, v.ForTime, v.Silence)
	}
	if res.RetCode != 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	if len(v.AdjustID) > 0 {
		apmSvc.OpsmindRemove(c, v.AdjustID, "OpsMindPolicyAdjustRemove")
	}
	c.JSON(res, err)
}

func databusQuery(c *bm.Context) {
	v := new(struct {
		Group string `form:"group" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	res, err := apmSvc.OpsmindQuery(c, v.Group, "OpsMindAdjustQuery")
	if res.RetCode != 0 {
		err = ecode.RequestErr
		c.JSON(nil, err)
		return
	}
	c.JSON(res, err)
}

func databusOpsmindRemove(c *bm.Context) {
	v := new(struct {
		AdjustID string `form:"adjust_id" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	res, err := apmSvc.OpsmindRemove(c, v.AdjustID, "OpsMindPolicyAdjustRemove")
	if err != nil || res.RetCode != 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(nil, err)
}

func databusMsgFetch(c *bm.Context) {
	v := new(struct {
		Cluster string `form:"cluster" validate:"required"`
		Topic   string `form:"topic" validate:"required"`
		Key     string `form:"key"`
		Start   int64  `form:"start"`
		End     int64  `form:"end"`
		Limit   int    `form:"limit"`
	})
	if err := c.Bind(v); err != nil {
		return
	}
	if v.Limit <= 0 {
		v.Limit = 20
	} else if v.Limit > 100 {
		v.Limit = 100
	}
	// {Topic}APM-MainCommonArch-S
	group := fmt.Sprintf("%sAPM-MainCommonArch-S", strings.Replace(v.Topic, "-T", "", -1))
	if v.Start > 0 {
		kc, ok := conf.Conf.Kafka[v.Cluster]
		if !ok {
			c.JSON(nil, ecode.NothingFound)
			return
		}
		client, err := service.NewClient(kc.Brokers, v.Topic, group)
		if err != nil {
			log.Error("service.NewClient() error(%v)\n", err)
			c.JSON(nil, err)
			return
		}
		if err = client.NewTime(v.Start); err != nil {
			client.Close()
			log.Error("client.NewTime(%d) error(%v)\n", v.Start, err)
			c.JSON(nil, err)
			return
		}
		client.Close()
		time.Sleep(time.Millisecond * 200)
	}
	res, err := service.FetchMessage(c, v.Cluster, v.Topic, group, v.Key, v.Start, v.End, v.Limit)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(res, err)
}

func databusAlarmInit(c *bm.Context) {
	v := new(struct {
		Name string `form:"name"`
		MS   int    `form:"ms"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	group := []*databus.Group{}
	if len(v.Name) > 0 {
		err = apmSvc.DBDatabus.Where("`group` LIKE ?", v.Name+"%").Find(&group).Error
	} else {
		err = apmSvc.DBDatabus.Find(&group).Error
	}
	if err != nil {
		log.Error("databusAlarmInit group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	app := []*databus.App{}
	if err = apmSvc.DBDatabus.Find(&app).Error; err != nil {
		log.Error("databusAlarmInit app error(%v)", err)
		c.JSON(nil, err)
		return
	}
	project := make(map[int]string)
	for _, val := range app {
		project[val.ID] = val.Project
	}
	data := make(map[string]int64)
	go func() {
		for _, val := range group {
			res, _ := apmSvc.Opsmind(context.Background(), project[val.AppID], val.Group, "OpsMindAdjustCreate", "", 80, 300, false)
			data[val.Group] = res.RetCode
			if v.MS > 0 {
				time.Sleep(time.Millisecond * time.Duration(v.MS))
			}
		}
		log.Info("databusAlarmInit data info(%v)", data)
	}()
	c.JSON(nil, err)
}

func databusAlarmAllEdit(c *bm.Context) {
	v := new(struct {
		Alarm      int8   `form:"alarm" validate:"required"`
		Percentage string `form:"percentage" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	ups := map[string]interface{}{
		"percentage": v.Percentage,
	}
	if err = apmSvc.DBDatabus.Model(&databus.Group{}).Where("alarm = ?", v.Alarm).Updates(ups).Error; err != nil {
		log.Error("apmSvc.databusAlarmAllEdit error(%v)", err)
		c.JSON(nil, err)
		return
	}
	sqlLog := &map[string]interface{}{
		"SQLType": "update",
		"Where":   "alarm = ?",
		"Value1":  v.Alarm,
		"Update":  ups,
	}
	username := name(c)
	apmSvc.SendLog(*c, username, 0, 2, 0, "apmSvc.databusAlarmAllEdit", sqlLog)
	c.JSON(nil, err)
}

func databusConsumerAddrs(c *bm.Context) {
	v := new(struct {
		Group string `form:"group" validate:"required"`
	})
	var err error
	if err = c.Bind(v); err != nil {
		return
	}
	c.JSON(apmSvc.DatabusConsumerAddrs(c, v.Group))
}
