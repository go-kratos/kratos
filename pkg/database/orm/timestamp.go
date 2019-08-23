package orm

import (
	"github.com/jinzhu/gorm"
)

// updateTimeStampForCreateCallback will set `ctime`, `mtime` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		now := gorm.NowFunc()

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

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("mtime", gorm.NowFunc())
	}
}
