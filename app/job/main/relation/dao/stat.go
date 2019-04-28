package dao

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strconv"

	"go-common/app/job/main/relation/model"
	"go-common/library/log"
)

// DelStatCache is
func (d *Dao) DelStatCache(mid int64) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(context.TODO(), d.clearStatPath, "", params, &res); err != nil {
		log.Error("d.client.Post error(%v)", err)
		return
	}
	if res.Code != 0 {
		log.Error("url(%s) res code(%d) or res.result(%v)", d.clearStatPath+"?"+params.Encode(), res.Code)
	}
	return
}

// FollowerAchieve is
func (d *Dao) FollowerAchieve(c context.Context, mid, follower int64) {
	// 不为 0 结尾的就不检查了
	if follower%10 != 0 {
		return
	}

	flag := model.AchieveFromFollower(follower)
	if flag <= 0 {
		log.Warn("No achieve flag achieved with mid: %d, follower: %d", mid, follower)
		return
	}
	effected, err := d.UserSetAchieveFlag(c, mid, uint64(flag))
	if err != nil {
		log.Error("Failed to set user achieve flag: mid: %d, flag: %d: %+v", mid, flag, err)
		return
	}
	if effected <= 0 {
		log.Info("Already achieved with mid: %d, flag: %d", mid, flag)
		return
	}
	msg := func() string {
		switch flag {
		case model.FollowerAchieve1k:
			return `恭喜您，您的粉丝已经达到1000粉！`
		case model.FollowerAchieve5k:
			return `恭喜您，您的粉丝已经达到5000粉！`
		case model.FollowerAchieve10k:
			return `恭喜您，您的粉丝已经达到1万粉！您将有机会获得UP主粉丝成就奖“一万粉丝成就奖励”， #{戳我领取吧！}{"https://www.bilibili.com/blackboard/activity-zxIQ8otdK.html#/"}`
		case model.FollowerAchieve100k:
			return `恭喜您，您的粉丝已经达到10万粉！您将有机会获得UP主粉丝成就奖“十万粉丝成就奖励”， #{戳我领取吧！}{"https://www.bilibili.com/blackboard/activity-zxIQ8otdK.html#/"}`
		case model.FollowerAchieve1000k:
			return `恭喜您，您的粉丝已经达到100万粉！您将有机会获得UP主粉丝成就奖“百万粉丝成就奖励”， #{戳我领取吧！}{"https://www.bilibili.com/blackboard/activity-zxIQ8otdK.html#/"}`
		}
		if flag >= model.FollowerAchieve100k {
			return fmt.Sprintf(`恭喜您，您的粉丝已达%d万粉！`, int64((math.Log2(float64(flag))-2)*100000/10000))
		}
		return ""
	}()
	if msg != "" {
		log.Info("Follower achieve send message to mid: %d: %s", mid, msg)
		d.SendMsg(c, mid, "粉丝增长通知", msg)
	}
	d.ensureAllFollowerAchieve(c, mid, follower)
}

func (d *Dao) ensureAllFollowerAchieve(c context.Context, mid int64, follower int64) {
	flags := model.AllAchieveFromFollower(follower)
	v := model.AchieveFlag(0)
	for _, f := range flags {
		v |= f
	}
	effected, err := d.UserSetAchieveFlag(c, mid, uint64(v))
	if err != nil {
		log.Error("Failed to ensure user achieve flag: mid: %d, flags: %+v, follower: %d: %+v", mid, flags, follower, err)
		return
	}
	if effected >= 0 {
		log.Warn("Achieve missed on mid: %d, follower: %d, flags: %+v", mid, follower, flags)
		return
	}
}

// SendMsg send message.
func (d *Dao) SendMsg(c context.Context, mid int64, title string, context string) (err error) {
	params := url.Values{}
	params.Set("mc", "2_5_1")
	params.Set("title", title)
	params.Set("data_type", "4")
	params.Set("context", context)
	params.Set("mid_list", fmt.Sprintf("%d", mid))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.followersNotify, "", params, &res); err != nil || res.Code != 0 {
		log.Error("sendMsgURL(%s) code(%d) error(%v)", d.followersNotify+"?"+params.Encode(), res.Code, err)
		return
	}
	log.Info("d.sendMsgURL url(%s) res(%d)", d.followersNotify+"?"+params.Encode(), res.Code)
	return
}
