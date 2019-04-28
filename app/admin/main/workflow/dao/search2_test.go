package dao

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/search"

	"github.com/smartystreets/goconvey/convey"
)

func TestSearchAuditLogGroup(t *testing.T) {
	convey.Convey("SearchAuditLogGroup", t, func() {
		cond := &search.AuditLogGroupSearchCond{
			Group:      []string{"oid"},
			Businesses: []int64{3},
			Order:      "ctime",
			PS:         1000,
			PN:         1,
			Sort:       "desc",
			//Oids:       []int64{24648472, 24694284, 24706061, 24583668, 24703363},
			Oids: []int64{10110151},
		}
		result, err := d.SearchAuditLogGroup(context.TODO(), cond)
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeEmpty)
		fmt.Println(len(result.Data.Result))
		fmt.Println(result.Data.Result)
	})
}

func TestSearchChallengeMultiPage(t *testing.T) {
	convey.Convey("SearchChallengeMultiPage", t, func() {
		cond := &search.ChallSearchCommonCond{
			Fields: []string{"id", "gid", "mid", "state", "ctime"},
			//IDs:    []int64{1869, 1872},
			IDs: []int64{1967, 1966, 1960},
		}
		result, err := d.SearchChallengeMultiPage(context.TODO(), cond)
		convey.So(err, convey.ShouldBeNil)
		fmt.Println("len", len(result))
		for _, c := range result {
			tc := &model.TinyChall{
				Cid: c.ID,
				Gid: c.Gid,
				Mid: c.Mid,
			}
			var t time.Time
			if t, err = time.Parse("2006-01-02 15:04:05", c.CTime); err != nil {
				fmt.Printf("time.Parse() error(%v)\n", err)
			}
			if err = tc.CTime.Scan(t); err != nil {
				fmt.Printf("xtime.Scan() error(%v)\n", err)
			}
			if str, ok := c.State.(string); ok {
				st, _ := strconv.Atoi(str)
				tc.State = int8(st)
			}
			if f, ok := c.State.(float64); ok {
				tc.State = int8(math.Floor(f))
			}
			fmt.Printf("%+v\n", tc)
		}
	})
}

func TestSearchWorkflowReportLog(t *testing.T) {
	convey.Convey("SearchWkf", t, func() {
		cond := &search.AuditReportSearchCond{
			AppID:    search.LogAuditAction,
			Fields:   []string{"uname", "extra_data"},
			Business: 11,
			Order:    "ctime",
			Sort:     "desc",
			UID:      []int64{6},
		}
		result, err := d.SearchAuditReportLog(context.TODO(), cond)
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeEmpty)
		fmt.Println(len(result.Result))
		fmt.Printf("%+v\n", result.Result[0])
	})
}

func TestSearchArchiveReportLog(t *testing.T) {
	convey.Convey("SearchArc", t, func() {
		cond := &search.AuditReportSearchCond{
			AppID:         search.LogAuditAction,
			Fields:        []string{"oid", "ctime"},
			IndexTimeType: "month",
			IndexTimeFrom: time.Now().AddDate(0, -6, 0),
			IndexTimeEnd:  time.Now(),
			Business:      3,
			Order:         "ctime",
			Sort:          "desc",
			Oid:           []int64{1},
			Distinct:      "oid",
		}
		result, err := d.SearchAuditReportLog(context.TODO(), cond)
		convey.So(err, convey.ShouldBeNil)
		convey.So(result, convey.ShouldNotBeEmpty)
	})
}
