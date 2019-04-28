package dao

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/push-archive/conf"
	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
)

// FanGroup 粉丝分组
type FanGroup struct {
	// 组名
	Name string
	// 粉丝与up主的关注关系
	RelationType int
	Hitby        string
	// 限制条数
	Limit         int
	PerUpperLimit int
	LimitExpire   int32
	// 本组获取粉丝的hbase信息
	HBaseTable      string
	HBaseFamily     []string
	MsgTemplateDesc string
	MsgTemplate     string
}

func fanGroupKey(relationType int, name string) string {
	return fmt.Sprintf(`%d#%s`, relationType, name)
}

// NewFanGroups 实例化，验证配置, 若配置错误，则panic
func NewFanGroups(config *conf.Config) (grp map[string]*FanGroup) {
	grp = make(map[string]*FanGroup)
	for _, g := range config.ArcPush.FanGroup {
		if g.Name == "" {
			log.Error("NewFanGroups config ArcPush.FanGroup.Name/hitby must not be empty")
			break
		}
		// 粉丝和up主的关系配置验证
		if g.RelationType != model.RelationAttention && g.RelationType != model.RelationSpecial {
			log.Error("NewFanGroups config ArcPush.FanGroup.RelationType not exist(%d)", g.RelationType)
			break
		}
		if g.Hitby != model.GroupDataTypeDefault && g.Hitby != model.GroupDataTypeHBase &&
			g.Hitby != model.GroupDataTypeAbtest && g.Hitby != model.GroupDataTypeAbComparison {
			log.Error("NewFanGroups config ArcPush.FanGroup.hitby(%s) must in [default,hbase]", g.Hitby)
			break
		}
		key := fanGroupKey(g.RelationType, g.Name)
		if _, ok := grp[key]; ok {
			log.Error("NewFanGroups config ArcPush.FanGroup.relationtype(%d) and name(%s) must be unique", g.RelationType, g.Name)
			break
		}
		// hbase配置
		if g.HBaseTable != "" && len(g.HBaseFamily) == 0 {
			log.Error("NewFanGroups config ArcPush.FanGroup.HbaseTable(%s) & HbaseFamily(%v) must exist togather", g.HBaseTable, g.HBaseFamily)
			break
		}
		msgTemp, err := decodeMsgTemplate(g.Name, g.MsgTemplate)
		if err != nil {
			log.Error("NewFanGroups config ArcPush.FanGroup.MsgTemplate(%s) decodeMsgTemplate error(%v)", g.MsgTemplate, err)
			break
		}
		if msgTemp != g.MsgTemplateDesc {
			log.Error("NewFanGroups config ArcPush.FanGroup.MsgTemplate decodeMsgTemplate(%s) must equal to MsgTemplateDesc(%s)", msgTemp, g.MsgTemplateDesc)
			break
		}
		if len(strings.SplitN(msgTemp, "\r\n", 2)) != 2 {
			log.Error("NewFanGroups config ArcPush.FanGroup.MsgTemplate(%s) decodeMsgTemplate(%s) must contains `\r\n`", g.MsgTemplate, msgTemp)
			break
		}
		grp[key] = &FanGroup{
			Name:            strings.TrimSpace(g.Name),
			RelationType:    g.RelationType,
			Hitby:           strings.TrimSpace(g.Hitby),
			Limit:           g.Limit,
			PerUpperLimit:   g.PerUpperLimit,
			LimitExpire:     int32(time.Duration(g.LimitExpire) / time.Second),
			HBaseTable:      strings.TrimSpace(g.HBaseTable),
			HBaseFamily:     g.HBaseFamily,
			MsgTemplateDesc: g.MsgTemplateDesc,
			MsgTemplate:     msgTemp,
		}
	}
	if len(grp) < len(config.ArcPush.FanGroup) {
		fmt.Printf("NewFanGroups failed\r\n\r\n")
		os.Exit(1)
	}
	return
}

// decodeMsgTemplate 将ascii格式的文案模版，解码成中文格式---防止某些服务器不支持中文配置
func decodeMsgTemplate(groupName string, temp string) (decode string, err error) {
	if temp == "" {
		return
	}
	b, err := hex.DecodeString(temp)
	if err != nil {
		log.Error("DecodeMsgTemplate hex.DecodeString error(%v) groupName(%s), temp(%s)", err, groupName, temp)
		return
	}
	buf := new(bytes.Buffer)
	temp = string(b)
	rows := strings.Split(temp[1:len(temp)-1], "\\r\\n")
	lenRows := len(rows) - 1
	for k, row := range rows {
		parts := strings.Split(row, "%s")
		lenParts := len(parts) - 1
		for kp, str := range parts {
			words := strings.Split(str, "\\u")
			for _, w := range words {
				if len(w) < 1 {
					continue
				}
				wi, err := strconv.ParseInt(w, 16, 32)
				if err != nil {
					log.Error("DecodeMsgTemplate error(%v) groupName(%s), decode(%s), word(%s)", err, groupName, temp, w)
					return "", err
				}
				buf.WriteString(fmt.Sprintf("%c", wi))
			}
			if kp >= lenParts {
				continue
			}
			buf.WriteString("%s")
		}
		if k >= lenRows {
			continue
		}
		buf.WriteString("\r\n")
	}
	decode = buf.String()
	return
}

// FansByHBase hbase表中查询粉丝所关联的up主，过滤up不在hbase结果中的粉丝
func (d *Dao) FansByHBase(upper int64, fanGroupKey string, fans *[]int64) (result []int64, excluded []int64) {
	g := d.FanGroups[fanGroupKey]
	// 不过滤
	if len(g.HBaseTable) == 0 {
		result = *fans
		return
	}
	params := model.NewBatchParam(map[string]interface{}{
		"base":     upper,
		"table":    g.HBaseTable,
		"family":   g.HBaseFamily,
		"result":   &result,
		"excluded": &excluded,
		"handler":  d.filterFanByUpper,
	}, nil)
	Batch(fans, 100, 1, params, d.FilterFans)
	return
}

// FansByActiveTime 配置了默认活跃时间，则批量过滤粉丝是否在活跃时间段内，否则不推送；未配置则不过滤活跃时间；若希望没有默认活跃时间但希望过滤活跃时间，配置成[0]
func (d *Dao) FansByActiveTime(hour int, fans *[]int64) (result []int64, excluded []int64) {
	// 未配置则不过滤活跃时间
	if len(d.ActiveDefaultTime) <= 0 {
		result = *fans
		excluded = []int64{}
		return
	}
	params := model.NewBatchParam(map[string]interface{}{
		"base":     hour,
		"table":    "dm_member_push_active_hour",
		"family":   []string{"p"},
		"result":   &result,
		"excluded": &excluded,
		"handler":  d.filterFanByActive,
	}, nil)
	Batch(fans, 100, 1, params, d.FilterFans)
	return
}
