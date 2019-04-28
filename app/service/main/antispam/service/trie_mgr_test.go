package service

import (
	"math/rand"
	"testing"
	"time"

	"go-common/app/service/main/antispam/util"
	"go-common/app/service/main/antispam/util/trie"
)

func TestConcurrentTrieFind(t *testing.T) {
	conTrie := NewConcurrentTrie()
	conTrie.trier = trie.New()
	conTrie.Put("jjjj", &KeywordLimitInfo{KeywordID: 111})
	k, v := conTrie.find("jjjj")
	t.Logf("k: %v, v: %v", k, v)
	if v == nil {
		t.FailNow()
	}
}

func BenchmarkConcurrentTrieFind(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	tr := NewConcurrentTrie()
	tr.trier = trie.New()

	for i := 0; i < b.N; i++ {
		tr.Put(util.RandStr(20), &KeywordLimitInfo{})
	}

	tr.Put("地方考虑saDFFDSALK", 8888)
	tr.Put("都说了开发贷款", 888)

	for i := 0; i < b.N; i++ {
		go func() {
			tr.Put(util.RandStr(20), &KeywordLimitInfo{})
		}()
		tr.find("dfa啥都发生地方的施工费按发的噶是打发士大夫撒旦噶尔尕热狗怕的是结果来看；砥节奉公；来人速度感而过;sfdsfas fsadf asd fsad f asd都说了sdfs gfdgd jimmy开发速度来发噶都说了开发贷款时间范德萨了空间发的是 jimmy按时到路口发生撒地方考虑saDFFDSALKFDFASDFASDFSDFSADFRGEWTRETGERG")
	}
}
