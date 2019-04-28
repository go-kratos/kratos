package orm

import (
	. "github.com/jinzhu/gorm"
)

func init() {
	DefaultCallback.Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	DefaultCallback.Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
}

// updateTimeStampForCreateCallback will set `ctime`, `mtime` when creating
func updateTimeStampForCreateCallback(scope *Scope) {
	if !scope.HasError() {
		now := NowFunc()

		if createdAtField, ok := scope.FieldByName("ctime"); ok {
			if createdAtField.IsBlank {
				createdAtField.Set(now)
			}
		}

		if updatedAtField, ok := scope.FieldByName("mtime"); ok {
			if updatedAtField.IsBlank {
				updatedAtField.Set(now)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("mtime", NowFunc())
	}
}
