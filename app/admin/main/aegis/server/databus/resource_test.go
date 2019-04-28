package databus

import (
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	InitAegis(nil)
	defer CloseAegis()
	err := Add(&AddInfo{
		BusinessID: 1,
		NetID:      1,
		OID:        "1008612",
		MID:        110,
		Content:    "内容",
		Extra6:     6,
		Extra4s:    "4s",
		MetaData:   `{"cover": "bfs/1.japg", "title": "标题啊"}`,
		ExtraTime1: time.Now().Add(-24 * time.Hour),
		OCtime:     time.Now().Add(-24 * time.Hour),
		Ptime:      time.Now().Add(-24 * time.Hour),
	})

	if err != nil {
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	InitAegis(nil)
	defer CloseAegis()
	err := Update(&UpdateInfo{
		BusinessID: 1,
		NetID:      1,
		OID:        "1008612",
		Update: map[string]interface{}{
			"mid":      119,
			"content":  "内容2",
			"extra6":   60,
			"ptime":    "2018-01-1 15:01:02",
			"metadata": `{"cover": "bfs/1.japg", "title": "标题啊"}`,
		},
	})

	if err != nil {
		t.Fail()
	}
}

func TestCancel(t *testing.T) {
	InitAegis(nil)
	defer CloseAegis()
	err := Cancel(&CancelInfo{
		BusinessID: 1,
		Oids:       []string{"1008612"},
		Reason:     "理由",
	})

	if err != nil {
		t.Fail()
	}
}
