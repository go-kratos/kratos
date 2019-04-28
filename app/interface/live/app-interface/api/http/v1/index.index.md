## 首页大接口
 首页换一换接口

`GET http://api.live.bilibili.com/xlive/app-interface/v1/index/getAllList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string|平台|
|device|是|string|设备|
|scale|是|string|分辨率|
|build|是|integer|版本号|
|relation_page|是|integer|关注页码|
|module_id|否|integer|模块id（可选）|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/app-interface/v1/index/change`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|module_id|是|integer| 模块id|
|attention_room_id|是|string||
|platform|否|string| 平台|
|device|否|string|设备|
|scale|否|string|分辨率|
|build|否|integer|版本号|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "module_list": [
            {
                "module_info": {
                    //  模块id
                    "id": 0,
                    //  标题
                    "title": "",
                    //  图标
                    "pic": "",
                    //  list数据类型  1: banner 2: 导航栏 3: 运营推荐分区-标准 4: 运营推荐分区-方 5：排行榜（小时榜） 6: 推荐主播-标准 7: 推荐主播-方 8:我的关注(用户相关) 9：一级分区-标准 10：一级分区-方 11: 活动卡片 12：常用标签推荐入口(用户相关) 13：常用标签推荐房间列表(用户相关) 14：大航海提示入口
                    "type": 0,
                    //  跳转链接
                    "link": "",
                    //  该模块数据总数
                    "count": 0,
                    "is_sky_horse_gray": 0
                },
                //  注意：可能是 PicList{id,pic,link,title}，需要根据ModuleInfo里的type判断
                "list": [
                    {
                        "roomid": 0,
                        "title": "",
                        "uname": "",
                        "online": 0,
                        "cover": "",
                        "link": "",
                        "face": "",
                        "area_v2_parent_id": 0,
                        "area_v2_parent_name": "",
                        "area_v2_id": 0,
                        "area_v2_name": "",
                        "play_url,omitempty": "",
                        "play_url_h265,omitempty": "",
                        "current_quality,omitempty": 0,
                        "broadcast_type": 0,
                        "pendent_ru": "",
                        "pendent_ru_pic": "",
                        "pendent_ru_color": "",
                        "rec_type": 0,
                        "pk_id": 0,
                        "accept_quality,omitempty": [
                            0
                        ]
                    }
                ]
            }
        ]
    }
}
```

