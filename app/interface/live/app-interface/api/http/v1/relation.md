## [app端关注二级页][全量]正在直播接口

`GET http://api.live.bilibili.com/xlive/app-interface/v1/relation/liveAnchor`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|buyaofangqizhiliao|否|string| 调试咒语|
|platform|否|string| 平台|
|device|否|string| 设备|
|build|否|string| 版本号|
|sortRule|否|integer| 排序类型|
|filterRule|否|integer| 筛选类型|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "rooms": [
            {
                //  房间id
                "roomid": 0,
                //  用户id
                "uid": 0,
                //  用户昵称
                "uname": "",
                //  用户头像
                "face": "",
                //  直播间标题
                "title": "",
                //  直播间标签
                "live_tag_name": "",
                //  开始直播时间
                "live_time": 0,
                //  人气值
                "online": 0,
                //  秒开url
                "playurl": "",
                //  可选清晰度
                "accept_quality": [
                    0
                ],
                //  当前清晰度
                "current_quality": 0,
                //  pk_id
                "pk_id": 0,
                //  特别关注标志
                "special_attention": 0,
                //  老的分区id
                "area": 0,
                //  老的分区名
                "area_name": "",
                //  子分区id
                "area_v2_id": 0,
                //  子分区名
                "area_v2_name": "",
                //  父分区名
                "area_v2_parent_name": "",
                //  父分区id
                "area_v2_parent_id": 0,
                //  广播适配标志
                "broadcast_type": 0,
                //  官方认证标志
                "official_verify": 0,
                //  直播间跳转链接
                "link": "",
                //  直播间封面
                "cover": "",
                //  角标文字
                "pendent_ru": "",
                //  角标颜色
                "pendent_ru_color": "",
                //  角标背景图
                "pendent_ru_pic": ""
            }
        ],
        "total_count": 0,
        "card_type": 0,
        "big_card_type": 0
    }
}
```

## [app端关注二级页][分页]暂未开播接口

`GET http://api.live.bilibili.com/xlive/app-interface/v1/relation/unliveAnchor`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|buyaofangqizhiliao|否|string| 调试咒语|
|page|否|integer| 分页号|
|pagesize|否|integer| 页大小|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "rooms": [
            {
                //  上次开播描述
                "live_desc": "",
                //  房间id
                "roomid": 0,
                //  用户id
                "uid": 0,
                //  用户昵称
                "uname": "",
                //  用户头像
                "face": "",
                //  特别关注标志
                "special_attention": 0,
                //  官方认证标志
                "official_verify": 0,
                //  直播状态标志
                "live_status": 0,
                //  广播适配标志
                "broadcast_type": 0,
                //  老的分区id
                "area": 0,
                //  粉丝数
                "attentions": 0,
                //  老的分区名
                "area_name": "",
                //  子分区id
                "area_v2_id": 0,
                //  子分区名
                "area_v2_name": "",
                //  父分区名
                "area_v2_parent_name": "",
                //  父分区id
                "area_v2_parent_id": 0,
                //  直播间跳转链接
                "link": "",
                //  房间页公告
                "announcement_content": "",
                //  房间页公告发布时间
                "announcement_time": ""
            }
        ],
        "total_count": 0,
        "no_room_count": 0,
        "has_more": 0
    }
}
```

