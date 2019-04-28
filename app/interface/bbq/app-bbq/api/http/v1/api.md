# HTTP API文档

导航
---

* [版本](#版本)
* [说明](#说明)
* [接口说明](#接口说明)
    + [[首页视频列表]](#首页视频列表)
    + [[批量视频地址]](#批量视频地址接口)

版本
---

| 版本 | 时间       | 修订者 | 备注           |
| :--- | :--------- | :----- | :------------- |
| v0.1 | 2018.10.12 | Cheney | 初始化接口列表 |

说明
---

BBQ视频相关接口

域名
---

* **线上：** bbq.bilibili.com

* **自动化文档路径:** http://172.16.38.91/swagger


接口说明
-------

### 首页视频列表

***接口说明***

提供视频列表相关基础信息和播放信息

***请求URL***

http://DOMAIN/bbq/app-bbq/sv/list

***请求方式***

HTTP/GET

***请求参数***

| 参数名   | 必选 | 类型 | 说明                    |
| :------- | :--- | :--- | :---------------------- |
| qn       | 是   | int  | 首选清晰度              |
| pagesize | 是   | int  | 返回数据量（range:0-20) |

***qn参数枚举***

目前PlayUrl支持UGC、PGC和Mobile H5的业务, 对各个业务的支持情况如下:  
***<span style="color:red">Note: 对于分组为同一个的清晰度, PlayUrl服务端只会给出其中一个清晰度作为结果返回.***  
##### UGC
| 分组 | 清晰度qn | 类型 | 名称       | 描述         | 访问权限             |
| :--- | :------- | :--- | :--------- | :----------- | :------------------- |
| 1#   | 116      | FLV  | flv_p60    | 高清 1080P60 | 大会员或此视频的UP主 |
| 1#   | 112      | FLV  | hdflv2     | 高清 1080P+  | 大会员或此视频的UP主 |
| 2#   | 74       | FLV  | flv720_p60 | 高清 720P60  | 大会员或此视频的UP主 |
| 3#   | 80       | FLV  | flv        | 高清 1080    | 都可访问             |
| 4#   | 64       | FLV  | flv720     | 高清 720P    | 都可访问             |
| 4#   | 48       | MP4  | hdmp4      | 高清 720P    | 都可访问             |
| 5#   | 32       | FLV  | flv480     | 清晰 480P    | 都可访问             |
| 6#   | 15       | FLV  | flv360     | 流畅 360P    | 都可访问             |
| 6#   | 16       | MP4  | mp4        | 流畅 360P    | 都可访问             |
| 7#   | 6        | MP4  | mp4        | 极速 240P    | 当且仅当type=mp4访问 |

##### PGC
| 分组 | 清晰度qn | 类型 | 名称   | 描述        | 访问权限             |
| :--- | :------- | :--- | :----- | :---------- | :------------------- |
| 1#   | 112      | FLV  | hdflv2 | 高清 1080P+ | 都可访问             |
| 2#   | 80       | FLV  | flv    | 高清 1080   | 都可访问             |
| 3#   | 64       | FLV  | flv720 | 高清 720P   | 都可访问             |
| 3#   | 48       | MP4  | hdmp4  | 高清 720P   | 都可访问             |
| 4#   | 32       | FLV  | flv480 | 清晰 480P   | 都可访问             |
| 5#   | 15       | FLV  | flv360 | 流畅 360P   | 都可访问             |
| 5#   | 16       | MP4  | mp4    | 流畅 360P   | 都可访问             |
| 6#   | 6        | MP4  | mp4    | 极速 240P   | 当且仅当type=mp4访问 |

##### Mobile H5
| 分组 | 清晰度qn | 类型 | 名称   | 描述      | 访问权限             |
| :--- | :------- | :--- | :----- | :-------- | :------------------- |
| 1#   | 15       | FLV  | flv360 | 流畅 360P | 都可访问             |
| 1#   | 16       | MP4  | mp4    | 流畅 360P | 都可访问             |
| 2#   | 6        | MP4  | mp4    | 极速 240P | 当且仅当type=mp4访问 |


*** 返回字段说明 ***

| 字段名    | 字段类型                 | 字段说明         |
| :-------: | :----------------------: | :--------------: |
| svid      | int                      | bbq视频id        |
| title     | string                   | 视频标题         |
| mid       | int                      | 发布用户uid      |
| duration  | int                      | 时长             |
| pubtime   | string                   | 发布时间         |
| ctime     | int                      | 创建时间         |
| avid      | int                      | avid             |
| cid       | int                      | cid              |
| from      | int                      | 来源渠道         |
| tag       | string                   | 后端首选tag      |
| tags      | [][VideoTag](#VideoTag)  | 全部tag          |
| pic       | string                   | 图片预留字段     |
| like      | int                      | 点赞数           |
| reply     | int                      | 评论数           |
| share     | int                      | 分享数           |
| user_info | *UserCard(object)        | 发布用户信息     |
| play      | *[VideoPlay](#VideoPlay) | 播放信息         |
| is_like   | bool                     | 当前用户是否点赞 |


***返回字段示例***

```javascript
{
    "code": 0,
    "message": "0",
    "ttl": 1,
    "data": [
        {
            "svid": 198,
            "title": "体重88kg,深蹲200KG成功了,用了2年时间,记录一下",
            "content": "",
            "mid": 98628543,
            "duration": 0,
            "pubtime": "2018-10-12 12:45:51",
            "ctime": 1539319551,
            "avid": 27406598,
            "cid": 47265445,
            "from": 0,
            "tag": "",
            "tags": [],
            "pic": "",
            "like": 0,
            "reply": 0,
            "share": 0,
            "user_info": {
                "mid": 98628543,
                "name": "笑圣贤",
                "sex": "保密",
                "rank": 10000,
                "face": "http://i0.hdslb.com/bfs/face/fac0c8dc176ce9f8d3e176d759378233788994d4.jpg",
                "sign": "一个练力量举的UP,1981年的大叔,体重88kg,三大项540KG",
                "level": 0,
                "vip_info": {
                    "type": 0,
                    "status": 0,
                    "due_date": 0
                }
            },
            "play": {
                "cid": 47265445,
                "expire_time": 1539330912,
                "file_info": [
                    {
                        "ahead": "Egg=",
                        "filesize": 2332329,
                        "timelength": 40128,
                        "vhead": "AWQAHv/hABpnZAAerNlBcFHlkhAAAAMAEAAAAwMg8WLZYAEABWjr7PI8"
                    }
                ],
                "fnval": 0,
                "fnver": 0,
                "quality": 15,
                "support_description": [
                    "高清 720P",
                    "清晰 480P",
                    "流畅 360P"
                ],
                "support_formats": [
                    "flv720",
                    "flv480",
                    "flv360"
                ],
                "support_quality": [
                    64,
                    32,
                    15
                ],
                "url": "http://upos-hz-mirrorcos.acgvideo.com/upgcxcode/45/54/47265445/47265445-1-15.flv?um_deadline=1539334512&platform=&rate=98807&oi=0&um_sign=45b60b850c9e577261feb079ced1bb3d&gen=playurl&os=cos&trid=3db463de4aa2476cb56900e87f733323",
                "video_codecid": 7,
                "video_project": true,
                "current_time": 1539327312
            },
            "is_like": false
        },
        {
            "svid": 1740,
            "title": "【朱一龙舔屏】全部都是你(๑> ᴗ<๑)❤【竖屏丨手机福利】",
            "content": "手机观看效果最佳！送闺蜜@灰 的居仔舔屏视频，熬夜终于在七夕搞定了，就是想摸鱼剪剪竖屏玩。祝大家七夕情人节快乐！快点进来收获爱情吧！ヾ(๑╹◡╹)ﾉ\"",
            "mid": 761997,
            "duration": 0,
            "pubtime": "2018-10-12 12:52:06",
            "ctime": 1539319926,
            "avid": 29574294,
            "cid": 51427230,
            "from": 0,
            "tag": "",
            "tags": [
                {
                    "id": 797,
                    "name": "明星",
                    "type": 3
                },
                {
                    "id": 798,
                    "name": "娱乐",
                    "type": 1
                },
                {
                    "id": 805,
                    "name": "朱一龙",
                    "type": 3
                },
                {
                    "id": 1729,
                    "name": "剪辑",
                    "type": 3
                },
                {
                    "id": 4103,
                    "name": "七夕",
                    "type": 3
                }
            ],
            "pic": "",
            "like": 0,
            "reply": 0,
            "share": 0,
            "user_info": {
                "mid": 761997,
                "name": "無駄無駄",
                "sex": "保密",
                "rank": 10000,
                "face": "http://i2.hdslb.com/bfs/face/fbe53ad204dd36d268a9e6e887918787956914fd.gif",
                "sign": "【渣渣水平，蜗牛手速，感谢赏脸！】微博：@-6plus7-   【B站ID是抢注的舍不得改，没错，是隐藏的JOJO厨！】",
                "level": 0,
                "vip_info": {
                    "type": 2,
                    "status": 1,
                    "due_date": 1570809600000
                }
            },
            "play": {
                "cid": 51427230,
                "expire_time": 1539330912,
                "file_info": [
                    {
                        "ahead": "EZA=",
                        "filesize": 6204058,
                        "timelength": 102254,
                        "vhead": "AWQAHv/hABlnZAAerNlBcFHl4QAAAwABAAADADwPFi2WAQAFaOvs8jw="
                    }
                ],
                "fnval": 0,
                "fnver": 0,
                "quality": 15,
                "support_description": [
                    "高清 1080P60",
                    "高清 720P60",
                    "高清 1080P",
                    "高清 720P",
                    "清晰 480P",
                    "流畅 360P"
                ],
                "support_formats": [
                    "flv_p60",
                    "flv720_p60",
                    "flv",
                    "flv720",
                    "flv480",
                    "flv360"
                ],
                "support_quality": [
                    116,
                    74,
                    80,
                    64,
                    32,
                    15
                ],
                "url": "http://upos-hz-mirrorcos.acgvideo.com/upgcxcode/30/72/51427230/51427230-1-15.flv?um_deadline=1539334512&platform=&rate=103144&oi=0&um_sign=305264b6bc32d61fd1020a71c24f3544&gen=playurl&os=cos&trid=3db463de4aa2476cb56900e87f733323",
                "video_codecid": 7,
                "video_project": true,
                "current_time": 1539327312
            },
            "is_like": false
        }
    ]
}
```

### 批量视频地址接口

***接口说明***

提供可替换视频地址服务

***请求URL***

http://DOMAIN/bbq/app-bbq/sv/playlist

***请求方式***

HTTP/GET

***请求参数***

| 参数名 | 必选 | 类型   | 说明                |
| :----- | :--- | :----- | :------------------ |
| qn     | 是   | int    | 首选清晰度          |
| cids   | 是   | string | 请求cid（逗号分隔） |

*** 返回字段说明 ***

| 字段名 | 字段类型                   | 字段说明          |
| :----: | :------------------------: | :---------------: |
| data   | []*[VideoPlay](#VideoPlay) | bvc playurl 结构 |

***返回字段示例***

```javascript
{
    "code": 0,
    "message": "0",
    "ttl": 1,
    "data": [
        {
            "cid": 49159587,
            "expire_time": 1539331184,
            "file_info": [
                {
                    "ahead": "EZBW5QA=",
                    "filesize": 12879106,
                    "timelength": 258118,
                    "vhead": "AWQAH//hAB1nZAAfrNnA2D3n//AoACfxAAADA+kAAOpgDxgxngEABWjpuyyL"
                }
            ],
            "fnval": 0,
            "fnver": 0,
            "quality": 32,
            "support_description": [
                "高清 1080P",
                "高清 720P",
                "清晰 480P",
                "流畅 360P"
            ],
            "support_formats": [
                "flv",
                "flv720",
                "flv480",
                "flv360"
            ],
            "support_quality": [
                80,
                64,
                32,
                15
            ],
            "url": "http://upos-hz-mirrorcos.acgvideo.com/upgcxcode/87/95/49159587/49159587-1-32.flv?um_deadline=1539334784&platform=&rate=84823&oi=0&um_sign=2969b624cf3abf67505e60d78babbc70&gen=playurl&os=cos&trid=60941200ead244fbbcbd6090c5b337d0",
            "video_codecid": 7,
            "video_project": true,
            "current_time": 1539327585
        },
        {
            "cid": 52172330,
            "expire_time": 1539331184,
            "file_info": [
                {
                    "ahead": "EZA=",
                    "filesize": 2377519,
                    "timelength": 15181,
                    "vhead": "AWQAHv/hABlnZAAerNlA2D3n4QAAAwABAAADADIPFi2WAQAFaOvs8jw="
                }
            ],
            "fnval": 0,
            "fnver": 0,
            "quality": 32,
            "support_description": [
                "高清 1080P+",
                "高清 1080P",
                "高清 720P",
                "清晰 480P",
                "流畅 360P"
            ],
            "support_formats": [
                "hdflv2",
                "flv",
                "flv720",
                "flv480",
                "flv360"
            ],
            "support_quality": [
                112,
                80,
                64,
                32,
                15
            ],
            "url": "http://upos-hz-mirrorcos.acgvideo.com/upgcxcode/30/23/52172330/52172330-1-32.flv?um_deadline=1539334784&platform=&rate=266239&oi=0&um_sign=1533acfdd923bffab96f65d63f64ad9c&gen=playurl&os=cos&trid=60941200ead244fbbcbd6090c5b337d0",
            "video_codecid": 7,
            "video_project": true,
            "current_time": 1539327585
        }
    ]
}
```


*附
--

***通用字段***

##### VideoPlay

| 字段名              | 字段类型                 | 字段说明                                            |
| :-----------------: | :----------------------: | :-------------------------------------------------: |
| cid                 | int64                    | cid                                                 |
| expire_time         | int64                    | 过期时间                                            |
| file_info           | []*[FileInfo](#FileInfo) | 分片信息                                            |
| fnval               | int64                    | 播放器请求端使用的, 功能标识, 每位(为1)标识一个功能 |
| fnver               | int64                    | 播放器请求端使用的, 功能版本号                      |
| quality             | int64                    | 清晰度                                              |
| support_description | []string                 | 支持清晰度描述                                      |
| support_formats     | []string                 | 支持格式                                            |
| support_quality     | []int64                  | 支持清晰度                                          |
| url                 | string                   | 基础url                                             |
| video_codecid       | int64                    | 对应视频调度的结果编码值                            |
| video_project       | bool                     | 是否可投影                                          |
| current_time        | int64                    | 当前时间戳                                          |

##### FileInfo

| 字段名     | 字段类型 | 字段说明                             |
| :--------: | :------: | :----------------------------------: |
| ahead      | int      | 视频分片文件的音频头信息, BASE64编码 |
| filesize   | int      | 视频的对应分片文件的大小             |
| timelength | int      | 视频的对应分片文件的时长             |
| vhead      | int      | 视频分片文件的视频头信息, BASE64编码 |

##### VideoTag

| 字段名 | 字段类型 | 字段说明 |
| :----: | :------: | :------: |
| id     | int      | tagid    |
| name   | int      | tag名字  |
| type   | int      | tag类型  |
