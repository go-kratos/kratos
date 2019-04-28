## 根据uid查询直播关键历史记录
> 需要登录

`GET http://api.live.bilibili.com/xlive/web-ucenter/v1/history/get_history_by_uid`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "title": "",
        "count": 0,
        "list": [
            {
                "roomid": 0,
                "uid": 0,
                "uname": "",
                "user_cover": "",
                "title": "",
                "face": "",
                "tags": "",
                "live_status": 0,
                "fans_num": 0,
                "is_attention": 0,
                "area_v2_id": 0,
                "area_v2_name": "",
                "area_v2_parent_name": "",
                "area_v2_parent_id": 0
            }
        ]
    }
}
```

## 删除直播历史记录
> 需要登录

`POST http://api.live.bilibili.com/xlive/web-ucenter/v1/history/del_history`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

