package trie

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"go-common/app/service/main/antispam/util"
)

func TestRuneTrieAdd(t *testing.T) {
	tr := NewRuneTrie()
	tr.Put("jimmy", 1)
	tr.Put("anny", 87)
	tr.Put("xxxx", 23)
	tr.Put("jim", 2)

	if v := tr.Get("jimxx"); v != nil {
		t.Errorf("expected nil, got %v", v)
	}
	if v := tr.Get("jimmy"); v.(int) != 1 {
		t.Errorf("expected val, got %v", v)
	}
	if v := tr.Get("anny"); v.(int) != 87 {
		t.Errorf("expected val, got %v", v)
	}
	if v := tr.Get("xxxx"); v.(int) != 23 {
		t.Errorf("expected val, got %v", v)
	}
	if v := tr.Get("jim"); v.(int) != 2 {
		t.Errorf("expected val, got %v", v)
	}
}

func BenchmarkRuneTriePut(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	tr := NewRuneTrie()
	for i := 0; i < b.N; i++ {
		tr.Put(util.RandStr(10), 845)
	}
}

func TestRuneTrieFind(t *testing.T) {
	tr := NewRuneTrie()
	tr.Put("我才是大佬", 2)
	tr.Put("我才是大佬", 88)
	tr.Put("mm", 1)
	tr.Put("mmp", 2)
	tr.Put("my name is jimmymmp", 2)
	tr.Put("xxx", 88)
	tr.Put("jimmy xxx, hhjhmmp", 2)

	cases := []struct {
		content     string
		expectKey   string
		expectValue int
	}{
		{
			content:     "mm",
			expectKey:   "mm",
			expectValue: 1,
		},
		{
			content:     "m都xx发生地方范德萨发爱迪生刚发的否多少发生的否阿萨德否收到符文大师否xxxmy name is jimy, hhjhmp",
			expectKey:   "xxx",
			expectValue: 88,
		},
		{
			content:     "m都mxx发生地方范德萨发爱迪生刚发的否多少发生的否阿萨德否收到符文大师否xxxmy name is jimy, hhjhmp",
			expectKey:   "xxx",
			expectValue: 88,
		},
		{
			content:     "我才是大佬",
			expectKey:   "我才是大佬",
			expectValue: 88,
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			k, v := tr.Find(c.content, "")
			if v == nil {
				t.Fatal(errors.New("val is nil"))
			}

			if k != c.expectKey || v.(int) != c.expectValue {
				t.Errorf("want key: %s, val:%v, got key:%s, val:%v", c.expectKey, c.expectValue, k, v)
			}
		})
	}
}

//BenchmarkTrieListFind-4		   50000     25347 ns/op
func BenchmarkRuneTrieFind(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	tr := NewRuneTrie()
	for i := 0; i < b.N; i++ {
		tr.Put(util.RandStr(10), i)
	}
	tr.Put("地方考虑saDFFDSALK", 8888)
	tr.Put("都说了开发贷款", 7512)

	for i := 0; i < b.N; i++ {
		tr.Find("dfa啥都发生地方的施工费按发的噶是打发士大夫撒旦噶尔尕热狗怕的是结果来看；砥节奉公；来人速度感而过;sfdsfas fsadf asd fsad f asd都说了sdfs gfdgd jimmy开发速度来发噶都说了开发贷款时间范德萨了空间发的是 jimmy按时到路口发生撒地方考虑saDFFDSALKFDFASDFASDFSDFSADFRGEWTRETGERG", "")
	}
}
