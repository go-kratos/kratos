package dao

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"go-common/app/admin/ep/marthe/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_buglyOkCode            = 200
	_buglyCookieExpiredCode = 100006
	_issueDetailCode        = 100000

	_issueDetailList    = "/v2/lastCrashInfo/appId/%s/platformId/%s/issues/%s/crashDataType/null?offsetTop=56&fsn=6d0260aa-331f-48b9-8557-c2aaf6e0be90"
	_issueList          = "/v2/issueList?sortOrder=desc&sortField=uploadTime&rows=50&fsn=45cdb5aa-eb6f-4bda-9bba-ba0b264bfc93&appId=%s&platformId=%s&version=%s&start=%s&rows=%s&exceptionTypeList=%s"
	_issueVersionList   = "/v2/getSelector/appId/%s/platformId/%s?types=version&fsn=8b8782b5-053d-4f58-bc17-d5c43d7f5ece"
	_issueExceptionList = "/v2/issueInfo/appId/%s/platformId/%s/issueId/%s/exceptionTypeList/Crash,Native,ExtensionCrash?fsn=114a8d02-586d-4fe4-8c23-79003fbe6882"
)

// BuglyVersion Bugly Version .
func (d *Dao) BuglyVersion(c context.Context, buglyCookie *model.BuglyCookie, projectID, platformID string) (ret []*model.BugVersion, err error) {
	var (
		req     *http.Request
		res     *model.BugVersionResponse
		hostStr string
	)

	hostStr = d.c.Bugly.Host + fmt.Sprintf(_issueVersionList, projectID, platformID)

	if req, err = d.newRequest("GET", hostStr, nil); err != nil {
		return
	}

	req.Header.Set("Cookie", buglyCookie.Cookie)
	req.Header.Set("x-token", buglyCookie.Token)
	req.Header.Set("content-type", "application/json;charset=utf-8")
	req.Header.Set("x-csrf-token", "undefined")

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.BugVersion url(%s) err(%v)", "BugVersion", err)
		return
	}

	if res.Status != _buglyOkCode {
		if res.Code == _buglyCookieExpiredCode {
			err = ecode.MartheCookieExpired
			log.Error("maybe need to update cookie and token")
		} else {
			err = ecode.MartheBuglyErr
		}
		log.Error("Status url(%s) res(%v) err(%v)", "BugVersion", res, err)
		return
	}

	ret = res.Ret.BugVersionList
	return
}

// BuglyIssueAndRetry Bugly Issue And Retry.
func (d *Dao) BuglyIssueAndRetry(c context.Context, buglyCookie *model.BuglyCookie, bugIssueRequest *model.BugIssueRequest) (ret *model.BugRet, err error) {
	for i := 0; i < d.c.Bugly.UrlRetryCount; i++ {
		if ret, err = d.BuglyIssue(c, buglyCookie, bugIssueRequest); err == nil {
			break
		}
	}
	return
}

// BuglyIssue Get Issue.
func (d *Dao) BuglyIssue(c context.Context, buglyCookie *model.BuglyCookie, bugIssueRequest *model.BugIssueRequest) (ret *model.BugRet, err error) {
	var (
		req     *http.Request
		res     *model.BugIssueResponse
		hostStr string
	)

	hostStr = d.c.Bugly.Host + fmt.Sprintf(_issueList, bugIssueRequest.ProjectID, bugIssueRequest.PlatformID, bugIssueRequest.Version, strconv.Itoa(bugIssueRequest.StartNum), strconv.Itoa(bugIssueRequest.Rows), bugIssueRequest.ExceptionType)

	if req, err = d.newRequest("GET", hostStr, nil); err != nil {
		return
	}

	req.Header.Set("Cookie", buglyCookie.Cookie)
	req.Header.Set("x-token", buglyCookie.Token)
	req.Header.Set("content-type", "application/json;charset=utf-8")
	req.Header.Set("x-csrf-token", "undefined")

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.BuglyIssue url(%s) err(%v)", "BuglyIssue", err)
		return
	}

	if res.Status != _buglyOkCode {
		if res.Code == _buglyCookieExpiredCode {
			err = ecode.MartheCookieExpired
			log.Error("maybe need to update cookie and token")
		} else {
			err = ecode.MartheBuglyErr
		}
		log.Error("Status url(%s) res(%v) err(%v)", "BuglyIssue", res, err)
		return
	}

	ret = res.Ret
	return
}

// BuglyIssueDetailAndRetry Bugly Issue Detail And Retry.
func (d *Dao) BuglyIssueDetailAndRetry(c context.Context, buglyCookie *model.BuglyCookie, projectID, platformID, issueNo string) (bugIssueDetail *model.BugIssueDetail, err error) {
	for i := 0; i < d.c.Bugly.UrlRetryCount; i++ {
		if bugIssueDetail, err = d.BuglyIssueDetail(c, buglyCookie, projectID, platformID, issueNo); err == nil {
			break
		}
	}
	return
}

// BuglyIssueDetail Get Issue Detail.
func (d *Dao) BuglyIssueDetail(c context.Context, buglyCookie *model.BuglyCookie, projectID, platformID, issueNo string) (bugIssueDetail *model.BugIssueDetail, err error) {
	var (
		req     *http.Request
		res     *model.BugIssueDetailResponse
		hostStr string
	)

	hostStr = d.c.Bugly.Host + fmt.Sprintf(_issueDetailList, projectID, platformID, issueNo)

	if req, err = d.newRequest("GET", hostStr, nil); err != nil {
		return
	}

	req.Header.Set("Cookie", buglyCookie.Cookie)
	req.Header.Set("x-token", buglyCookie.Token)
	req.Header.Set("content-type", "application/json;charset=utf-8")
	req.Header.Set("x-csrf-token", "undefined")

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.BuglyIssue url(%s) err(%v)", "BuglyIssue", err)
		return
	}

	if res.Code != _issueDetailCode {
		err = ecode.MartheBuglyErr
		log.Error("Status url(%s) res(%v) err(%v)", "BuglyIssue", res, err)
		return
	}

	bugIssueDetail = res.Data
	return
}

// BuglyIssueExceptionList Bugly Issue Exception List.
func (d *Dao) BuglyIssueExceptionList(c context.Context, buglyCookie *model.BuglyCookie, projectID, platformID, issueNo string) (bugIssueException *model.IssueException, err error) {
	var (
		req     *http.Request
		res     *model.BugIssueExceptionListResponse
		hostStr string
	)

	hostStr = d.c.Bugly.Host + fmt.Sprintf(_issueExceptionList, projectID, platformID, issueNo)

	if req, err = d.newRequest("GET", hostStr, nil); err != nil {
		return
	}

	req.Header.Set("Cookie", buglyCookie.Cookie)
	req.Header.Set("x-token", buglyCookie.Token)
	req.Header.Set("content-type", "application/json;charset=utf-8")
	req.Header.Set("x-csrf-token", "undefined")

	if err = d.httpClient.Do(c, req, &res); err != nil {
		log.Error("d.BuglyIssueExceptionList url(%s) err(%v)", "BuglyIssueExceptionList", err)
		return
	}

	if res.Status != _buglyOkCode {
		err = ecode.MartheBuglyErr
		log.Error("Status url(%s) res(%v) err(%v)", "BuglyIssueExceptionList", res, err)
		return
	}

	if res.Ret != nil && len(res.Ret.IssueException) != 0 && res.Ret.IssueException[0].IssueID == issueNo {
		bugIssueException = res.Ret.IssueException[0]
	}

	return
}
