package main

import (
	"flag"

	"go-common/app/service/ep/footman/conf"
	"go-common/app/service/ep/footman/model"
	"go-common/app/service/ep/footman/service"
	"go-common/library/cache/memcache"
	"go-common/library/container/pool"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"
	"go-common/library/net/netutil/breaker"
	"go-common/library/time"
)

func main() {
	var (
		workspaceType string
		reportType    string
		workspaceID   string
		iteration     string
		reportPath    string
		isTimeStr     string

		startTime       string
		endTime         string
		importPath      string
		importSheetName string
	)

	flag.StringVar(&workspaceType, "t", "", "端类型，android或ios")
	flag.StringVar(&reportType, "r", "", "报表类型")
	flag.StringVar(&workspaceID, "w", "", "项目ID")
	flag.StringVar(&iteration, "i", "", "迭代名称")
	flag.StringVar(&reportPath, "p", "", "报表文件路径")
	flag.StringVar(&isTimeStr, "d", "n", "是否详细时间")

	flag.StringVar(&startTime, "s", "", "开始时间 2018-01-01")
	flag.StringVar(&endTime, "e", "", "结束时间 2019-01-01")
	flag.StringVar(&importPath, "f", "", "导入文件")
	flag.StringVar(&importSheetName, "sn", "", "导入文件sheet name")
	flag.Parse()

	c := &conf.Config{
		Tapd: &conf.Tapd{
			RetryTime: 5,
			WaitTime:  time.Duration(100000),
		},

		HTTPClient: &xhttp.ClientConfig{
			App: &xhttp.App{
				Key:    "c05dd4e1638a8af0",
				Secret: "7daa7f8c06cd33c5c3067063c746fdcb",
			},
			Dial:      time.Duration(2000000000),
			Timeout:   time.Duration(10000000000),
			KeepAlive: time.Duration(60000000000),
			Breaker: &breaker.Config{
				Window:  time.Duration(10000000000),
				Sleep:   time.Duration(2000000000),
				Bucket:  10,
				Ratio:   0.5,
				Request: 100,
			},
		},
		Mail: &conf.Mail{
			Host:        "smtp.exmail.qq.com",
			Port:        465,
			Username:    "merlin@bilibili.com",
			Password:    "",
			NoticeOwner: []string{"fengyifeng@bilibili.com"},
		},
		Memcache: &conf.Memcache{
			Expire: time.Duration(10000000),
			Config: &memcache.Config{
				Name:         "merlin",
				Proto:        "tcp",
				Addr:         "172.22.33.137:11216",
				DialTimeout:  time.Duration(1000),
				ReadTimeout:  time.Duration(1000),
				WriteTimeout: time.Duration(1000),
				Config: &pool.Config{
					Active:      10,
					IdleTimeout: time.Duration(1000),
				},
			},
		},
	}
	s := service.New(c)

	switch reportType {

	//纯测试时长统计列表区分_ios_android
	case "testtime":
		if err := s.TestTimeReport(workspaceID, workspaceType, iteration, reportPath, model.Test, model.IPS, model.SPS, model.SCPS); err != nil {
			log.Error("Error happened when generating test time reprot:error(%v)", err)
		}

		//待测试时长统计 不区分ios android
	case "waittest":
		if err := s.WaitTimeReport(workspaceID, iteration, reportPath, model.Test, model.IPS, model.SPS, model.SCPS); err != nil {
			log.Error("Error happened when generating wait for test time reprot:error(%v)", err)
		}

		//产品验收时长统计 不区分ios android
	case "experiencetime":
		if err := s.TestTimeReport(workspaceID, workspaceType, iteration, reportPath, model.Experience, model.IPS, model.SPS, model.SCPS); err != nil {
			log.Error("Error happened when generating experience time reprot:error(%v)", err)
		}
	case "delayedstory":
		if err := s.DelayedStoryReport(workspaceID, workspaceType, iteration, reportPath, model.IPS, model.SPS, model.SCPS); err != nil {
			log.Error("Error happened when generating delayed story reprot:error(%v)", err)
		}

		//测试中打回需求统计
	case "testrejected":
		if err := s.RejectedStoryReport(workspaceID, workspaceType, iteration, reportPath, model.Test, model.IPS, model.SPS, model.SCPS); err != nil {
			log.Error("Error happened when generating test rejected story reprot:error(%v)", err)
		}

		//产品验收打回需求统计
	case "experiencerejected":
		if err := s.RejectedStoryReport(workspaceID, workspaceType, iteration, reportPath, model.Experience, model.IPS, model.SPS, model.SCPS); err != nil {
			log.Error("Error happened when generating experience rejected story reprot:error(%v)", err)
		}

		//故事墙
	case "storywall":
		isTime := false
		if isTimeStr == "y" {
			isTime = true
		}
		if err := s.StoryWallReport(workspaceID, workspaceType, iteration, reportPath, isTime, model.IPS, model.SPS, model.SCPS, model.CPS); err != nil {
			log.Error("Error happened when generating story wall reprot:error(%v)", err)
		}

	case "storywallreport":
		if err := s.GenStoryReport(workspaceType, importPath, importSheetName, reportPath, iteration, startTime, endTime); err != nil {
			log.Error("Error happened when generating experience rejected story reprot:error(%v)", err)
		}
	}

}
