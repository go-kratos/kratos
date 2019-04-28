package retrieve

import (
	"context"

	"github.com/json-iterator/go"

	rpc "go-common/app/service/bbq/recsys/api/grpc/v1"
	"go-common/app/service/bbq/recsys/model"
	searchv1 "go-common/app/service/bbq/search/api/grpc/v1"
	"go-common/library/log"
	"strconv"
	"strings"
)

//召回策略
const (
	Hot       = "hot"
	Selection = "selection"
	Relevant  = "relevant"
	Tag       = "tag"
)

// RetrieverManager manages multiple retrieve functions
type RetrieverManager struct {
	Retrievers     []Retriever
	PreRetrievers  []Retriever
	PostRetrievers []Retriever
	RetrieveFunc   RetrieverFunc
}

//NewRetrieverManager ...
func NewRetrieverManager() (m *RetrieverManager) {
	m = &RetrieverManager{
		Retrievers:     make([]Retriever, 0),
		PreRetrievers:  make([]Retriever, 0),
		PostRetrievers: make([]Retriever, 0),
		RetrieveFunc:   DefaultRetrieveFunc,
	}
	hot := &hotRetriever{}
	tag := &tagRetriever{}
	operation := &operationRetriever{}
	relevant := &relevantRetriever{}
	m.Retrievers = append(m.Retrievers, hot, tag, operation, relevant)

	m.PostRetrievers = append(m.PostRetrievers, hot)

	return
}

//Retrievers ...
type Retrievers struct {
	MethodName string
	Retrieve   RetrieverFunc
}

//RetrieverFunc ...
type RetrieverFunc func(c context.Context, r Retriever, searchClient searchv1.SearchClient, request rpc.RecsysRequest, userProfile *model.UserProfile, response chan rpc.RecsysResponse)

//ColdStartRetriever ...
type ColdStartRetriever struct {
	Retrievers
}

//Merge ...
func (m *RetrieverManager) Merge(response *rpc.RecsysResponse) {

	records := make([]*rpc.RecsysRecord, 0)
	set := map[int64]int{}
	for _, record := range response.List {
		if count, ok := set[record.Svid]; ok {
			set[record.Svid] = count + 1
		} else {
			set[record.Svid] = 1
			records = append(records, record)
		}
	}
	response.List = records

	records = make([]*rpc.RecsysRecord, 0)
	set = map[int64]int{}
	for _, record := range response.List {
		if avid, ok := record.Map[model.AVID]; ok {
			avidInt, _ := strconv.ParseInt(avid, 10, 64)
			if count, ok := set[avidInt]; ok {
				set[avidInt] = count + 1
			} else {
				set[avidInt] = 1
				records = append(records, record)
			}
		}
	}
	response.List = records

}

//Retriever .
type Retriever interface {
	name() (name string)

	queryRewrite(c context.Context, request rpc.RecsysRequest, userProfile *model.UserProfile) (req *searchv1.RecVideoDataRequest, err error)
}

//Query .
type Query struct {
	Calc   *Calc                  `json:"calc"`
	Where  *Where                 `json:"where"`
	Filter map[string]interface{} `json:"filter"`
	From   int                    `json:"from"`
	Size   int                    `json:"size"`
}

//Calc .
type Calc struct {
	Open       int64   `json:"open"`
	PlayRatio  float64 `json:"play_ratio"`
	FavRatio   float64 `json:"fav_ratio"`
	LikeRatio  float64 `json:"like_ratio"`
	PubRatio   float64 `json:"pub_ratio"`
	CoinRatio  float64 `json:"coin_ratio"`
	ReplyRatio float64 `json:"reply_ratio"`
	ShareRatio float64 `json:"share_ratio"`
}

//Where .
type Where struct {
	In    map[string][]interface{} `json:"in"`
	NotIn map[string][]interface{} `json:"not_in"`
}

