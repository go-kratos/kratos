package gorm_test

import (
	"context"
	"testing"
)

type PointerStruct struct {
	ID   int64
	Name *string
	Num  *int
}

type NormalStruct struct {
	ID   int64
	Name string
	Num  int
}

func TestPointerFields(t *testing.T) {
	DB.DropTable(&PointerStruct{})
	DB.AutoMigrate(&PointerStruct{})
	var name = "pointer struct 1"
	var num = 100
	pointerStruct := PointerStruct{Name: &name, Num: &num}
	if DB.Create(context.Background(), &pointerStruct).Error != nil {
		t.Errorf("Failed to save pointer struct")
	}

	var pointerStructResult PointerStruct
	if err := DB.First(context.Background(), &pointerStructResult, "id = ?", pointerStruct.ID).Error; err != nil || *pointerStructResult.Name != name || *pointerStructResult.Num != num {
		t.Errorf("Failed to query saved pointer struct")
	}

	var tableName = DB.NewScope(&PointerStruct{}).TableName()

	var normalStruct NormalStruct
	DB.Table(tableName).First(context.Background(), &normalStruct)
	if normalStruct.Name != name || normalStruct.Num != num {
		t.Errorf("Failed to query saved Normal struct")
	}

	var nilPointerStruct = PointerStruct{}
	if err := DB.Create(context.Background(), &nilPointerStruct).Error; err != nil {
		t.Error("Failed to save nil pointer struct", err)
	}

	var pointerStruct2 PointerStruct
	if err := DB.First(context.Background(), &pointerStruct2, "id = ?", nilPointerStruct.ID).Error; err != nil {
		t.Error("Failed to query saved nil pointer struct", err)
	}

	var normalStruct2 NormalStruct
	if err := DB.Table(tableName).First(context.Background(), &normalStruct2, "id = ?", nilPointerStruct.ID).Error; err != nil {
		t.Error("Failed to query saved nil pointer struct", err)
	}

	var partialNilPointerStruct1 = PointerStruct{Num: &num}
	if err := DB.Create(context.Background(), &partialNilPointerStruct1).Error; err != nil {
		t.Error("Failed to save partial nil pointer struct", err)
	}

	var pointerStruct3 PointerStruct
	if err := DB.First(context.Background(), &pointerStruct3, "id = ?", partialNilPointerStruct1.ID).Error; err != nil || *pointerStruct3.Num != num {
		t.Error("Failed to query saved partial nil pointer struct", err)
	}

	var normalStruct3 NormalStruct
	if err := DB.Table(tableName).First(context.Background(), &normalStruct3, "id = ?", partialNilPointerStruct1.ID).Error; err != nil || normalStruct3.Num != num {
		t.Error("Failed to query saved partial pointer struct", err)
	}

	var partialNilPointerStruct2 = PointerStruct{Name: &name}
	if err := DB.Create(context.Background(), &partialNilPointerStruct2).Error; err != nil {
		t.Error("Failed to save partial nil pointer struct", err)
	}

	var pointerStruct4 PointerStruct
	if err := DB.First(context.Background(), &pointerStruct4, "id = ?", partialNilPointerStruct2.ID).Error; err != nil || *pointerStruct4.Name != name {
		t.Error("Failed to query saved partial nil pointer struct", err)
	}

	var normalStruct4 NormalStruct
	if err := DB.Table(tableName).First(context.Background(), &normalStruct4, "id = ?", partialNilPointerStruct2.ID).Error; err != nil || normalStruct4.Name != name {
		t.Error("Failed to query saved partial pointer struct", err)
	}
}
