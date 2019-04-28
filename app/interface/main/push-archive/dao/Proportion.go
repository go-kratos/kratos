package dao

import (
	"strconv"
	"strings"

	"go-common/app/interface/main/push-archive/conf"
	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
)

// Proportion 普通关注粉丝的灰度比例
type Proportion struct {
	// 粉丝的后2位的最值
	MinValue int
	MaxValue int
}

// NewProportion new
func NewProportion(config []conf.Proportion) (ps []Proportion) {
	var ppt float64
	for _, g := range config {
		valueStartFloat, err := strconv.ParseFloat(strings.TrimSpace(g.ProportionStartFrom), 64)
		if err != nil {
			log.Error("NewProportions config ArcPush.FanGroup.ProportionStartFrom strconv.ParseFloat(%s) error(%v)", g.ProportionStartFrom, err)
			return
		}
		valueStartFrom := int(valueStartFloat)
		// 比例验证
		prop, err := strconv.ParseFloat(strings.TrimSpace(g.Proportion), 64)
		if err != nil {
			log.Error("NewProportions config ArcPush.FanGroup.Proportion(%s) strconv.ParseFloat err(%v)", g.Proportion, err)
			return nil
		}
		if prop*100-float64(int(prop*100)) != 0 {
			// 比例最多保留2位小数
			log.Error("NewProportions config ArcPush.FanGroup.Proportion(%s) must keep at most 2 bits", g.Proportion)
			return nil
		}
		ppt += prop
		if prop <= 0 || prop > 1 || ppt > 1 || ppt <= 0 {
			// 单个数在(0,1]区间，总和在(0, 1]区间
			log.Error("NewProportions config ArcPush.FanGroup.Proportion(%s) must in (0, 1] and sum(%f) in (0, 1]", g.Proportion, ppt)
			return nil
		}
		// 起始值和比例之和必须在00～99之间
		maxValue := int(100*prop-1) + valueStartFrom
		if maxValue >= 100 {
			log.Error("NewProportions config ArcPush.FanGroup.Proportion(%s)+ProportionStartFrom must in [0, 99]", g.Proportion)
			return
		}
		p := Proportion{
			MinValue: valueStartFrom,
			MaxValue: maxValue,
		}
		ps = append(ps, p)
	}
	return
}

// FansByProportion 根据比例分配该关注类型的粉丝, 以全站所有用户作为分母
func (d *Dao) FansByProportion(upper int64, fans map[int64]int) (attentions []int64, specials []int64) {
	for mid, relationType := range fans {
		if relationType == model.RelationSpecial {
			specials = append(specials, mid)
			continue
		}
		if len(d.Proportions) == 0 {
			attentions = append(attentions, mid)
			continue
		}
		// mid最后2位是否在抽样区间内
		last2Digits := int(mid % 100)
		for _, g := range d.Proportions {
			if last2Digits >= g.MinValue && last2Digits <= g.MaxValue {
				attentions = append(attentions, mid)
				break
			}
		}
	}
	return
}
