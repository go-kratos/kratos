package dao

import (
	"context"
	"errors"
	"fmt"

	"go-common/app/service/bbq/video/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/queue/databus"

	"github.com/json-iterator/go"
)

const (
	_queryCmsVideoRepository  = "select id from `video_repository` where `cid` = ?;"
	_insertCmsVideoRepository = "insert into `video_repository`(`avid`, `cid`, `mid`, `title`, `from`, `content`, `pubtime`, `duration`, `state`, `tid`, `sub_tid`, `svid`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);"
)

// ArchiveSub .
func (d *Dao) ArchiveSub() (*model.Archive, error) {
	if d.archiveSub == nil {
		return nil, ecode.ArchiveDatabusNilErr
	}
	msg, ok := <-d.archiveSub.Messages()
	if !ok {
		return nil, errors.New("chan <- failed")
	}
	return d.archiveProcess(msg)
}

func (d *Dao) archiveProcess(msg *databus.Message) (a *model.Archive, err error) {
	defer msg.Commit()

	an := new(model.ArchiveNotify)
	if err = jsoniter.Unmarshal(msg.Value, an); err != nil {
		return
	}

	if an.Action == "update" && an.New.Videos <= an.Old.Videos {
		return
	}

	if an.Table != "archive" || an.New == nil {
		return
	}

	if d.archiveFilters.DoFilter(an.New) {
		return
	}

	a = an.New

	return
}

// ArchiveKickOff .
func (d *Dao) ArchiveKickOff(c context.Context, svid int64, a *model.Archive) (err error) {
	row := d.cmsdb.QueryRow(c, _queryCmsVideoRepository, a.CID)
	tmp := 0
	if err = row.Scan(&tmp); err != nil && err != xsql.ErrNoRows {
		return
	}
	if tmp != 0 {
		err = fmt.Errorf("ArchiveKickOff cid existed [%d]", a.CID)
		return
	}
	_, err = d.cmsdb.Exec(c, _insertCmsVideoRepository, a.AID, a.CID, a.MID, a.Title, 0, a.Content, a.PubTime, a.Duration, a.State, a.TID, a.SubTID, svid)
	return
}