//DefaultRetrieveFunc is default retrieve function
func DefaultRetrieveFunc(c context.Context, r Retriever, searchClient searchv1.SearchClient, request rpc.RecsysRequest, userProfile *model.UserProfile, response chan rpc.RecsysResponse) {

	result := rpc.RecsysResponse{}

	req, err := r.queryRewrite(c, request, userProfile)
	if err != nil {
		log.Error("query rewrite error: ", err)
		response <- result
		return
	}
	if req == nil {
		response <- result
		return
	}

	res, err := searchClient.RecVideoData(c, req)
	if err != nil {
		log.Error("Retrieve error: ", err)
		response <- result
		return
	}

	for _, videoEsInfo := range res.List {
		record := &rpc.RecsysRecord{Map: make(map[string]string)}
		fillRecord(record, videoEsInfo)
		record.Map[model.Retriever] = r.name()
		result.List = append(result.List, record)
	}

	if r.name() == Relevant {
		likeTags := make(map[string]float64)
		posTags := make(map[string]float64)
		negTags := make(map[string]float64)

		for _, record := range result.List {
			svid := record.Svid
			if _, ok := userProfile.LikeVideos[svid]; ok {
				if itemTags, ok := record.Map[model.TagsName]; ok {
					for _, tag := range strings.Split(itemTags, "|") {
						if tagScore, ok := likeTags[tag]; ok {
							likeTags[tag] = tagScore + 1.0
						} else {
							likeTags[tag] = 1.0
						}
					}
				}
			}
			if _, ok := userProfile.PosVideos[svid]; ok {
				if itemTags, ok := record.Map[model.TagsName]; ok {
					for _, tag := range strings.Split(itemTags, "|") {
						if tagScore, ok := posTags[tag]; ok {
							posTags[tag] = tagScore + 1.0
						} else {
							posTags[tag] = 1.0
						}
					}
				}
			}

			if _, ok := userProfile.NegVideos[svid]; ok {
				if itemTags, ok := record.Map[model.TagsName]; ok {
					for _, tag := range strings.Split(itemTags, "|") {
						if tagScore, ok := negTags[tag]; ok {
							negTags[tag] = tagScore + 1.0
						} else {
							negTags[tag] = 1.0
						}
					}
				}
			}
		}

		userProfile.LikeTags = likeTags
		userProfile.PosTags = posTags
		userProfile.NegTags = negTags

		calc := &Calc{
			Open:       1,
			PlayRatio:  0.3,
			FavRatio:   0.05,
			LikeRatio:  0.15,
			PubRatio:   0.1,
			ShareRatio: 0.1,
			CoinRatio:  0.1,
			ReplyRatio: 0.2,
		}

		where := new(Where)
		where.In = make(map[string][]interface{})

		where.NotIn = make(map[string][]interface{})
		for _, id := range userProfile.DedupVideos {
			where.NotIn[model.CID] = append(where.NotIn[model.CID], id)
		}
		for _, id := range userProfile.PosVideos {
			where.NotIn[model.SVID] = append(where.NotIn[model.SVID], id)
		}
		for _, id := range userProfile.NegVideos {
			where.NotIn[model.SVID] = append(where.NotIn[model.SVID], id)
		}

		hasTag := false
		for tag, score := range posTags {
			if score > 0 {
				where.In[model.TagsName] = append(where.In[model.TagsName], tag)
				hasTag = true
			}
		}

		if !hasTag {
			log.Error("Relevant videos has no tag: ")
			response <- rpc.RecsysResponse{}
			return
		}

		filter := make(map[string]interface{})
		filter["buvid"] = request.BUVID
		filter["mid"] = request.MID

		query := Query{
			Calc:   calc,
			Where:  where,
			Filter: filter,
			From:   0,
			Size:   100,
		}

		queryBody, _ := jsoniter.Marshal(query)
		log.Info(r.name(), string(queryBody))
		req = &searchv1.RecVideoDataRequest{Query: string(queryBody)}

		res, err := searchClient.RecVideoData(c, req)
		if err != nil {
			log.Error("Retrieve error: ", err)
			response <- rpc.RecsysResponse{}
			return
		}

		result = rpc.RecsysResponse{}
		for _, videoEsInfo := range res.List {
			record := &rpc.RecsysRecord{Map: make(map[string]string)}
			fillRecord(record, videoEsInfo)
			record.Map[model.Retriever] = r.name()
			result.List = append(result.List, record)
		}

	}
	response <- result

}

type hotRetriever struct {
	Retriever
}

func (r *hotRetriever) name() (name string) {
	name = Hot
	return
}

