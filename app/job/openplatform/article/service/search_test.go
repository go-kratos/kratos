package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Search(t *testing.T) {
	Convey("work", t, func() {
		a := `
<div><em><br>AlphaGO战胜柯洁、Google翻译准确度大幅提升、自动驾驶技术日趋成熟；这些无一不在告诉我们人工智能技术正飞速发展。那么如果让人工智能来玩STG会怎么样呢？<br></em><br></div><div><a href="http://www.bilibili.com/video/av12432">视频链接</a><br>STG也称为弹幕游戏。往往由于满屏幕都是子弹，而使一部分玩家望而却步。长久以来，初学者觉得游戏难度无法想象，简直反人类，而高级玩家却觉得还不够难，不够过瘾。这对矛盾已经成为了影响STG发展的主要矛盾。<br>av12679</div><div>av23412<br>现今，大多数STG都会为了解决或缓解这一矛盾而采取了一些措施，比如大多数游戏会设置多个难度，使初学者玩简单难度，而高手可以玩疯狂难度；有的游戏会设置教学关卡，指导初学者避弹；有一些游戏设置了自动雷击功能来为初学者降低难度。<br><br></div><div><strong><br>然而就在最近，我们在游戏《弹幕音乐绘》中，提出了一种更加别出心裁的解决方案——自动避弹。<br></strong><em><br>AlphaGO战胜柯洁、Google翻译准确度大幅提升、自动驾驶技术日趋成熟；这些无一不在告诉我们人工智能技术正飞速发展。那么如果让人工智能来玩STG会怎么样呢？<br></em><br></div><div><br>STG也称为弹幕游戏。往往由于满屏幕都是子弹，而使一部分玩家望而却步。长久以来，初学者觉得游戏难度无法想象，简直反人类，而高级玩家却觉得还不够难，不够过瘾。这对矛盾已经成为了影响STG发展的主要矛盾。<br><br></div><div><br>现今，大多数STG都会为了解决或缓解这一矛盾而采取了一些措施，比如大多数游戏会设置多个难度，使初学者玩简单难度，而高手可以玩疯狂难度；有的游戏会设置教学关卡，指导初学者避弹；有一些游戏设置了自动雷击功能来为初学者降低难度。<br><br></div><div><strong><br>然而就在最近，我们在游戏《弹幕音乐绘》中，提出了一种更加别出心裁的解决方案——自动避弹。<br></strong><br></div>
`
		res, err := extractText(a)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		// fmt.Println(res)
	})
}
