package chinese

import (
	"context"
	"testing"
)

func TestConvert(t *testing.T) {
	Init()
	in := `请不要怀疑,这是一个由人工智能推荐的频道。`
	out := Convert(context.Background(), in)

	t.Logf("in:%s,out:%s", in, out)
	in = `说起来你可能不信,我是考试考进来的`
	out = Convert(context.Background(), in)
	t.Logf("in:%s,out:%s", in, out)
}

func BenchmarkConvert(b *testing.B) {
	var testcase = []string{
		"说起来你可能不信,我是考试考进来的",
		"说起来你可能不信,我是花钱找关系进来的",
		"请不要怀疑,这是一个由人工智能推荐的频道",
		"我开挖掘机拆屋的时候听特别带感",
		"1990年真实记录，当时的秋名山的日常",
		"1990年藤原豆腐店成了连锁店 没错这些车都是送豆腐的",
		"Go语言,从底层到应用，视Golang的环境搭建、基础知识、进阶知识、项目实践、Redis基础及其项目实践（海量用户通讯系统）、算法与数据结构基础知识的golang实现。",
	}
	Init()
	for i := 0; i < b.N; i++ {
		out := Convert(context.Background(), testcase[i%len(testcase)])

		b.Logf("in:%s,out:%s", testcase[i%len(testcase)], out)
	}
}
func BenchmarkConverts(b *testing.B) {
	var testcase = []string{
		"说起来你可能不信,我是考试考进来的",
		"说起来你可能不信,我是花钱找关系进来的",
		"请不要怀疑,这是一个由人工智能推荐的频道",
		"我开挖掘机拆屋的时候听特别带感",
		"1990年真实记录，当时的秋名山的日常",
		"1990年藤原豆腐店成了连锁店 没错这些车都是送豆腐的",
		"Go语言,从底层到应用，视Golang的环境搭建、基础知识、进阶知识、项目实践、Redis基础及其项目实践（海量用户通讯系统）、算法与数据结构基础知识的golang实现。",
	}
	Init()
	var out map[string]string
	for i := 0; i < b.N; i++ {
		out = Converts(context.Background(), testcase...)
	}
	b.Logf("in:%s,out:%s", testcase, out)
}
