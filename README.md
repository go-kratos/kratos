### Kratos

##### 项目规范
1,每个目录 需要有独立的README.md  CHANGELOG.md CONTRIBUTORS.md，具体可以参考：
http://git.bilibili.co/platform/go-common/tree/master/business/service/archive

2,以后每个业务或者基础组件维护自己的版本号，在CHANGELOG.md中，rider 构建以后的tag关联成自己的版本号；

3,整个大仓库不再有tag，只有master 主干分支，所有mr发送前，一定要注意先merge master；

4,使用Rider构建以后（retag），回滚可以基于Rider的retag来回滚，而不是回滚大仓库的代码；

5,提供RPC内部服务放置在business/service中，任务队列放置在business/job中，对外网关服务放置在business/interface，管理后台服务放置在business/admin

6,每个业务自建cmd文件夹,将main.go文件和test配置文件迁移进去

7,构建的时候自定义脚本选择krotos_buil.sh,自定义参数选择自己所在业务的路径 （ps：例如 interface/web-show）

8,大仓库的mr合并方式为，在mr中留言"+merge"，鉴权依据服务根目录下 CONTRIBUTORS.md 文件解析，具体可以参考：
http://info.bilibili.co/pages/viewpage.action?pageId=7539410

## 负责人信息
<details>
<summary>展开查看</summary>
<pre><code>.
├── Owner: maojian,haoguanwei
├── app
│   ├── Owner: maojian,haoguanwei,linmiao
│   ├── admin
│   │   ├── ep
│   │   │   ├── merlin
│   │   │   │   └── Owner: maojian,yuanmin,fengyifeng,xuneng
│   │   │   └── saga
│   │   │       └── Owner: tangyongqiang
│   │   ├── main
│   │   │   ├── activity
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei
│   │   │   ├── answer
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── apm
│   │   │   │   └── Owner: haoguanwei,lintanghui
│   │   │   ├── app
│   │   │   │   └── Owner: haoguanwei,peiyifei
│   │   │   ├── appstatic
│   │   │   │   └── Owner: liweijia,renwei
│   │   │   ├── bfs-apm
│   │   │   │   └── Owner: wangweizhen
│   │   │   ├── block
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── cache
│   │   │   │   └── Owner: lintanghui
│   │   │   ├── config
│   │   │   │   └── Owner: haoguanwei,lintanghui
│   │   │   ├── coupon
│   │   │   │   └── Owner: yubaihai,zhaogangtao
│   │   │   ├── creative
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── credit
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── dm
│   │   │   │   └── Owner: liangkai,renwei
│   │   │   ├── esports
│   │   │   │   └── Owner: liweijia,renwei
│   │   │   ├── filter
│   │   │   │   └── Owner: zhaogangtao,muyang
│   │   │   ├── growup
│   │   │   │   └── Owner: gaopeng
│   │   │   ├── laser
│   │   │   │   └── Owner: haoguanwei,shencen,wangzhe01
│   │   │   ├── manager
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei
│   │   │   ├── member
│   │   │   │   └── Owner: linmiao,haoguanwei,zhoujiahui,zhoujixiang,chenjianrong
│   │   │   ├── point
│   │   │   │   └── Owner: yubaihai,zhaogangtao
│   │   │   ├── push
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei
│   │   │   ├── relation
│   │   │   │   └── Owner: linmiao,zhoujiahui
│   │   │   ├── reply
│   │   │   │   └── Owner: chenzhihui,lujinhui
│   │   │   ├── search
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei,guanhuaxin
│   │   │   ├── sms
│   │   │   │   └── Owner: renwei,zhapuyu
│   │   │   ├── spy
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── tag
│   │   │   │   └── Owner: renwei,renyashun
│   │   │   ├── tv
│   │   │   │   └── Owner: liweijia,renwei
│   │   │   ├── up
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── upload
│   │   │   │   └── Owner: haoguanwei,zhapuyu
│   │   │   ├── usersuit
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── videoup
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── videoup-task
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── vip
│   │   │   │   └── Owner: zhaogangtao
│   │   │   └── workflow
│   │   │       └── Owner: haoguanwei,zhapuyu,zhuangzhewei,zhoushuguang
│   │   └── openplatform
│   │       └── sug
│   │           └── Owner: changxuanran,xucheng
│   ├── common
│   │   └── openplatform
│   │       └── Owner: liuzhan,huangshancheng
│   ├── interface
│   │   ├── live
│   │   │   ├── Owner: liuzhen
│   │   │   └── push-live
│   │   │       └── Owner: kuangxibin
│   │   └── main
│   │       ├── account
│   │       │   └── Owner: wanghuan01,zhoujiahui,zhaogangtao,chenjianrong,zhoujixiang
│   │       ├── activity
│   │       │   └── Owner: liweijia
│   │       ├── answer
│   │       │   └── Owner: zhaogangtao
│   │       ├── app-channel
│   │       │   └── Owner: peiyifei
│   │       ├── app-feed
│   │       │   └── Owner: peiyifei
│   │       ├── app-interface
│   │       │   └── Owner: peiyifei
│   │       ├── app-player
│   │       │   └── Owner: peiyifei
│   │       ├── app-resource
│   │       │   └── Owner: peiyifei
│   │       ├── app-show
│   │       │   └── Owner: peiyifei
│   │       ├── app-tag
│   │       │   └── Owner: peiyifei
│   │       ├── app-view
│   │       │   └── Owner: peiyifei
│   │       ├── app-wall
│   │       │   └── Owner: peiyifei
│   │       ├── article
│   │       │   └── Owner: changxuanran,lijiadong,qiuliang
│   │       ├── broadcast
│   │       │   └── Owner: chenzhihui,caoguoliang,guhao
│   │       ├── captcha
│   │       │   └── Owner: chenzhihui
│   │       ├── creative
│   │       │   └── Owner: shencen,wangzhe01
│   │       ├── credit
│   │       │   └── Owner: zhaogangtao
│   │       ├── dm
│   │       │   └── Owner: liangkai,renwei
│   │       ├── dm2
│   │       │   └── Owner: liangkai,renwei
│   │       ├── esports
│   │       │   └── Owner: liweijia,zhapuyu
│   │       ├── favorite
│   │       │   └── Owner: chenzhihui,lujinhui
│   │       ├── feedback
│   │       │   └── Owner: peiyifei
│   │       ├── growup
│   │       │   └── Owner: gaopeng
│   │       ├── history
│   │       │   └── Owner: renwei,wangxu01
│   │       ├── kvo
│   │       │   └── Owner: liweijia,zhapuyu
│   │       ├── laser
│   │       │   └── Owner: haoguanwei,shencen
│   │       ├── player
│   │       │   └── Owner: liweijia,zhapuyu
│   │       ├── playlist
│   │       │   └── Owner: liweijia
│   │       ├── push
│   │       │   └── Owner: renwei,zhapuyu
│   │       ├── push-archive
│   │       │   └── Owner: zhapuyu,shencen,renwei,liweijia,wangzhe01
│   │       ├── reply
│   │       │   └── Owner: lujinhui,chenzhihui,caoguoliang
│   │       ├── report-click
│   │       │   └── Owner: zhangshengchao,chenzhihui,renyashun
│   │       ├── shorturl
│   │       │   └── Owner: peiyifei,zhapuyu
│   │       ├── space
│   │       │   └── Owner: liweijia,zhapuyu
│   │       ├── spread
│   │       │   └── Owner: zhapuyu,renwei
│   │       ├── tag
│   │       │   └── Owner: renwei,renyashun
│   │       ├── tv
│   │       │   └── Owner: renwei,liweijia
│   │       ├── upload
│   │       │   └── Owner: peiyifei,zhapuyu
│   │       ├── videoup
│   │       │   └── Owner: shencen,wangzhe01
│   │       ├── web
│   │       │   └── Owner: liweijia,zhapuyu
│   │       ├── web-feed
│   │       │   └── Owner: zhapuyu,liweijia,renwei
│   │       ├── web-goblin
│   │       │   └── Owner: liweijia,renwei
│   │       └── web-show
│   │           └── Owner: liweijia
│   ├── job
│   │   ├── live
│   │   │   ├── Owner: liuzhen
│   │   │   └── wallet
│   │   │       └── Owner: lixiang,zhouzhichao
│   │   ├── main
│   │   │   ├── account-notify
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── account-summary
│   │   │   │   └── Owner: zhoujiahui
│   │   │   ├── activity
│   │   │   │   └── Owner: liweijia
│   │   │   ├── answer
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── app
│   │   │   │   └── Owner: peiyifei
│   │   │   ├── app-wall
│   │   │   │   └── Owner: peiyifei,renwei,haoguanwei
│   │   │   ├── archive
│   │   │   │   └── Owner: peiyifei
│   │   │   ├── archive-kisjd
│   │   │   │   └── Owner: peiyifei
│   │   │   ├── article
│   │   │   │   └── Owner: qiuliang,changxuanran,lijiadong
│   │   │   ├── block
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── broadcast
│   │   │   │   └── Owner: chenzhihui,caoguoliang,guhao
│   │   │   ├── click
│   │   │   │   └── Owner: peiyifei
│   │   │   ├── coin
│   │   │   │   └── Owner: lintanghui,linmiao,zhapuyu
│   │   │   ├── coupon
│   │   │   │   └── Owner: zhaogangtao,yubaihai
│   │   │   ├── creative
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── credit
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── credit-timer
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── dm
│   │   │   │   └── Owner: liangkai,renwei
│   │   │   ├── dm2
│   │   │   │   └── Owner: liangkai,renwei
│   │   │   ├── favorite
│   │   │   │   └── Owner: lujinhui,chenzhihui
│   │   │   ├── feed
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei
│   │   │   ├── figure
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── figure-timer
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── growup
│   │   │   │   └── Owner: gaopeng
│   │   │   ├── history
│   │   │   │   └── Owner: renwei,wangxu01
│   │   │   ├── identify
│   │   │   │   └── Owner: linmiao,wanghuan01
│   │   │   ├── member
│   │   │   │   └── Owner: chenjianrong,zhoujiahui,linmiao,zhoujixiang
│   │   │   ├── passport
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── passport-auth
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── passport-encrypt
│   │   │   │   └── Owner: linmiao
│   │   │   ├── passport-game-cloud
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── passport-game-data
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── passport-game-local
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── playlist
│   │   │   │   └── Owner: liweijia
│   │   │   ├── point
│   │   │   │   └── Owner: yubaihai,zhaogangtao
│   │   │   ├── push
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei
│   │   │   ├── relation
│   │   │   │   └── Owner: linmiao,zhoujiahui
│   │   │   ├── reply
│   │   │   │   └── Owner: chenzhihui,lujinhui,caoguoliang
│   │   │   ├── search
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei,guanhuaxin
│   │   │   ├── sms
│   │   │   │   └── Owner: renwei,zhapuyu
│   │   │   ├── spy
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── stat
│   │   │   │   └── Owner: peiyifei
│   │   │   ├── tag
│   │   │   │   └── Owner: renwei,renyashun
│   │   │   ├── tv
│   │   │   │   └── Owner: renwei,liweijia
│   │   │   ├── upload
│   │   │   │   └── Owner: zhapuyu,renwei,zhuangzhewei
│   │   │   ├── usersuit
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── videoup
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── videoup-report
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── vip
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── web-goblin
│   │   │   │   └── Owner: liweijia,renwei
│   │   │   └── workflow
│   │   │       └── Owner: haoguanwei,zhapuyu
│   │   └── openplatform
│   │       └── open-market
│   │           └── Owner: changxuanran,liuyan02,qiuliang
│   ├── service
│   │   ├── ep
│   │   │   └── saga-agent
│   │   │       └── Owner: muyang,tangyongqiang,fangrongchang
│   │   ├── live
│   │   │   ├── Owner: liuzhen
│   │   │   ├── userexp
│   │   │   │   └── Owner: kuangxibing
│   │   │   └── wallet
│   │   │       └── Owner: lixiang,zhouzhichao
│   │   ├── main
│   │   │   ├── account
│   │   │   │   └── Owner: wanghuan01,zhoujiahui
│   │   │   ├── antispam
│   │   │   │   └── Owner: chenzhihui,lujinhui
│   │   │   ├── archive
│   │   │   │   └── Owner: haoguanwei,peiyifei
│   │   │   ├── assist
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── block
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── bns
│   │   │   │   └── Owner: haoguawnei weicheng
│   │   │   ├── broadcast
│   │   │   │   └── Owner: chenzhihui,caoguoliang,guhao
│   │   │   ├── canal
│   │   │   │   └── Owner: haoguanwei
│   │   │   ├── coin
│   │   │   │   └── Owner: lintanghui,linmiao,zhapuyu
│   │   │   ├── config
│   │   │   │   └── Owner: maojian
│   │   │   ├── coupon
│   │   │   │   └── Owner: zhaogangtao,yubaihai
│   │   │   ├── dapper
│   │   │   │   └── Owner: maojian,haoguanwei
│   │   │   ├── databus
│   │   │   │   └── Owner: haoguanwei
│   │   │   ├── discovery
│   │   │   │   └── Owner: haoguanwei,lintanghui
│   │   │   ├── dynamic
│   │   │   │   └── Owner: liweijia,zhapuyu
│   │   │   ├── favorite
│   │   │   │   └── Owner: chenzhihui,lujinhui
│   │   │   ├── feed
│   │   │   │   └── Owner: renwei,zhapuyu
│   │   │   ├── figure
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── filter
│   │   │   │   └── Owner: zhaogangtao,muyang
│   │   │   ├── identify
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── identify-game
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── location
│   │   │   │   └── Owner: peiyifei,haoguanwei
│   │   │   ├── member
│   │   │   │   └── Owner: zhaogangtao,wanghuan01,zhoujiahui,chenjianrong,zhoujixiang
│   │   │   ├── msm
│   │   │   │   └── Owner: maojian
│   │   │   ├── notify
│   │   │   │   └── Owner: haoguanwei,lintanghui
│   │   │   ├── passport
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── passport-auth
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── passport-game
│   │   │   │   └── Owner: wanghuan01
│   │   │   ├── point
│   │   │   │   └── Owner: yubaihai,zhaogangtao
│   │   │   ├── push
│   │   │   │   └── Owner: renwei,zhapuyu
│   │   │   ├── push-strategy
│   │   │   │   └── Owner: renwei,zhapuyu
│   │   │   ├── relation
│   │   │   │   └── Owner: linmiao,zhoujiahui
│   │   │   ├── resource
│   │   │   │   └── Owner: haoguanwei,peiyifei
│   │   │   ├── search
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei,guanhuaxin
│   │   │   ├── secure
│   │   │   │   └── Owner: zhaogangtao,lintanghui
│   │   │   ├── seq-server
│   │   │   │   └── Owner: peiyifei
│   │   │   ├── share
│   │   │   │   └── Owner: peiyifei,haoguanwei
│   │   │   ├── sms
│   │   │   │   └── Owner: renwei,zhapuyu
│   │   │   ├── spy
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── tag
│   │   │   │   └── Owner: renwei,renyashun
│   │   │   ├── thumbup
│   │   │   │   └── Owner: liweijia,zhapuyu,renwei
│   │   │   ├── up
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── upcredit
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── usersuit
│   │   │   │   └── Owner: zhaogangtao
│   │   │   ├── videoup
│   │   │   │   └── Owner: shencen,wangzhe01
│   │   │   ├── vip
│   │   │   │   └── Owner: lintanghui,zhaogangtao
│   │   │   └── workflow
│   │   │       └── Owner: haoguanwei,zhapuyu,zhoushuguang
│   │   └── openplatform
│   │       ├── abtest
│   │       │   └── Owner: lijiadong,qiuliang
│   │       ├── anti-fraud
│   │       │   └── Owner: wanglitao,wangminda,jiayanxiang
│   │       ├── ticket-item
│   │       │   └── Owner: yangyucheng
│   │       └── ticket-sales
│   │           └── Owner: liuzhan,yangyucheng
│   └── tool
│       ├── cache
│       │   └── Owner: zhapuyu
│       ├── ci
│       │   └── Owner: tangyongqiang
│       ├── creater
│       │   └── Owner: chenjianrong
│       ├── gdoc
│       │   └── Owner: lintanghui
│       ├── saga
│       │   └── Owner: muyang,tangyongqiang
│       └── warden
│           └── Owner: weicheng
└── library
    ├── cache
    │   ├── memcache
    │   │   └── Owner: maojian
    │   └── redis
    │       └── Owner: maojian
    ├── container
    │   └── pool
    │       └── Owner: zhapuyu
    ├── database
    │   ├── elastic
    │   │   └── Owner: haoguanwei,renwei,zhapuyu
    │   └── sql
    │       └── Owner: 
    ├── ecode
    │   ├── Owner: all
    │   └── tip
    │       └── Owner: all
    ├── exp
    │   └── feature
    │       └── Owner: zhoujiahui
    ├── log
    │   └── Owner: maojian
    ├── naming
    │   └── discovery
    │       └── Owner: lintanghui,caoguoliang
    ├── net
    │   ├── http
    │   │   ├── Owner: maojian
    │   │   └── blademaster
    │   │       ├── Owner: 
    │   │       ├── middleware
    │   │       │   ├── Owner: 
    │   │       │   ├── antispam
    │   │       │   │   └── Owner: 
    │   │       │   ├── auth
    │   │       │   │   └── Owner: maojian,zhoujiahui
    │   │       │   ├── cache
    │   │       │   │   └── Owner: 
    │   │       │   ├── identify
    │   │       │   │   └── Owner: 
    │   │       │   ├── limit
    │   │       │   │   └── aqm
    │   │       │   │       └── Owner: 
    │   │       │   ├── proxy
    │   │       │   │   └── Owner: 
    │   │       │   ├── rate
    │   │       │   │   └── Owner: 
    │   │       │   ├── supervisor
    │   │       │   │   └── Owner: 
    │   │       │   ├── tag
    │   │       │   │   └── Owner: maojian
    │   │       │   └── verify
    │   │       │       └── Owner: maojian,zhoujiahui
    │   │       └── render
    │   │           └── Owner: 
    │   ├── metadata
    │   │   └── Owner: 
    │   ├── netutil
    │   │   └── breaker
    │   │       └── Owner: 
    │   ├── rpc
    │   │   └── warden
    │   │       ├── Owner: maojian,caoguoliang
    │   │       ├── balancer
    │   │       │   └── wrr
    │   │       │       └── Owner: caoguoliang
    │   │       └── resolver
    │   │           └── Owner: caoguoliang
    │   └── trace
    │       └── Owner: maojian
    ├── rate
    │   └── limit
    │       └── bench
    │           └── stress
    │               └── Owner: lintanghui
    ├── stat
    │   └── sys
    │       └── cpu
    │           └── Owner: caoguoliang
    └── sync
        └── errgroup
            └── Owner: 
</code></pre>
</details>
