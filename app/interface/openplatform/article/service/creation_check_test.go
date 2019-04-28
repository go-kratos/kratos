package service

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_contentStripSize(t *testing.T) {
	Convey("emoji size 2", t, func() {
		size := s.contentStripSize("😀")
		So(size, ShouldEqual, 2)
	})
	Convey("chinese and english", t, func() {
		size := s.contentStripSize("中a")
		So(size, ShouldEqual, 2)
	})
	Convey("ignore normal blank char", t, func() {
		size := s.contentStripSize("中 a \n \t b")
		So(size, ShouldEqual, 3)
	})
	Convey("img size 1", t, func() {
		size := s.contentStripSize("<img></img>")
		So(size, ShouldEqual, 1)
	})
	Convey("truely data", t, func() {
		data := `<p>音乐卡：</p><figure class="img-box" contenteditable="false"><img src="//uat-i0.hdslb.com/bfs/article/0aae45bcb008157ba5c7765ab8d952284d12fcad.png" aid="au75" width="1320" height="188" class="music-card" type="normal"/></figure><p>商品卡：</p><figure class="img-box" contenteditable="false"><img src="//uat-i0.hdslb.com/bfs/article/999065dfd84193ecbbd590a6a6fd46a374d2a840.png" aid="sp886" width="1320" height="208" class="shop-card" type="normal"/></figure><p>票务卡：</p><figure class="img-box" contenteditable="false"><img src="//uat-i0.hdslb.com/bfs/article/458aec77c8523fcb5e846b128e68804ee875cc26.png" aid="pw100" width="1320" height="208" class="shop-card" type="normal"/></figure><p><br/></p>`
		size := s.contentStripSize(data)
		So(size, ShouldEqual, 15)
	})
	Convey("truely data 2", t, func() {
		data := `<p><br/></p><figure class="img-box" contenteditable="false"><img src="//i0.hdslb.com/bfs/article/690a4cdd2d652c04b32aa737f9653895b909c8da.png" width="745" height="289"/><figcaption class="caption" contenteditable="true">-</figcaption></figure><p><br/></p><p><br/></p><p><span class="font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <span class="color-blue-04 font-size-14">海奉是一个风景优美的地方，但并不在沿海。数量众多的旅行家笔记显示，海奉是一片死火山群。那里坐落着世界上最高的山峰——奈文摩尔峰，峰顶终年积雪。其它沉睡的火山围坐在他的周围，高低不同，错落有致。火山口往往积蓄湖水，形成湖泊，当地人称之为“镜湖”。每到雨季，经过连续的降雨，湖中的水便会溢出，从山顶冲下，形成“水山爆发”的情景。山脚下是海奉人的村落，那里的房子全部以木头搭建，巧妙的避开河水的必经之路。海奉人以木工闻名，无论是精巧的木头机械还是美丽的木雕都不在话下。此外，每一个海奉人都戴着一枚木制的十字架，那是由海奉独有的铁木制成，绝不出售给外人，因而成为海奉人的标志。</span></span></p><p><span class="font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 但是故事并不发生在海奉，这些描写仅是因为主角是海奉人。</span></p><p style="text-align: left;"><span class="font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; <span class="color-blue-04 font-size-14">船还在航行。天色昏暗，雨从来没有停过。船舱紧闭，窗口透出一丝微弱的光。</span></span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; “您是海奉人吗？”山本真奈美借着微弱的灯光盯着他的十字架。</span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; “我家乡在海奉。”任川吃着生鱼片，随手用一根铁钎拨弄油灯的灯芯。</span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; “您为什么离开家乡呢？”</span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; 任川没有回答。</span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; “我只是好奇，毕竟海奉是个风景如画的地方。”</span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; “那可是个十分可笑的理由，”任川叹了口气，“你不理解也没关系，请答应不要打断我吧。”</span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp; &nbsp; “嗯。”</span></p><p style="text-align: left;"><span class="color-blue-04 font-size-14">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;`
		size := s.contentStripSize(data)
		So(size, ShouldEqual, 496)
	})
	Convey("unicode blank char should eq 0", t, func() {
		So(s.contentStripSize("\u200B"), ShouldEqual, 0)
		So(s.contentStripSize("\u00a0"), ShouldEqual, 0)
	})
}
func Test_checkTitle(t *testing.T) {
	Convey("en 80 should ok", t, func() {
		title := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		res, ok := s.checkTitle(title)
		So(res, ShouldEqual, title)
		So(ok, ShouldBeTrue)
	})
	Convey("chinese 40 should ok", t, func() {
		title := "好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好"
		res, ok := s.checkTitle(title)
		So(res, ShouldEqual, title)
		So(ok, ShouldBeTrue)
	})
	Convey("chinese 30 and en 21 should be wrong", t, func() {
		title := "好好好好好好好好好好好好好好好好好好好好好好好好好好好好好好aaaaaaaaaaaaaaaaaaaaa"
		_, ok := s.checkTitle(title)
		So(ok, ShouldBeFalse)
	})
}
