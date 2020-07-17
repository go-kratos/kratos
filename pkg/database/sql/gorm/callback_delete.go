package gorm

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/pkg/net/trace"
)

// Define callbacks for deleting
func init() {
	DefaultCallback.Delete().Register("gorm:begin_transaction", beginTransactionCallback)
	DefaultCallback.Delete().Register("gorm:before_delete", beforeDeleteCallback)
	DefaultCallback.Delete().Register("gorm:delete", deleteCallback)
	DefaultCallback.Delete().Register("gorm:after_delete", afterDeleteCallback)
	DefaultCallback.Delete().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}

// beforeDeleteCallback will invoke `BeforeDelete` method before deleting
func beforeDeleteCallback(scope *Scope) {
	if scope.DB().HasBlockGlobalUpdate() && !scope.hasConditions() {
		scope.Err(errors.New("missing WHERE clause while deleting"))
		return
	}
	if !scope.HasError() {
		scope.CallMethod("BeforeDelete")
	}
}

// deleteCallback used to delete data from database or set deleted_at to current time (when using with soft delete)
func deleteCallback(scope *Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedAtField, hasDeletedAtField := scope.FieldByName("DeletedAt")

		var sql string
		if !scope.Search.Unscoped && hasDeletedAtField {
			sql = fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(scope.db.nowFunc()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)
		} else {
			sql = fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)
		}
		if scope.t != nil {
			scope.t = scope.t.Fork(_family, "sql exec")
			scope.t.SetTag(trace.String(trace.TagComment, sql))
			scope.t.SetTag(trace.String(trace.TagCommentExt, fmt.Sprintf("%+v", scope.SQLVars)))
		}
		scope.Raw(sql).Exec()

	}
}

// afterDeleteCallback will invoke `AfterDelete` method after deleting
func afterDeleteCallback(scope *Scope) {
	if !scope.HasError() {
		scope.CallMethod("AfterDelete")
		if scope.t != nil {
			scope.t.Finish(&scope.db.Error)
		}
	}
}
