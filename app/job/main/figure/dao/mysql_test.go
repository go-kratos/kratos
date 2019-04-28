package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/job/main/figure/model"
)

func TestSaveFigure(t *testing.T) {
	once.Do(startService)
	f := &model.Figure{
		Mid:             761223,
		Score:           100,
		LawfulScore:     d.c.Figure.Lawful,
		WideScore:       d.c.Figure.Wide,
		FriendlyScore:   d.c.Figure.Friendly,
		BountyScore:     d.c.Figure.Bounty,
		CreativityScore: d.c.Figure.Creativity,
		Ver:             1,
		Ctime:           time.Now(),
		Mtime:           time.Now(),
	}
	if id, err := d.SaveFigure(context.TODO(), f); err != nil {
		t.Errorf("figure err (%v)", err)
	} else {
		t.Logf("id(%d)", id)
	}
}

func TestGetFigure(t *testing.T) {
	once.Do(startService)
	var mid int64 = 7593623
	id, _ := d.ExistFigure(context.TODO(), mid)
	t.Logf("id:%d", id)
}
