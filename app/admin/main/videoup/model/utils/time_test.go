package utils

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_format(t *testing.T) {
	type A struct {
		Ct FormatTime `json:"ct"`
	}

	a := &A{Ct: ""}
	if !a.Ct.TimeValue().IsZero() {
		t.Fatal("1")
	}

	a = &A{Ct: "0000-00-00 00:00:00"}
	fmt.Println("ct:", a.Ct)
	fmt.Println("value:", a.Ct.TimeValue())
	fmt.Println("zero:", a.Ct.TimeValue().IsZero())
	if !a.Ct.TimeValue().IsZero() {
		t.Fatal("2")
	}

	a = &A{Ct: "0001-01-01 00:00:00"}
	fmt.Println("ct:", a.Ct)
	fmt.Println("value:", a.Ct.TimeValue())
	fmt.Println("zero:", a.Ct.TimeValue().IsZero())
	if !a.Ct.TimeValue().IsZero() {
		t.Fatal("3")
	}

	Tt, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-06-01 22:59:16", time.Local)
	a = &A{}
	J := `{"ct":"2018-06-01 22:59:16"}`
	err := json.Unmarshal([]byte(J), a)
	if err != nil || string(a.Ct) != "2018-06-01 22:59:16" || a.Ct.TimeValue() != Tt {
		t.Fatal("4")
	}
	fmt.Println("ct:", a.Ct)
	fmt.Println("value:", a.Ct.TimeValue())
	fmt.Println("zero:", a.Ct.TimeValue().IsZero())

	J = `{"ct":"2018-06-01T22:59:16.437367789+08:00"}`
	err = json.Unmarshal([]byte(J), a)
	if err != nil || string(a.Ct) != "2018-06-01 22:59:16" || a.Ct.TimeValue() != Tt {
		t.Fatal("5")
	}
	fmt.Println("ct:", a.Ct)
	fmt.Println("value:", a.Ct.TimeValue())
	fmt.Println("zero:", a.Ct.TimeValue().IsZero())
}
