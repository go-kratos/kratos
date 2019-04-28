package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go-common/app/service/bbq/search/api/grpc/v1"
	"go-common/app/service/bbq/search/model"
	"google.golang.org/grpc"
	"os"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:9000", "The server address in the format of host:port")
)

func main() {

	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("11111", err)
		return
	}
	defer conn.Close()

	salesClient := v1.NewSearchClient(conn)
	time1 := time.Now().UnixNano()
	data := new(v1.SaveVideoRequest)
	for i := 0; i < 1; i++ {
		tmp := &v1.VideoESInfo{
			SVID:     int64(4 + i),
			MID:      1,
			CID:      1,
			AVID:     1,
			Title:    "小姐姐 wo 好看 cityy" + string(i),
			Content:  "占位占位",
			Pubtime:  12,
			Duration: 0,
			Original: 123,
			State:    1,
			From:     1,
			VerID:    12,
			Ver:      "lasfjn123",
			Ctime:    123,
			Mtime:    123,
		}
		tags := make(map[int64]string)
		//tags[1] = "舞蹈"
		tags[2] = "直播"
		//tags[3] = "开心"
		for i, v := range tags {
			tmp.Tags = append(tmp.Tags, &v1.VideoESTags{ID: i, Name: v})
		}
		data.List = append(data.List, tmp)
	}
	//res1,err := salesClient.SaveVideo(context.Background(), data)
	//fmt.Println("CreateIndex",res1,err)

	calc := &model.Calc{
		Open:       1,
		PlayRatio:  0.3,
		FavRatio:   0.05,
		LikeRatio:  0.15,
		CoinRatio:  0.1,
		ReplyRatio: 0.2,
		ShareRatio: 0.1,
	}

	hotTags := []string{"美女"}
	where := new(model.Where)
	where.In = make(map[string][]interface{})
	for _, tag := range hotTags {
		where.In["tags.name"] = append(where.In["tags.name"], tag)
	}
	where.In["state"] = append(where.In["state"], 3)
	where.NotIn = make(map[string][]interface{})
	where.NotIn["avid"] = append(where.NotIn["avid"], 27035488)

	where.Lte = make(map[string]int64)
	where.Lte["svid"] = 1168

	where.Gte = make(map[string]int64)
	where.Gte["avid"] = 28457770

	filter := make(map[string]interface{})
	filter["buvid"] = "bbqtestbuvid"
	filter["mid"] = 123
	query := model.Query{
		Calc:   calc,
		Where:  where,
		From:   0,
		Size:   1,
		Filter: filter,
	}
	queryBody, err := json.Marshal(query)
	fmt.Println(query.Where)
	fmt.Println(string(queryBody))

	//del := new(v1.DelVideoBySVIDRequest)
	//del.SVIDs = append(del.SVIDs, 84)
	//res4, err := salesClient.DelVideoBySVID(context.Background(), del)
	//fmt.Println(res4, err)

	res3, err := salesClient.RecVideoData(context.Background(), &v1.RecVideoDataRequest{Query: string(queryBody)})
	fmt.Println(res3)
	fmt.Println(err)
	//return

	//res2,err := salesClient.RecVideoData(context.Background(), &v1.RecVideoDataRequest{PageNum:0,PageSize:2})
	//res2,err := salesClient.RecVideoData(context.Background(), &v1.RecVideoDataRequest{Query:"{\"calc\":{\"open\":1,\"fav_ratio\":1.0,\"like_ratio\":0.5,\"pub_ratio\":0.3},\"where\":{\"in\":{\"title\":[\"舞蹈\",\"美女\"],\"tag.Name\":[\"小姐姐\"]}},\"limit\":2}"})
	res2, err := salesClient.RecVideoData(context.Background(), &v1.RecVideoDataRequest{Query: string(queryBody)})
	//res2,err := salesClient.RecVideoData(context.Background(), &v1.RecVideoDataRequest{Query:"{}"})
	//fmt.Println("VideoData",res2,err)

	var out bytes.Buffer
	var b []byte
	b, _ = json.Marshal(res2)
	json.Indent(&out, b, "", "\t")
	out.WriteTo(os.Stdout)

	time2 := time.Now().UnixNano()

	fmt.Println((time2 - time1) / 1e6)
	//fmt.Println(queryBody)
	//salesClientMis := v1.NewSalesMisClient(conn)
	//
	//
	//res2,err := salesClientMis.GetGroupOrdersMis(context.Background(), &v1.GetGroupOrdersMisRequest{OrderID:0})
	//fmt.Println("222",res2,err)
	fmt.Println(err)
	return
}
