package email

import (
	"context"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"

	"go-common/app/job/main/videoup-report/model/email"
	"go-common/library/cache/redis"
	"sync"
)

var tplTest = &email.Template{
	Headers: map[string][]string{
		email.TO:      {"chenxi01@bilibili.com"},
		email.SUBJECT: {"nothing at all"},
	},
	Body:        "testhahaha",
	ContentType: "text/plain",
}

func TestEmailSendMail(t *testing.T) {
	convey.Convey("SendMail", t, func(ctx convey.C) {
		tplTest.Headers[email.FROM] = []string{d.email.Username}
		d.SendMail(tplTest)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestEmailsendEmailLog(t *testing.T) {
	var (
		to     = []string{"chenxi01@bilibili.com"}
		cc     = []string{""}
		result = "成功"
	)
	convey.Convey("sendEmailLog", t, func(ctx convey.C) {
		d.sendEmailLog(tplTest, to, cc, result)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

//分片
func batch(tk []int64, length int) (path [][]int64) {
	ll := len(tk) / length
	if len(tk)%length > 0 {
		ll++
	}
	path = [][]int64{}
	item := []int64{}
	for i := 0; i < len(tk); i++ {
		if i > 0 && i%length == 0 {
			path = append(path, item)
			item = []int64{}
		}
		item = append(item, tk[i])
	}
	if len(item) > 0 {
		path = append(path, item)
	}

	return
}

func TestEmailPushToRedis(t *testing.T) {
	uid := int64(123)
	uids := []int64{}
	speedThreshold := d.c.Mail.SpeedThreshold
	overlimit := speedThreshold * d.c.Mail.OverspeedThreshold
	for i := 0; i < overlimit*2; i++ {
		uids = append(uids, uid)
	}

	tplTest.UID = uid
	path := batch(uids, speedThreshold)
	len1 := len(uids)
	path = append(path, batch(uids, speedThreshold-1)...)
	len2 := 2 * len1
	path = append(path, batch(uids, speedThreshold)...)
	cnt := 0
	convey.Convey("连续发送邮件，间隔出现超限名额", t, func(ctx convey.C) {
		tplTest.Headers[email.FROM] = []string{d.email.Username}
		for index, task := range path {
			now := time.Now().UnixNano()
			for i := range task {
				cnt++
				isfast, key, err := d.PushToRedis(context.TODO(), tplTest)
				//_, _, err := d.PushToRedis(context.email.TODO(), tplTest)
				convey.So(err, convey.ShouldBeNil)
				t.Logf("cnt=%d,index=%d, i=%d, detector=%+v", cnt, index, i, d.detector)

				if cnt < overlimit+speedThreshold { //快速，探查阶段
					convey.So(isfast, convey.ShouldEqual, false)
					convey.So(key, convey.ShouldEqual, email.MailKey)
				} else if cnt == overlimit+speedThreshold { //快速，确认为超限，提供超限名单
					convey.So(isfast, convey.ShouldEqual, true)
					convey.So(key, convey.ShouldEqual, email.MailFastKey)
				} else if cnt < len1+speedThreshold { //快速，探查阶段，保留上一次的超限名单
					convey.So(isfast, convey.ShouldEqual, false)
					convey.So(key, convey.ShouldEqual, email.MailFastKey)
				} else if cnt < len2+overlimit+speedThreshold { //慢速/快速探查阶段，第一次慢速时清空上一次的超限名单
					convey.So(isfast, convey.ShouldEqual, false)
					convey.So(key, convey.ShouldEqual, email.MailKey)
				} else if cnt == len2+overlimit+speedThreshold { //快速，确认为超限，提供超限名单
					convey.So(isfast, convey.ShouldEqual, true)
					convey.So(key, convey.ShouldEqual, email.MailFastKey)
				} else { //快速，探查阶段，保留上一次的超限名单
					convey.So(isfast, convey.ShouldEqual, false)
					convey.So(key, convey.ShouldEqual, email.MailFastKey)
				}
			}

			if diff := now + 1e9 - time.Now().UnixNano(); diff > 0 {
				time.Sleep(time.Duration(diff))
			}
		}
	})
}

func TestEmailStart(t *testing.T) {
	convey.Convey("email Start", t, func(ctx convey.C) {
		err := d.Start(email.MailKey)
		convey.So(err, convey.ShouldBeNil)

		err = d.Start(email.MailKey + "_1")
		convey.So(err, convey.ShouldEqual, redis.ErrNil)
	})
}

func TestEmailBatchStart(t *testing.T) {
	wg := sync.WaitGroup{}
	convey.Convey("email Start\r\n", t, func(ctx convey.C) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				err := d.Start(email.MailKey)
				t.Logf("start to push normal email, time=%d\r\n", time.Now().Unix())
				if err == redis.ErrNil {
					t.Logf("normal email stopped\r\n")
					break
				}
				ctx.So(err, convey.ShouldBeNil)
			}
		}()

		go func() {
			defer wg.Done()
			for {
				err := d.Start(email.MailFastKey)
				t.Logf("start to push fast email, time=%d\r\n", time.Now().Unix())
				if err == redis.ErrNil {
					t.Logf("fast email stopped\r\n")
					break
				}
				ctx.So(err, convey.ShouldBeNil)
			}
		}()

		wg.Wait()
		t.Logf("end")
	})
}