func (r *hotRetriever) queryRewrite(c context.Context, request rpc.RecsysRequest, userProfile *model.UserProfile) (req *searchv1.RecVideoDataRequest, err error) {
	calc := &Calc{
		Open:       1,
		PlayRatio:  0.3,
		FavRatio:   0.05,
		LikeRatio:  0.15,
		PubRatio:   0.1,
		ShareRatio: 0.1,
		CoinRatio:  0.1,
		ReplyRatio: 0.2,
	}

	//hotTags := []string{"美女", "性感", "女神", "英雄联盟", "电子竞技", "小姐姐", "LOL"}
	//hotTags := []string{"舞蹈", "宅舞", "mmd", "英雄联盟", "电子竞技", "小姐姐", "LOL"}
	//hotTags := []string{"mmd"}

	where := new(Where)
	where.In = make(map[string][]interface{})
	where.In[model.State] = append(where.In[model.State], model.State1, model.State0, model.State3, model.State4, model.State5)

	where.NotIn = make(map[string][]interface{})
	for _, id := range userProfile.DedupVideos {
		where.NotIn[model.CID] = append(where.NotIn[model.CID], id)
	}

	filter := make(map[string]interface{})
	filter["buvid"] = request.BUVID
	filter["mid"] = request.MID

	query := Query{
		Calc:   calc,
		Where:  where,
		Filter: filter,
		From:   0,
		Size:   100,
	}

	queryBody, err := jsoniter.Marshal(query)
	log.Info(r.name(), string(queryBody))
	req = &searchv1.RecVideoDataRequest{Query: string(queryBody)}
	return
}

type tagRetriever struct {
	Retriever
}

func (r *tagRetriever) name() (name string) {
	name = Tag
	return
}

func (r *tagRetriever) queryRewrite(c context.Context, request rpc.RecsysRequest, userProfile *model.UserProfile) (req *searchv1.RecVideoDataRequest, err error) {

	if len(userProfile.BiliTags) <= 0 && len(userProfile.Zones1) == 0 && len(userProfile.Zones2) == 0 {
		err = nil
		return
	}

	calc := &Calc{
		Open:       1,
		PlayRatio:  0.3,
		FavRatio:   0.05,
		LikeRatio:  0.15,
		PubRatio:   0.1,
		ShareRatio: 0.1,
		CoinRatio:  0.1,
		ReplyRatio: 0.2,
	}

	where := new(Where)
	where.In = make(map[string][]interface{})
	where.In[model.State] = append(where.In[model.State], model.State1, model.State0, model.State3, model.State4, model.State5)

	for tag := range userProfile.BiliTags {
		tagID, _ := strconv.ParseInt(tag, 10, 64)
		where.In[model.TagsID] = append(where.In[model.TagsID], tagID)
	}
	for tag := range userProfile.Zones1 {
		tagID, _ := strconv.ParseInt(tag, 10, 64)
		where.In[model.TagsID] = append(where.In[model.TagsID], tagID)
	}
	for tag := range userProfile.Zones2 {
		tagID, _ := strconv.ParseInt(tag, 10, 64)
		where.In[model.TagsID] = append(where.In[model.TagsID], tagID)
	}

	where.NotIn = make(map[string][]interface{})
	for _, id := range userProfile.DedupVideos {
		where.NotIn[model.CID] = append(where.NotIn[model.CID], id)
	}

	filter := make(map[string]interface{})
	filter["buvid"] = request.BUVID
	filter["mid"] = request.MID

	query := Query{
		Calc:   calc,
		Where:  where,
		Filter: filter,
		From:   0,
		Size:   100,
	}

	queryBody, err := jsoniter.Marshal(query)
	log.Info(r.name(), string(queryBody))
	req = &searchv1.RecVideoDataRequest{Query: string(queryBody)}

	return
}

type operationRetriever struct {
	Retriever
}

func (r *operationRetriever) name() (name string) {
	name = Selection
	return
}

func (r *operationRetriever) queryRewrite(c context.Context, request rpc.RecsysRequest, userProfile *model.UserProfile) (req *searchv1.RecVideoDataRequest, err error) {

	calc := &Calc{
		Open:       1,
		PlayRatio:  0.3,
		FavRatio:   0.05,
		LikeRatio:  0.15,
		PubRatio:   0.1,
		ShareRatio: 0.1,
		CoinRatio:  0.1,
		ReplyRatio: 0.2,
	}

	where := new(Where)
	where.In = make(map[string][]interface{})
	where.In[model.State] = append(where.In[model.State], model.State5)

	where.NotIn = make(map[string][]interface{})
	for _, id := range userProfile.DedupVideos {
		where.NotIn[model.CID] = append(where.NotIn[model.CID], id)
	}

	filter := make(map[string]interface{})
	filter["buvid"] = request.BUVID
	filter["mid"] = request.MID

	query := Query{
		Calc:   calc,
		Where:  where,
		Filter: filter,
		From:   0,
		Size:   100,
	}

	queryBody, err := jsoniter.Marshal(query)
	log.Info(r.name(), string(queryBody))
	req = &searchv1.RecVideoDataRequest{Query: string(queryBody)}

	return
}

