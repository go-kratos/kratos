package dao

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_sagaWechatURL = "http://uat-saga-admin.bilibili.co/ep/admin/saga/v2/wechat"
	_gitCommitsAPI = "http://git.bilibili.co/api/v4/projects/682/repository/commits"
)

// ParseUTFiles parse html to get specific file
func (d *Dao) ParseUTFiles(c context.Context, url string) (pkgs []*ut.PkgAnls, err error) {
	var (
		req   *http.Request
		html  []byte
		files []string
	)
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		log.Error("apmSvc.ParseUTFiless error (%v)", err)
		return
	}
	if html, err = d.client.Raw(c, req); err != nil {
		log.Error("apmSvc.ParseUTFiles error (%v)", err)
		return
	}
	reg := regexp.MustCompile(`<option(.*)</option>`)
	files = reg.FindAllString(string(html), -1)
	for _, file := range files {
		cov, _ := strconv.ParseFloat(file[strings.Index(file, "(")+1:strings.Index(file, "%)")], 64)
		pkg := &ut.PkgAnls{
			PKG:      file[strings.Index(file, "go-common") : strings.Index(file, ".go")+3],
			Coverage: cov,
			HTMLURL:  url + "#" + file[strings.Index(file, `="`)+2:strings.Index(file, `">`)],
		}
		pkgs = append(pkgs, pkg)
	}
	return
}

// SendWechatToUsers send msg to multiple users by saga-admin api
func (d *Dao) SendWechatToUsers(c context.Context, users []string, msg string) (err error) {
	var (
		req *http.Request
		b   = &bytes.Buffer{}
		url = _sagaWechatURL + "/message/send"
		res struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}
		body = &ut.WechatUsersMsg{
			ToUser:  users,
			Content: msg,
		}
	)
	if err = json.NewEncoder(b).Encode(body); err != nil {
		log.Error("apmSvc.SendWechatToUsers Error(%v)", err)
		return
	}
	if req, err = http.NewRequest(http.MethodPost, url, b); err != nil {
		log.Error("apmSvc.SendWechatToUsers Error(%v)", err)
		return
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("apmSvc.SendWechatToUsers Error(%v)", err)
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("Http response Code(%v)!=0", res.Code)
		log.Error("apmSvc.SendWechatToUsers Error(%v)", err)
		return
	}
	return
}

// SendWechatToGroup send msg to a group by saga-admin api
func (d *Dao) SendWechatToGroup(c context.Context, chatid, msg string) (err error) {
	var (
		num  int
		req  *http.Request
		b    = &bytes.Buffer{}
		url  = _sagaWechatURL + "/appchat/send"
		body = &ut.WechatGroupMsg{
			ChatID:  chatid,
			MsgType: "text",
			Safe:    0,
		}
	)
	msgBlock := strings.Split(msg, "\n")
	if len(msgBlock)%40 == 0 {
		num = len(msgBlock)/40 - 1
	} else {
		num = len(msgBlock) / 40
	}
	for i := 0; i <= num; i++ {
		var (
			res struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}
		)
		start, end := 40*i, 40*(i+1)
		if end > len(msgBlock) {
			end = len(msgBlock)
		}
		body.Text = &ut.TextContent{
			Content: strings.Join(msgBlock[start:end], "\n") + fmt.Sprintf("\n(%d/%d)", i+1, num+1),
		}
		if err = json.NewEncoder(b).Encode(body); err != nil {
			log.Error("apmSvc.SendWechatToGroup Error(%v)", err)
			return
		}
		if req, err = http.NewRequest(http.MethodPost, url, b); err != nil {
			log.Error("apmSvc.SendWechatToGroup Error(%v)", err)
			return
		}
		if err = d.client.Do(c, req, &res); err != nil {
			log.Error("apmSvc.SendWechatToGroup Error(%v)", err)
			return
		}
		if res.Code != 0 {
			err = fmt.Errorf("Http response Code(%v)!=0", res.Code)
			log.Error("apmSvc.SendWechatToGroup Error(%v)", err)
			return
		}
	}
	return
}

// GitLabCommits transfer gitlab uri,now support get method
func (d *Dao) GitLabCommits(c context.Context, commitID string) (commit *ut.GitlabCommit, err error) {
	var req *http.Request
	params := url.Values{}
	params.Set("private_token", conf.Conf.Gitlab.Token)
	if req, err = http.NewRequest(http.MethodGet, _gitCommitsAPI+"/"+commitID+"?"+params.Encode(), nil); err != nil {
		log.Error("GitLabCommits http.NewRequest error(%v) params(%s)", err, params.Encode())
		return
	}
	if err = d.client.Do(c, req, &commit); err != nil {
		log.Error("GitLabCommits d.client.Do error(%v) params(%s)", err, params.Encode())
		return
	}
	return
}

// GetCoverage get the none-file coverage by commitID and pkg (pkg cannot be fileName)
func (d *Dao) GetCoverage(c context.Context, commitID, pkg string) (cov float64, err error) {
	var (
		count = strings.Count(pkg, "/")
		file  = &ut.File{}
	)
	if len(pkg) == 0 {
		err = fmt.Errorf("The pkg should not be empty")
		return
	}
	if pkg[len(pkg)-1] != '/' {
		count++
	}
	err = d.DB.Select(`commit_id,substring_index(name,"/",?) as pkg,round(sum(covered_statements)/sum(statements)*100,2) as coverage`, count).Group(fmt.Sprintf(`commit_id,substring_index(name,"/",%d)`, count)).Having("commit_id=? and pkg=?", commitID, pkg).First(file).Error
	if err == gorm.ErrRecordNotFound {
		cov, err = 0.00, nil
		return
	} else if err != nil {
		log.Error("dao.GetCoverage commitID(%s) pkg(%s) error(%v)", commitID, pkg, err)
		return
	}
	cov = file.Coverage
	return
}
