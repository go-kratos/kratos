package regexp

import (
	"regexp"
	"runtime"
	"testing"

	pcre "github.com/GRbit/go-pcre"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 4)
}

type Check struct {
	Expect  bool
	Content string
	RegStr  string
}

var (
	routineSize = runtime.NumCPU() * 100
	checks      = []Check{
		Check{true, "朝鮮金三胖竊取虛擬貨幣", `(朝鮮).{0,5}(掃黃|性服務|黃片|緊急斷電|査光盤|機關槍槍斃|色情影片|造假鈔|僞印|假幣|假鈔|印美元|僞造美元|造假美元|超級美元|假人民幣|僞造人民幣|人民幣假鈔|印刷人民幣|僞造|造假幣|製毒|販毒|製造毒品|走私毒品|罌粟|洗錢|竊取虛擬貨幣|阬大哥|皇長子|女人地位|底層女性|底層婦女)`},
		Check{false, "朝鮮人民主体思想", `(朝鮮).{0,5}(掃黃|性服務|黃片|緊急斷電|査光盤|機關槍槍斃|色情影片|造假鈔|僞印|假幣|假鈔|印美元|僞造美元|造假美元|超級美元|假人民幣|僞造人民幣|人民幣假鈔|印刷人民幣|僞造|造假幣|製毒|販毒|製造毒品|走私毒品|罌粟|洗錢|竊取虛擬貨幣|阬大哥|皇長子|女人地位|底層女性|底層婦女)`},
		Check{true, "求購你什么也不知道的2,5-二甲氧基苯乙胺", `(銷售|直銷|售|買賣|提供|齣售|訂購|賣|專賣|供應|批髮|貨到付欵|直供|提鍊|購買|求購).{0,10}(麥角二乙胺|2,5-二甲氧基苯乙胺|三氯甲烷|地索嗎啡|硫代芬太尼|右嗎拉胺|1－苯基－2－丙酮|乙醚|副甲氧基安非他明|2,5-二甲氧基-4-溴苯乙胺|醋氫可待囙|阿灋甲基芬太尼|二乙噻丁|貝齊米特|黃樟素|麥角新堿|噴他佐辛|苯甲麯秦|呱苯甲醕|四氫西泮|丙己君|紮來普隆|罌粟殼|氫嗎啡醕|舒芬太尼|對氟芬太尼|去甲嗎啡|可多尅辛|戊巴比妥|甲米雷司|芬乙茶堿|瑞芬太尼)`},
		Check{true, "萝莉控", `萝莉`},
		Check{true, "pan.baidu.com/s/1aFsD8wgR3ozVrO_ysdRE8A", `(?:http|https|www|pan|w w w|baidu\.com)(?:[\s\.:\/\/]{1,})([\w%+:\s\/\.?=]{1,}[^:：[" ]{1,}[(๑•. •๑)´_ゝ｀﹏ŏ*≧▽≦ツ┏━┓Σ°△|︴；"▔□:з」∠\.\w\/]*\w{1,})`},
		Check{false, "不命中气死偶累", `(?:http|https|www|pan|w w w|baidu\.com)(?:[\s\.:\/\/]{1,})([\w%+:\s\/\.?=]{1,}[^:：[" ]{1,}[(๑•. •๑)´_ゝ｀﹏ŏ*≧▽≦ツ┏━┓Σ°△|︴；"▔□:з」∠\.\w\/]*\w{1,})`},
		Check{true, "狗狗2", `saya|shaya|hha|hh1|狗狗2`},
	}
)

func BenchmarkStdParallel(b *testing.B) {
	regs := make([]*regexp.Regexp, len(checks))
	for i, c := range checks {
		regs[i] = regexp.MustCompile(c.RegStr)
	}
	b.SetParallelism(routineSize)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < len(checks); i++ {
				regs[i].MatchString(checks[i].Content)
			}
		}
	})
}

func BenchmarkPCREParallel(b *testing.B) {
	regs := make([]pcre.Regexp, len(checks))
	for i, c := range checks {
		regs[i] = pcre.MustCompile(c.RegStr, pcre.UTF8)
	}
	b.SetParallelism(routineSize)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < len(checks); i++ {
			}
		}
	})
}

func BenchmarkPCREWithJITParallel(b *testing.B) {
	regs := make([]pcre.Regexp, len(checks))
	for i, c := range checks {
		regs[i] = pcre.MustCompileJIT(c.RegStr, pcre.UTF8, 0)
	}
	b.SetParallelism(routineSize)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < len(checks); i++ {
			}
		}
	})
}

func TestStd(t *testing.T) {
	regs := make([]*regexp.Regexp, len(checks))
	for i, c := range checks {
		regs[i] = regexp.MustCompile(c.RegStr)
	}
	for i := 0; i < len(checks); i++ {
		if regs[i].MatchString(checks[i].Content) != checks[i].Expect {
			t.Fatal()
		}
	}
}

func TestPCRE(t *testing.T) {
	regs := make([]pcre.Regexp, len(checks))
	for i, c := range checks {
		regs[i] = pcre.MustCompile(c.RegStr, pcre.UTF8)
	}
	for i := 0; i < len(checks); i++ {
		if regs[i].MatcherString(checks[i].Content, 0).Matches != checks[i].Expect {
			t.Fatal()
		}
	}
}

func TestPCREWithJIT(t *testing.T) {
	regs := make([]pcre.Regexp, len(checks))
	for i, c := range checks {
		regs[i] = pcre.MustCompileJIT(c.RegStr, pcre.UTF8, pcre.STUDY_JIT_PARTIAL_SOFT_COMPILE)
	}
	for i := 0; i < len(checks); i++ {
		if regs[i].MatcherString(checks[i].Content, 0).Matches != checks[i].Expect {
			t.Fatal()
		}
	}
}

func TestRegexp(t *testing.T) {
	Convey("regexp", t, func() {
		var (
			content = "狗狗2"
		)
		reg, err := Compile(`saya|shaya|hha|hh1|狗狗2`)
		if err != nil {
			t.Fatal(err)
		}
		So(reg.FindAllIndex([]byte(content)), ShouldNotBeEmpty)
		So(reg.MatchString(content), ShouldBeTrue)
	})
}

func TestPCREJitFindAllIndex(t *testing.T) {
	c := "(@|[43]G|[保佑祝]我|[修成登]仙|[养出][了肥]|[处處][女男])"
	re, err := pcre.CompileJIT(c, pcre.UTF8, pcre.STUDY_JIT_COMPILE)
	if err != nil {
		t.Errorf("3444 %+v", err)
	}
	i := re.FindAllIndex([]byte("@我人"), pcre.NOTBOL)
	m := re.Matcher([]byte("@我人"), 0)
	t.Error("FindIndex start", i, re.Groups(), []byte("@我人"), m.Group(1))
	b := re.MatchString("@我人", 0)
	t.Error(b)

}

func TestPCREFindAllIndex(t *testing.T) {
	c := "(@|[43]G|[保佑祝]我|[修成登]仙|[养出][了肥]|[处處][女男])"
	re, err := pcre.Compile(c, pcre.UTF8)
	if err != nil {
		t.Errorf("3444 %+v", err)
	}
	i := re.FindAllIndex([]byte("我人"), 0)
	m := re.Matcher([]byte("我人"), 0)
	t.Error("FindIndex start", i, re.Groups(), []byte("我人"), m.Group(1))
	b := re.MatchString("鬼喜欢用boo这词来吓人", 0)
	t.Error(b)

}