type relevantRetriever struct {
	Retriever
}

func (r *relevantRetriever) name() (name string) {
	name = Relevant
	return
}

func (r *relevantRetriever) queryRewrite(c context.Context, request rpc.RecsysRequest, userProfile *model.UserProfile) (req *searchv1.RecVideoDataRequest, err error) {

	//TODO
	//userProfile.SessionPosVideos
	if len(userProfile.PosVideos) == 0 && len(userProfile.LikeVideos) == 0 {
		return
	}

	calc := &Calc{
		Open:       1,
		PlayRatio:  0.3,
		FavRatio:   0.05,
		LikeRatio:  0.15,
		PubRatio:   0.1,
		ShareRatio: 0.1,
		CoinRatio:  0.1,
		ReplyRatio: 0.2,
	}

	where := new(Where)
	where.In = make(map[string][]interface{})
	where.NotIn = make(map[string][]interface{})
	where.In[model.State] = append(where.In[model.State], model.State1, model.State0, model.State3, model.State4, model.State5)

	for id := range userProfile.PosVideos {
		where.In[model.SVID] = append(where.In[model.SVID], id)
	}
	for id := range userProfile.NegVideos {
		where.In[model.SVID] = append(where.In[model.SVID], id)
	}
	for id := range userProfile.LikeVideos {
		where.In[model.SVID] = append(where.In[model.SVID], id)
	}

	query := Query{
		Calc:  calc,
		Where: where,
		From:  0,
		Size:  100,
	}

	queryBody, err := jsoniter.Marshal(query)
	log.Info(r.name(), string(queryBody))
	req = &searchv1.RecVideoDataRequest{Query: string(queryBody)}

	return
}

func fillRecord(record *rpc.RecsysRecord, videoEsInfo *searchv1.RecVideoInfo) {

	record.Svid = videoEsInfo.SVID

	record.Map[model.Title] = strings.Replace(videoEsInfo.Title, "\"", " ", -1)
	//方便解析成json格式
	record.Map[model.Content] = strings.Replace(videoEsInfo.Content, "\"", " ", -1)
	record.Map[model.AVID] = strconv.FormatInt(videoEsInfo.AVID, 10)
	record.Map[model.CID] = strconv.FormatInt(videoEsInfo.CID, 10)
	record.Map[model.SVID] = strconv.FormatInt(videoEsInfo.SVID, 10)

	record.Map[model.TID] = strconv.FormatInt(videoEsInfo.Tid, 10)
	record.Map[model.SubTid] = strconv.FormatInt(videoEsInfo.SubTid, 10)
	record.Map[model.PlayHive] = strconv.FormatInt(videoEsInfo.PlayHive, 10)
	record.Map[model.FavHive] = strconv.FormatInt(videoEsInfo.FavHive, 10)
	record.Map[model.LikesHive] = strconv.FormatInt(videoEsInfo.LikesHive, 10)
	record.Map[model.CoinHive] = strconv.FormatInt(videoEsInfo.CoinHive, 10)

	record.Map[model.State] = strconv.FormatInt(videoEsInfo.State, 10)
	record.Map[model.UperMid] = strconv.FormatInt(videoEsInfo.MID, 10)

	tagNames := make([]string, 0)
	tagTypes := make([]string, 0)
	tagIDs := make([]string, 0)
	for _, tag := range videoEsInfo.Tags {
		tagNames = append(tagNames, tag.Name)
		tagTypes = append(tagTypes, strconv.Itoa(int(tag.Type)))
		tagIDs = append(tagIDs, strconv.Itoa(int(tag.ID)))
	}
	record.Map[model.TagsName] = strings.Join(tagNames, "|")
	record.Map[model.TagsType] = strings.Join(tagTypes, "|")
	record.Map[model.TagsID] = strings.Join(tagIDs, "|")

}
