package model

import (
	xtime "go-common/library/time"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRestrictDate(t *testing.T) {
	xt1 := xtime.Time(-62169580800)
	t1 := xt1.Time()
	assert.True(t, t1.Year() < 0)

	xt2 := RestrictDate(xt1)
	t2 := xt2.Time()
	assert.Equal(t, 0, t2.Year())
	assert.Equal(t, t1.Month(), t2.Month())
	assert.Equal(t, t1.Day(), t2.Day())
}

func TestRealname(t *testing.T) {
	realnameIMG := RealnameApplyImgMessage{
		AddTimeStr: "2018-05-11 15:07:00",
		AddTimeDB:  time.Now(),
	}

	t.Log("addtime", realnameIMG.AddTime().Local())

	realnameApply := RealnameApplyMessage{
		CardDataCanal: "dGMyNUhTcXN1NS9MYkF1bGdoKzNkTGs0eHI5UlM4SS9NY2VQa25xV3czU2grY1Q3R0JiSjlaQWhPZ294TTlEQVV0ZFhuYzIrUnlXTVN3NWk5THFWdkpmTEJaYXJhSlFHWUQ3bVlWZ3liNU1IS1hQTXZ0RE9pR1d6UnpYcUtRUEFYTHZCcXIzZVFoVURwT3VieXJzV0c1Z0dJS2dQNEdYbUV0T1B3MXV6bEE4PQ==",
	}

	t.Log("card data", realnameApply.CardData())

	var (
		card = "340702199110120012"
	)
	md5 := cardMD5(card, 0, 0)
	t.Log("card md5", md5)
}
