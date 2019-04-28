package notification

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"go-common/app/tool/saga/dao"
	"go-common/app/tool/saga/model"
	"go-common/app/tool/saga/service/mail"
	"go-common/app/tool/saga/service/wechat"
	"go-common/library/log"
)

// WechatAuthor send wechat message to original author
func WechatAuthor(dao *dao.Dao, authorName string, url, sourceBranch, targetBranch string, comment string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("wechatAuthor: %+v %s", x, debug.Stack())
		}
	}()
	var (
		subject = fmt.Sprintf("[SAGA] MR ( %s ) merge 通知", sourceBranch)
		data    = fmt.Sprintf("MR : %s \n ( %s -> %s ) 状态 %s", url, sourceBranch, targetBranch, comment)
		wct     = wechat.New(dao)
		ctx     = context.Background()
	)

	return wct.PushMsg(ctx, []string{authorName}, fmt.Sprintf("%s\n\n%s", subject, data))
}

// MailPipeline ...
func MailPipeline(event *model.HookPipeline) (err error) {
	var (
		author = event.User.UserName
		branch = event.ObjectAttributes.Ref
		commit = event.ObjectAttributes.Sha
		title  = ""
		url    = ""
		status = "失败"
	)
	if event.Commit != nil {
		commitIndex := strings.LastIndex(event.Commit.URL, "commit")
		url = event.Commit.URL[:commitIndex] + "pipelines/" + strconv.FormatInt(event.ObjectAttributes.ID, 10)
	}

	if strings.Contains(event.Commit.Message, "\n") {
		title = event.Commit.Message[:strings.Index(event.Commit.Message, "\n")]
	} else {
		title = event.Commit.Message
	}

	if event.ObjectAttributes.Status == model.PipelineSuccess {
		status = "成功"
	}

	subject := fmt.Sprintf("[SAGA] Pipeline ( %s ) %s 通知", branch, status)
	data := fmt.Sprintf("Pipeline : %s \nCommit : %s\nCommitID : %s\n状态: 运行%s !", url, title, commit, status)
	addr := &model.MailAddress{
		Address: author + "@bilibili.com",
		Name:    author,
	}

	if err = mail.SendMail2(addr, subject, data); err != nil {
		log.Error("%+v", err)
	}
	return
}

// WechatPipeline ...
func WechatPipeline(dao *dao.Dao, event *model.HookPipeline) (err error) {
	var (
		wct    = wechat.New(dao)
		ctx    = context.Background()
		author = event.User.UserName
		branch = event.ObjectAttributes.Ref
		commit = event.ObjectAttributes.Sha
		title  = ""
		url    = ""
		status = "失败"
	)
	if event.Commit != nil {
		commitIndex := strings.LastIndex(event.Commit.URL, "commit")
		url = event.Commit.URL[:commitIndex] + "pipelines/" + strconv.FormatInt(event.ObjectAttributes.ID, 10)
	}

	if strings.Contains(event.Commit.Message, "\n") {
		title = event.Commit.Message[:strings.Index(event.Commit.Message, "\n")]
	} else {
		title = event.Commit.Message
	}

	if event.ObjectAttributes.Status == model.PipelineSuccess {
		status = "成功"
	}

	subject := fmt.Sprintf("[SAGA] Pipeline ( %s ) %s 通知", branch, status)
	data := fmt.Sprintf("Pipeline : %s \nCommit : %s\nCommitID : %s\n状态: 运行%s !", url, title, commit, status)

	if err = wct.PushMsg(ctx, []string{author}, fmt.Sprintf("%s\n\n%s", subject, data)); err != nil {
		log.Error("%+v", err)
	}
	return
}

// MailPlusOne ...
func MailPlusOne(author, reviewer, url, commit, sourceBranch, targetBranch string) (err error) {
	var (
		subject = fmt.Sprintf("[SAGA] MR ( %s ) review 通知", sourceBranch)
		data    = fmt.Sprintf("MR : %s \n ( %s -> %s ) commitID ( %s ) 已被 ( %s ) review!", url, sourceBranch, targetBranch, commit, reviewer)
		addr    = &model.MailAddress{
			Address: author + "@bilibili.com",
			Name:    author,
		}
	)
	if err = mail.SendMail2(addr, subject, data); err != nil {
		log.Error("%+v", err)
	}
	return
}

