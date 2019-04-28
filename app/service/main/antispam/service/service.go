package service

import (
	"context"
	"errors"

	"go-common/app/service/main/antispam/model"
)

const (
	// CacheKeyRegexps .
	CacheKeyRegexps = "regexps"
	// CacheKeyLocalCountsFormat .
	CacheKeyLocalCountsFormat = "resource_id:%d:keyword_id:%d:local_limit_counts"
	// CacheKeyTotalCountsFormat .
	CacheKeyTotalCountsFormat = "keyword_id:%d:total_counts"
	// CacheKeyGlobalCountsFormat .
	CacheKeyGlobalCountsFormat = "keyword_id:%d:global_limit_counts"
	// CacheKeyKeywordsSenderIDsFormat .
	CacheKeyKeywordsSenderIDsFormat = "keyword_id:%d:sender_ids"
	// CacheKeyRulesWithinAreaAndLimitTypeFormat .
	CacheKeyRulesWithinAreaAndLimitTypeFormat = "rule:area:%s:limit_type:%s"
)

var (
	// ErrOverRateLimit .
	ErrOverRateLimit = errors.New("over rate limit")
	// ErrIllegalOperation .
	ErrIllegalOperation = errors.New("illegal operation")
	// ErrResourceNotExist .
	ErrResourceNotExist = errors.New("resource not found")
)

// Service .
type Service interface {
	Close()
	Ping(context.Context) error

	Filter(context.Context, UserGeneratedContent) (*model.SuspiciousResp, error)
	RefreshRules(context.Context)
	RefreshRegexps(context.Context)

	GetKeywordsByCond(context.Context, *Condition) ([]*model.Keyword, int64, error)
	GetKeywordsByOffsetLimit(context.Context, *Condition) ([]*model.Keyword, error)
	GetKeywordByID(ctx context.Context, id int64) (*model.Keyword, error)
	GetSenderIDsByKeywordID(ctx context.Context, id int64) (*model.SenderList, error)
	OpKeyword(ctx context.Context, id int64, newTag string) (*model.Keyword, error)
	DeleteKeywords(ctx context.Context, ids []int64) ([]*model.Keyword, error)

	GetRegexpByID(context.Context, int64) (*model.Regexp, error)
	GetRegexpsByCond(context.Context, *Condition) ([]*model.Regexp, int64, error)
	GetRegexpByAreaAndContent(ctx context.Context, area, content string) (*model.Regexp, error)
	UpsertRegexp(context.Context, *model.Regexp) (*model.Regexp, error)
	DeleteRegexp(ctx context.Context, id int64, adminID int64) (*model.Regexp, error)

	GetRuleByAreaAndLimitTypeAndScope(ctx context.Context, area, limitType, limitScope string) (*model.Rule, error)
	GetRuleByArea(ctx context.Context, area string) ([]*model.Rule, error)
	UpsertRule(context.Context, *model.Rule) (*model.Rule, error)
	ExpireKeyword(context.Context, int64) error
}
