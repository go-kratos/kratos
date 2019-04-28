package main

import (
	"flag"
	"os"
	"strings"

	"go-common/app/tool/ci/lib/mail"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
)

type sendInfo struct {
	SenderName  string
	SendTitle   string
	SendContent string
	ExtraData   string
}
type receiverInfo struct {
	ReceiverName []string
}
type config struct {
	Title        string
	SendInfo     sendInfo
	ReceiverInfo receiverInfo
}

func mailToml(cPath string) (conf config, err error) {
	if _, err = toml.DecodeFile(cPath, &conf); err != nil {
		log.Error("Error(%v)", err)
	}
	return
}
func main() {
	var (
		filePath     string
		ciSendTo     string
		pipeStatus   string
		eConf        config
		sendTo       []string
		sendContent  string
		extraData    string
		err          error
		ciProjectURL = os.Getenv("CI_PROJECT_URL")
		ciPipelineId = os.Getenv("CI_PIPELINE_ID")
		ciUserEmail  = os.Getenv("GITLAB_USER_EMAIL")
		sourceBranch = os.Getenv("CI_COMMIT_REF_NAME")
	)
	//log init
	logConfig := &log.Config{
		Stdout: true,
	}
	log.Init(logConfig)

	//pipeline url
	if ciProjectURL == "" {
		log.Warn("Error: Not CI_PROJECT_URL")
	}
	pipelineURL := ciProjectURL + "/pipelines/" + ciPipelineId
	log.Info("url: %v", pipelineURL)

	//send email data from config files
	flag.StringVar(&filePath, "configPath", "", "config path, eg: /data/gitlab/email.toml")
	flag.StringVar(&ciSendTo, "sendTo", "", "send to email, eg: jiangkai@bilibili.com,tangyongqiang@bilibili.com")
	flag.StringVar(&pipeStatus, "pipeStatus", "failed", "pipeStatus, only failed or success")
	flag.StringVar(&extraData, "extraData", "", "email contents")
	flag.Parse()

	if ciSendTo != "" {
		matchList := strings.Split(ciSendTo, ",")
		sendTo = matchList
	} else {
		if filePath != "" {
			if eConf, err = mailToml(filePath); err != nil {
				log.Warn("Warn(%v)", err)
			}
			sendTo = eConf.ReceiverInfo.ReceiverName
			sendContent = eConf.SendInfo.SendContent
			extraData = eConf.SendInfo.ExtraData
		} else {
			sendTo = []string{ciUserEmail}
		}
	}

	// delete saga send mail
	if strings.Contains(ciUserEmail, "zzjs") {
		log.Info("Saga exc pipeline")
	} else {
		sendmail.SendMail(sendTo, pipelineURL, sendContent, sourceBranch, extraData, pipeStatus)
	}
}
