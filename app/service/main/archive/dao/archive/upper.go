package archive

import (
	"context"
	"fmt"

	"go-common/library/log"
	"go-common/library/time"
	"go-common/library/xstr"

	"github.com/pkg/errors"
)

const (
	_upCntSQL       = "SELECT COUNT(*) FROM archive WHERE mid=? AND (state>=0 or state=-6)"
	_upsCntSQL      = "SELECT mid,COUNT(*) FROM archive WHERE mid IN(%s) AND (state>=0 or state=-6) GROUP BY mid"
	_upPasSQL       = "SELECT aid,pubtime,copyright FROM archive WHERE mid=? AND state>=0 ORDER BY pubtime DESC"
	_upsPasSQL      = "SELECT aid,mid,pubtime,copyright FROM archive WHERE mid IN (%s) AND state>=0 ORDER BY pubtime DESC"
	_upRecommendSQL = "SELECT reco_aid FROM archive_recommend WHERE aid=? AND state=0"
)

// UppersCount get mids count
func (d *Dao) UppersCount(c context.Context, mids []int64) (uc map[int64]int, err error) {
	rows, err := d.resultDB.Query(c, fmt.Sprintf(_upsCntSQL, xstr.JoinInts(mids)))
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	uc = make(map[int64]int, len(mids))
	defer rows.Close()
	for rows.Next() {
		var (
			mid int64
			cnt int
		)
		if err = rows.Scan(&mid, &cnt); err != nil {
			err = errors.WithStack(err)
			return
		}
		uc[mid] = cnt
	}
	return
}

// UpperCount get the count of archives by mid of Up.
func (d *Dao) UpperCount(c context.Context, mid int64) (count int, err error) {
	d.infoProm.Incr("UpperCount")
	row := d.upCntStmt.QueryRow(c, mid)
	if err = row.Scan(&count); err != nil {
		log.Error("row.Scan error(%v)", err)
	}
	return
}

// UpperPassed get upper passed archives.
func (d *Dao) UpperPassed(c context.Context, mid int64) (aids []int64, ptimes []time.Time, copyrights []int8, err error) {
	d.infoProm.Incr("UpperPassed")
	rows, err := d.upPasStmt.Query(c, mid)
	if err != nil {
		log.Error("getUpPasStmt.Query(%d) error(%v)", mid, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			aid       int64
			ptime     time.Time
			copyright int8
		)
		if err = rows.Scan(&aid, &ptime, &copyright); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
		ptimes = append(ptimes, ptime)
		copyrights = append(copyrights, copyright)
	}
	return
}

// UppersPassed get uppers passed archives.
func (d *Dao) UppersPassed(c context.Context, mids []int64) (aidm map[int64][]int64, ptimes map[int64][]time.Time, copyrights map[int64][]int8, err error) {
	d.infoProm.Incr("UppersPassed")
	rows, err := d.resultDB.Query(c, fmt.Sprintf(_upsPasSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("UpsPassed error(%v)", err)
		return
	}
	defer rows.Close()
	aidm = make(map[int64][]int64, len(mids))
	ptimes = make(map[int64][]time.Time, len(mids))
	copyrights = make(map[int64][]int8, len(mids))
	for rows.Next() {
		var (
			aid, mid  int64
			ptime     time.Time
			copyright int8
		)
		if err = rows.Scan(&aid, &mid, &ptime, &copyright); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		aidm[mid] = append(aidm[mid], aid)
		ptimes[mid] = append(ptimes[mid], ptime)
		copyrights[mid] = append(copyrights[mid], copyright)
	}
	return
}

// UpperReommend get up recommend
func (d *Dao) UpperReommend(c context.Context, aid int64) (reAids []int64, err error) {
	d.infoProm.Incr("UpperRecommend")
	rows, err := d.arcReadDB.Query(c, _upRecommendSQL, aid)
	if err != nil {
		log.Error("d.arcReadDB.Query(%s, %d) error(%v)", _upRecommendSQL, aid)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var reAid int64
		if err = rows.Scan(&reAid); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		reAids = append(reAids, reAid)
	}
	return
}
