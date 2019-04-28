package dao

import (
	"context"
	"testing"

	artmdl "go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Rank(t *testing.T) {
	data := `{"code":0,"source_date":"2018-01-02","list":[{"aid":1,"mid":137952,"score":565918,"view":176708,"reply":2108,"favorites":1517,"coin":6816,"likes":10454},{"aid":2,"mid":144900177,"score":300536,"view":652823,"reply":2661,"favorites":10268,"coin":470,"likes":4130},{"aid":3,"mid":32708462,"score":241845,"view":485737,"reply":969,"favorites":7347,"coin":1290,"likes":5542},{"aid":4,"mid":124799,"score":188953,"view":46594,"reply":595,"favorites":797,"coin":1771,"likes":6268}],"num":4,"note":"统计7日内新投稿的数据综合得分"}`
	Convey("get data", t, func() {
		httpMock("GET", d.c.Article.RankHost+"/data/rank/article/all-7.json").Reply(200).JSON(data)
		ranks, err := d.Rank(context.TODO(), artmdl.RankWeek, "")
		So(err, ShouldBeNil)
		So(ranks, ShouldNotBeEmpty)
	})
}
