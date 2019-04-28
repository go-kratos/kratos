<!-- package=live.webucenter.v1 -->
- [/xlive/web-ucenter/v1/history/get_history_by_uid](#xliveweb-ucenterv1historyget_history_by_uid)  根据uid查询直播关键历史记录
- [/xlive/web-ucenter/v1/history/del_history](#xliveweb-ucenterv1historydel_history)  删除直播历史记录

## /xlive/web-ucenter/v1/history/get_history_by_uid
### 根据uid查询直播关键历史记录

> 需要登录

#### 方法：GET

#### 请求参数


#### 响应

```javascript
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


## /xlive/web-ucenter/v1/history/del_history
### 删除直播历史记录

> 需要登录

#### 方法：POST

#### 请求参数


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

