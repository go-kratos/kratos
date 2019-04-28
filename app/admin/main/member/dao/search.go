package dao

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/library/database/elastic"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	_numPattern = regexp.MustCompile(`^\d+$`)
)

func allnum(s string) bool {
	return _numPattern.MatchString(s)
}

// SearchMember is.
func (d *Dao) SearchMember(ctx context.Context, arg *model.ArgList) (*model.SearchMemberResult, error) {
	r := d.es.NewRequest("member_user").
		Fields("mid", "name").
		Index("user_base").
		Order("mid", elastic.OrderAsc).
		Ps(int(arg.PS)).
		Pn(int(arg.PN))
	if arg.Mid != 0 {
		r.WhereEq("mid", arg.Mid)
	}
	if arg.Keyword != "" {
		fields := []string{"name"}
		if allnum(arg.Keyword) {
			fields = append(fields, "mid")
		}
		r.WhereLike(fields, []string{arg.Keyword}, true, elastic.LikeLevelLow)
	}
	result := &model.SearchMemberResult{}
	if err := r.Scan(ctx, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}

// SearchLog is.
func (d *Dao) SearchLog(ctx context.Context, uid, oid int64, uname, action string) (*model.SearchLogResult, error) {
	nowYear := time.Now().Year()
	index := []string{
		fmt.Sprintf("log_audit_121_%d", nowYear),
		fmt.Sprintf("log_audit_121_%d", nowYear-1),
	}
	r := d.es.NewRequest("log_audit").
		Index(index...).
		Fields("uid", "uname", "oid", "type", "action", "str_0", "str_1", "str_2", "int_0", "int_1", "int_2", "ctime", "extra_data").
		Order("ctime", elastic.OrderDesc)
	// 默认查询的是第一页，每夜10条记录
	if uid > 0 {
		r.WhereEq("uid", strconv.FormatInt(uid, 10))
	}
	if oid > 0 {
		r.WhereEq("oid", strconv.FormatInt(oid, 10))
	}
	if uname != "" {
		r.WhereEq("uname", uname)
	}
	if action != "" {
		r.WhereIn("action", strings.Split(action, ","))
	}

	result := &model.SearchLogResult{}
	if err := r.Scan(ctx, result); err != nil {
		log.Error("Failed to SearchLog: Scan params(%s) error(%+v)", r.Params(), err)
		return nil, err
	}
	return result, nil
}

// SearchFaceCheckRes is.
func (d *Dao) SearchFaceCheckRes(ctx context.Context, fileName string) (*model.SearchLogResult, error) {
	index := fmt.Sprintf("log_audit_161_%s", time.Now().Format("2006_01"))
	r := d.es.NewRequest("log_audit").
		Index(index).
		Fields("uid", "uname", "oid", "type", "action", "str_0", "str_1", "str_2", "int_0", "int_1", "int_2", "ctime", "extra_data").
		WhereEq("str_1", "face").
		WhereEq("str_2", fileName).
		Order("ctime", elastic.OrderDesc).
		Pn(1).
		Ps(1)

	result := &model.SearchLogResult{}
	if err := r.Scan(ctx, result); err != nil {
		log.Error("Failed to SearchFaceCheckRes: Scan params(%s) error(%+v)", r.Params(), err)
		return nil, err
	}
	return result, nil
}

// SearchUserAuditLog is.
func (d *Dao) SearchUserAuditLog(ctx context.Context, mid int64) (*model.SearchLogResult, error) {
	nowYear := time.Now().Year()
	index := []string{
		fmt.Sprintf("log_audit_121_%d", nowYear),
		fmt.Sprintf("log_audit_121_%d", nowYear-1),
	}
	r := d.es.NewRequest("log_audit").
		Index(index...).
		Fields("uid", "uname", "oid", "type", "action", "str_0", "str_1", "str_2", "int_0", "int_1", "int_2", "ctime", "extra_data").
		WhereEq("action", "base_audit").
		WhereEq("oid", strconv.FormatInt(mid, 10)).
		Order("ctime", elastic.OrderDesc).
		Pn(1).
		Ps(5)

	result := &model.SearchLogResult{}
	if err := r.Scan(ctx, result); err != nil {
		log.Error("Failed to SearchUserAuditLog: Scan params(%s) error(%+v)", r.Params(), err)
		return nil, err
	}
	return result, nil
}

// SearchUserPropertyReview is.
func (d *Dao) SearchUserPropertyReview(ctx context.Context, mid int64, property []int, state []int, isMonitor, isDesc bool, operator, stime, etime string, pn, ps int) (*model.SearchUserPropertyReviewResult, error) {
	order := elastic.OrderAsc
	if isDesc {
		order = elastic.OrderDesc
	}
	monitor := 0
	if isMonitor {
		monitor = 1
	}
	r := d.es.NewRequest("user_property_review").
		Fields("id").
		Index("user_property_review").
		WhereEq("is_monitor", monitor).
		WhereRange("ctime", stime, etime, elastic.RangeScopeLcRc).
		Order("id", order).
		Ps(ps).
		Pn(pn)
	if mid > 0 {
		r.WhereEq("mid", mid)
	}
	if operator != "" {
		r.WhereEq("operator", operator)
	}
	if len(property) > 0 {
		r.WhereIn("property", property)
	}
	if len(state) > 0 {
		r.WhereIn("state", state)
	}
	result := &model.SearchUserPropertyReviewResult{}
	if err := r.Scan(ctx, result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result, nil
}
