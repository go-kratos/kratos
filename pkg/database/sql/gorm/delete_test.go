package gorm_test

import (
	"context"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	user1, user2 := User{Name: "delete1"}, User{Name: "delete2"}
	DB.Save(context.Background(), &user1)
	DB.Save(context.Background(), &user2)

	if err := DB.Delete(context.Background(), &user1).Error; err != nil {
		t.Errorf("No error should happen when delete a record, err=%s", err)
	}

	if !DB.Where("name = ?", user1.Name).First(context.Background(), &User{}).RecordNotFound() {
		t.Errorf("User can't be found after delete")
	}

	if DB.Where("name = ?", user2.Name).First(context.Background(), &User{}).RecordNotFound() {
		t.Errorf("Other users that not deleted should be found-able")
	}
}

func TestInlineDelete(t *testing.T) {
	user1, user2 := User{Name: "inline_delete1"}, User{Name: "inline_delete2"}
	DB.Save(context.Background(), &user1)
	DB.Save(context.Background(), &user2)

	if DB.Delete(context.Background(), &User{}, user1.Id).Error != nil {
		t.Errorf("No error should happen when delete a record")
	} else if !DB.Where("name = ?", user1.Name).First(context.Background(), &User{}).RecordNotFound() {
		t.Errorf("User can't be found after delete")
	}

	if err := DB.Delete(context.Background(), &User{}, "name = ?", user2.Name).Error; err != nil {
		t.Errorf("No error should happen when delete a record, err=%s", err)
	} else if !DB.Where("name = ?", user2.Name).First(context.Background(), &User{}).RecordNotFound() {
		t.Errorf("User can't be found after delete")
	}
}

func TestSoftDelete(t *testing.T) {
	type User struct {
		Id        int64
		Name      string
		DeletedAt *time.Time
	}
	DB.AutoMigrate(&User{})

	user := User{Name: "soft_delete"}
	DB.Save(context.Background(), &user)
	DB.Delete(context.Background(), &user)

	if DB.First(context.Background(), &User{}, "name = ?", user.Name).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	if err := DB.Unscoped().First(context.Background(), &User{}, "name = ?", user.Name).Error; err != nil {
		t.Errorf("Should be able to find soft deleted record with Unscoped, but err=%s", err)
	}

	DB.Unscoped().Delete(context.Background(), &user)
	if !DB.Unscoped().First(context.Background(), &User{}, "name = ?", user.Name).RecordNotFound() {
		t.Errorf("Can't find permanently deleted record")
	}
}

func TestSoftDeleteWithCustomizedDeletedAtColumnName(t *testing.T) {
	creditCard := CreditCard{Number: "411111111234567"}
	DB.Save(context.Background(), &creditCard)
	DB.Delete(context.Background(), &creditCard)

	if deletedAtField, ok := DB.NewScope(&CreditCard{}).FieldByName("DeletedAt"); !ok || deletedAtField.DBName != "deleted_time" {
		t.Errorf("CreditCard's DeletedAt's column name should be `deleted_time`")
	}

	if DB.First(context.Background(), &CreditCard{}, "number = ?", creditCard.Number).Error == nil {
		t.Errorf("Can't find a soft deleted record")
	}

	if err := DB.Unscoped().First(context.Background(), &CreditCard{}, "number = ?", creditCard.Number).Error; err != nil {
		t.Errorf("Should be able to find soft deleted record with Unscoped, but err=%s", err)
	}

	DB.Unscoped().Delete(context.Background(), &creditCard)
	if !DB.Unscoped().First(context.Background(), &CreditCard{}, "number = ?", creditCard.Number).RecordNotFound() {
		t.Errorf("Can't find permanently deleted record")
	}
}
