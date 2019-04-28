package recommend

import (
	"context"
	"strconv"
	"time"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// PubDislike is.
func (d *Dao) PubDislike(c context.Context, buvid, gt string, id, mid, reasonID, cmreasonID, feedbackID, upperID, rid, tagID int64, adcb string, now time.Time) (err error) {
	return d.pub(c, buvid, gt, id, mid, reasonID, cmreasonID, feedbackID, upperID, rid, tagID, adcb, 1, now)
}

// PubDislikeCancel is.
func (d *Dao) PubDislikeCancel(c context.Context, buvid, gt string, id, mid, reasonID, cmreasonID, feedbackID, upperID, rid, tagID int64, adcb string, now time.Time) (err error) {
	return d.pub(c, buvid, gt, id, mid, reasonID, cmreasonID, feedbackID, upperID, rid, tagID, adcb, 2, now)
}

func (d *Dao) pub(c context.Context, buvid, gt string, id, mid, reasonID, cmreasonID, feedbackID, upperID, rid, tagID int64, adcb string, state int8, now time.Time) (err error) {
	key := strconv.FormatInt(mid, 10)
	msg := struct {
		Buvid      string `json:"buvid"`
		Goto       string `json:"goto"`
		ID         int64  `json:"id"`
		Mid        int64  `json:"mid"`
		ReasonID   int64  `json:"reason_id"`
		CMReasonID int64  `json:"cm_reason_id"`
		FeedbackID int64  `json:"feedback_id"`
		UpperID    int64  `json:"upper_id"`
		Rid        int64  `json:"rid"`
		TagID      int64  `json:"tag_id"`
		ADCB       string `json:"ad_cb"`
		State      int8   `json:"state"`
		Time       int64  `json:"time"`
	}{Buvid: buvid, Goto: gt, ID: id, Mid: mid, ReasonID: reasonID, CMReasonID: cmreasonID, FeedbackID: feedbackID, UpperID: upperID, Rid: rid, TagID: tagID, ADCB: adcb, State: state, Time: now.Unix()}
	if err = d.databus.Send(c, key, msg); err != nil {
		err = errors.Wrapf(err, "%s %v", key, msg)
		return
	}
	log.Info("d.dataBus.Pub(%s,%v)", key, msg)
	return
}
