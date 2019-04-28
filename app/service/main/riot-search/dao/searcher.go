package dao

import (
	"go-common/app/service/main/riot-search/model"

	"github.com/go-ego/riot/types"
)

// SearchIDOnly only return aids
func (d *Dao) SearchIDOnly(arg *model.RiotSearchReq) *model.IDsResp {
	if arg.Keyword == "" {
		return nil
	}
	var docIDs map[uint64]bool
	if len(arg.IDs) != 0 {
		docIDs = make(map[uint64]bool, len(arg.IDs))
		for _, id := range arg.IDs {
			docIDs[id] = true
		}
	}
	output := d.searcher.Search(types.SearchReq{
		Text:    arg.Keyword,
		DocIds:  docIDs,
		Timeout: d.c.Riot.Timeout,
		RankOpts: &types.RankOpts{
			// 从第几条结果开始输出
			OutputOffset: (arg.Pn - 1) * arg.Ps,
			// 最大输出的搜索结果数，为 0 时无限制
			MaxOutputs: arg.Ps,
		},
	})
	docLength := len(output.Docs.(types.ScoredDocs))
	tokenLength := len(output.Tokens)
	res := &model.IDsResp{
		IDs:    make([]uint64, docLength),
		Tokens: make([]string, tokenLength),
		Page: &model.Page{
			PageNum:  arg.Pn,
			PageSize: arg.Ps,
			Total:    docLength,
		},
	}
	for i, doc := range output.Docs.(types.ScoredDocs) {
		res.IDs[i] = doc.DocId
	}
	copy(res.Tokens, output.Tokens)
	res.IDs, res.Page.Total = uniqueIDs(res.IDs)
	return res
}

func uniqueIDs(IDs []uint64) (uIDs []uint64, length int) {
	m := make(map[uint64]struct{})
	for _, ID := range IDs {
		if _, ok := m[ID]; !ok {
			m[ID] = struct{}{}
			uIDs = append(uIDs, ID)
			length++
		}
	}
	return
}

// Search return archives info
func (d *Dao) Search(arg *model.RiotSearchReq) *model.DocumentsResp {
	if arg.Keyword == "" {
		return nil
	}
	var docIDs map[uint64]bool
	if len(arg.IDs) != 0 {
		docIDs = make(map[uint64]bool, len(arg.IDs))
		for _, id := range arg.IDs {
			docIDs[id] = true
		}
	}
	output := d.searcher.Search(types.SearchReq{
		Text:    arg.Keyword,
		DocIds:  docIDs,
		Timeout: d.c.Riot.Timeout,
		RankOpts: &types.RankOpts{
			// 从第几条结果开始输出
			OutputOffset: (arg.Pn - 1) * arg.Ps,
			// 最大输出的搜索结果数，为 0 时无限制
			MaxOutputs: arg.Ps,
		},
	})
	docLength := len(output.Docs.(types.ScoredDocs))
	tokenLength := len(output.Tokens)
	res := &model.DocumentsResp{
		Documents: make([]model.Document, docLength),
		Tokens:    make([]string, tokenLength),
		Page: &model.Page{
			PageNum:  arg.Pn,
			PageSize: arg.Ps,
			Total:    docLength,
		},
	}
	for i, doc := range output.Docs.(types.ScoredDocs) {
		res.Documents[i].ID = doc.DocId
		res.Documents[i].Content = doc.Content
	}
	copy(res.Tokens, output.Tokens)
	return res
}

// Has return DocId exists
func (d *Dao) Has(id uint64) bool {
	return d.searcher.HasDoc(id)
}