// WechatPlusOne ...
func WechatPlusOne(dao *dao.Dao, author, reviewer, url, commit, sourceBranch, targetBranch string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("wechatPlusOne: %+v %s", x, debug.Stack())
		}
	}()
	var (
		subject = fmt.Sprintf("[SAGA] MR ( %s ) review 通知", sourceBranch)
		data    = fmt.Sprintf("MR : %s \n ( %s -> %s ) commitID ( %s ) 已被 ( %s ) review!", url, sourceBranch, targetBranch, commit, reviewer)
		wct     = wechat.New(dao)
		ctx     = context.Background()
	)

	if err = wct.PushMsg(ctx, []string{author}, fmt.Sprintf("%s\n\n%s", subject, data)); err != nil {
		log.Error("%+v", err)
	}
	return
}

// func notifyAssign(assign string, authorName string, url, sourceBranch, targetBranch string) (err error) {
// 	var (
// 		subject = fmt.Sprintf("[SAGA] MR ( %s ) assign 通知", sourceBranch)
// 		data    = fmt.Sprintf("MR : %s \n ( %s -> %s ) 开发者 ( %s ) 已指派给您 !", url, sourceBranch, targetBranch, authorName)
// 		addr    = &model.MailAddress{
// 			Address: assign + "@bilibili.com",
// 			Name:    assign,
// 		}
// 	)
// 	if err = mail.SendMail2(addr, subject, data); err != nil {
// 		log.Error("%+v", err)
// 	}
// 	return
// }

// func wechatAssign(dao *dao.Dao, assign string, authorName string, url, sourceBranch, targetBranch string) (err error) {
// 	var (
// 		subject = fmt.Sprintf("[SAGA] MR ( %s ) assign 通知", sourceBranch)
// 		data    = fmt.Sprintf("MR : %s \n ( %s -> %s ) 开发者 ( %s ) 已指派给您 !", url, sourceBranch, targetBranch, authorName)
// 		wct     = wechat.New(dao)
// 		ctx     = context.Background()
// 	)

// 	if err = wct.PushMsg(ctx, []string{assign}, fmt.Sprintf("%s\n\n%s", subject, data)); err != nil {
// 		log.Error("wechatAssign failed, err (%s)", err.Error())
// 	}
// 	return
// }

// func notifyReviewer(reviewers []string, authorName, url, sourceBranch, targetBranch string) (err error) {
// 	var (
// 		subject = fmt.Sprintf("[SAGA] MR ( %s ) review 请求", sourceBranch)
// 		data    = fmt.Sprintf("MR : %s \n ( %s -> %s ) 开发者 ( %s ) 需要您的 review !\n请在 review 后留言 +1 !", url, sourceBranch, targetBranch, authorName)
// 	)
// 	for _, reviewer := range reviewers {
// 		addr := &model.MailAddress{
// 			Address: reviewer + "@bilibili.com",
// 			Name:    reviewer,
// 		}
// 		if err = mail.SendMail2(addr, subject, data); err != nil {
// 			log.Error("%+v", err)
// 		}
// 	}
// 	return
// }

// func wechatReviewer(dao *dao.Dao, reviewers []string, authorName, url, sourceBranch, targetBranch string) (err error) {
// 	var (
// 		subject = fmt.Sprintf("[SAGA] MR ( %s ) review 请求", sourceBranch)
// 		data    = fmt.Sprintf("MR : %s \n ( %s -> %s ) 开发者 ( %s ) 需要您的 review !\n请在 review 后留言 +1 !", url, sourceBranch, targetBranch, authorName)
// 		wct     = wechat.New(dao)
// 		ctx     = context.Background()
// 	)

// 	if err = wct.PushMsg(ctx, reviewers, fmt.Sprintf("%s\n\n%s", subject, data)); err != nil {
// 		log.Error("%+v", err)
// 	}
// 	return
// }
