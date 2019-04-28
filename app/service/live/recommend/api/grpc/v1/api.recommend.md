## 获取n个推荐, 得到的结果是在线的房间
 去重，不会重复推荐
 如果没有足够推荐的结果则返回空的结果，调用方需要补位

`GET http://api.live.bilibili.com/xlive/recommend/v1/recommend/random_recs_by_user`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid|
|count|否|integer| 获取数量|
|exist_ids|否|多个integer| room_id去重|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        //  返回数量
        "count": 0,
        //  房间id
        "room_ids": [
            0
        ]
    }
}
```

## 清空推荐缓存，清空推荐过的集合

`GET http://api.live.bilibili.com/xlive/recommend/v1/recommend/clear_recommend_cache`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户uid|

```json
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

