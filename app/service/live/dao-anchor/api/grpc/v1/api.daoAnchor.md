<!-- package=live.daoanchor.v1 -->
- [/live.daoanchor.v1.DaoAnchor/FetchRoomByIDs](#live.daoanchor.v1.DaoAnchorFetchRoomByIDs)  FetchRoomByIDs 查询房间信息
- [/live.daoanchor.v1.DaoAnchor/RoomOnlineList](#live.daoanchor.v1.DaoAnchorRoomOnlineList)  RoomOnlineList 在线房间列表
- [/live.daoanchor.v1.DaoAnchor/RoomOnlineListByArea](#live.daoanchor.v1.DaoAnchorRoomOnlineListByArea)  RoomOnlineListByArea 分区在线房间列表(只返回room_id列表，不传分区，默认查找所有)
- [/live.daoanchor.v1.DaoAnchor/RoomOnlineListByAttrs](#live.daoanchor.v1.DaoAnchorRoomOnlineListByAttrs)  RoomOnlineListByAttrs 在线房间维度信息(不传attrs，不查询attr)
- [/live.daoanchor.v1.DaoAnchor/RoomCreate](#live.daoanchor.v1.DaoAnchorRoomCreate)  RoomCreate 房间创建
- [/live.daoanchor.v1.DaoAnchor/RoomUpdate](#live.daoanchor.v1.DaoAnchorRoomUpdate)  RoomUpdate 房间信息更新
- [/live.daoanchor.v1.DaoAnchor/RoomBatchUpdate](#live.daoanchor.v1.DaoAnchorRoomBatchUpdate)  RoomBatchUpdate 房间信息批量更新
- [/live.daoanchor.v1.DaoAnchor/RoomExtendUpdate](#live.daoanchor.v1.DaoAnchorRoomExtendUpdate)  RoomExtendUpdate 房间扩展信息更新
- [/live.daoanchor.v1.DaoAnchor/RoomExtendBatchUpdate](#live.daoanchor.v1.DaoAnchorRoomExtendBatchUpdate)  RoomExtendBatchUpdate 房间扩展信息批量更新
- [/live.daoanchor.v1.DaoAnchor/RoomExtendIncre](#live.daoanchor.v1.DaoAnchorRoomExtendIncre)  RoomExtendIncre 房间信息增量更新
- [/live.daoanchor.v1.DaoAnchor/RoomExtendBatchIncre](#live.daoanchor.v1.DaoAnchorRoomExtendBatchIncre)  RoomExtendBatchIncre 房间信息批量增量更新
- [/live.daoanchor.v1.DaoAnchor/RoomTagCreate](#live.daoanchor.v1.DaoAnchorRoomTagCreate)  RoomTagCreate 房间Tag创建
- [/live.daoanchor.v1.DaoAnchor/RoomAttrCreate](#live.daoanchor.v1.DaoAnchorRoomAttrCreate)  RoomAttrCreate 房间Attr创建
- [/live.daoanchor.v1.DaoAnchor/RoomAttrSetEx](#live.daoanchor.v1.DaoAnchorRoomAttrSetEx)  RoomAttrSetEx 房间Attr更新
- [/live.daoanchor.v1.DaoAnchor/AnchorUpdate](#live.daoanchor.v1.DaoAnchorAnchorUpdate)  AnchorUpdate 主播信息更新
- [/live.daoanchor.v1.DaoAnchor/AnchorBatchUpdate](#live.daoanchor.v1.DaoAnchorAnchorBatchUpdate)  AnchorBatchUpdate 主播信息批量更新
- [/live.daoanchor.v1.DaoAnchor/AnchorIncre](#live.daoanchor.v1.DaoAnchorAnchorIncre)  AnchorIncre 主播信息增量更新
- [/live.daoanchor.v1.DaoAnchor/AnchorBatchIncre](#live.daoanchor.v1.DaoAnchorAnchorBatchIncre)  AnchorBatchIncre 主播信息批量增量更新
- [/live.daoanchor.v1.DaoAnchor/FetchAreas](#live.daoanchor.v1.DaoAnchorFetchAreas)  FetchAreas 根据父分区号查询子分区
- [/live.daoanchor.v1.DaoAnchor/FetchAttrByIDs](#live.daoanchor.v1.DaoAnchorFetchAttrByIDs)  FetchAttrByIDs 批量根据房间号查询指标
- [/live.daoanchor.v1.DaoAnchor/DeleteAttr](#live.daoanchor.v1.DaoAnchorDeleteAttr)  DeleteAttr 删除某一个指标

## /live.daoanchor.v1.DaoAnchor/FetchRoomByIDs
### FetchRoomByIDs 查询房间信息

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_ids|否|多个integer||
|uids|否|多个integer||
|fields|否|多个string||
|default_fields|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "room_data_set": {
            "1": {
                "uid": 0,
                "room_id": 0,
                "short_id": 0,
                "title": "",
                "cover": "",
                "tags": "",
                "background": "",
                "description": "",
                "live_status": 0,
                "live_start_time": 0,
                "live_screen_type": 0,
                "live_mark": 0,
                "lock_status": 0,
                "lock_time": 0,
                "hidden_status": 0,
                "hidden_time": 0,
                "area_id": 0,
                "area_name": "",
                "parent_area_id": 0,
                "parent_area_name": "",
                "keyframe": "",
                "popularity_count": 0,
                "tag_list": [
                    {
                        "tag_id": 0,
                        "tag_sub_id": 0,
                        "tag_value": 0,
                        "tag_ext": "",
                        "tag_expire_at": 0
                    }
                ],
                "anchor_profile_type": 0,
                "anchor_level": {
                    //  当前等级
                    "level": 0,
                    //  当前等级颜色
                    "color": 0,
                    //  当前积分
                    "score": 0,
                    //  当前等级最小积分
                    "left": 0,
                    //  当前等级最大积分
                    "right": 0,
                    //  最大等级
                    "max_level": 0
                },
                "anchor_round_switch": 0,
                "anchor_round_status": 0,
                "anchor_record_switch": 0,
                "anchor_record_status": 0,
                "anchor_san": 0,
                //  0默认 1摄像头直播 2录屏直播 3语音直播
                "live_type": 0
            }
        }
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomOnlineList
### RoomOnlineList 在线房间列表

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|filter|否|string||
|sort|否|string||
|page|否|integer||
|page_size|否|integer||
|fields|否|多个string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "room_data_list": {
            "1": {
                "uid": 0,
                "room_id": 0,
                "short_id": 0,
                "title": "",
                "cover": "",
                "tags": "",
                "background": "",
                "description": "",
                "live_status": 0,
                "live_start_time": 0,
                "live_screen_type": 0,
                "live_mark": 0,
                "lock_status": 0,
                "lock_time": 0,
                "hidden_status": 0,
                "hidden_time": 0,
                "area_id": 0,
                "area_name": "",
                "parent_area_id": 0,
                "parent_area_name": "",
                "keyframe": "",
                "popularity_count": 0,
                "tag_list": [
                    {
                        "tag_id": 0,
                        "tag_sub_id": 0,
                        "tag_value": 0,
                        "tag_ext": "",
                        "tag_expire_at": 0
                    }
                ],
                "anchor_profile_type": 0,
                "anchor_level": {
                    //  当前等级
                    "level": 0,
                    //  当前等级颜色
                    "color": 0,
                    //  当前积分
                    "score": 0,
                    //  当前等级最小积分
                    "left": 0,
                    //  当前等级最大积分
                    "right": 0,
                    //  最大等级
                    "max_level": 0
                },
                "anchor_round_switch": 0,
                "anchor_round_status": 0,
                "anchor_record_switch": 0,
                "anchor_record_status": 0,
                "anchor_san": 0,
                //  0默认 1摄像头直播 2录屏直播 3语音直播
                "live_type": 0
            }
        }
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomOnlineListByArea
### RoomOnlineListByArea 分区在线房间列表(只返回room_id列表，不传分区，默认查找所有)

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|area_ids|否|多个integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "room_ids": [
            0
        ]
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomOnlineListByAttrs
### RoomOnlineListByAttrs 在线房间维度信息(不传attrs，不查询attr)

#### 方法：GET

#### 请求参数

```javascript
{
    "attrs": [
        {
            "attr_id": 0,
            "attr_sub_id": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "attrs": {
            "1": {
                "uid": 0,
                "room_id": 0,
                "area_id": 0,
                "parent_area_id": 0,
                "tag_list": [
                    {
                        "tag_id": 0,
                        "tag_sub_id": 0,
                        "tag_value": 0,
                        "tag_ext": "",
                        "tag_expire_at": 0
                    }
                ],
                "attr_list": [
                    {
                        "room_id": 0,
                        "attr_id": 0,
                        "attr_sub_id": 0,
                        "attr_value": 0
                    }
                ],
                "popularity_count": 0,
                "anchor_profile_type": 0
            }
        }
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomCreate
### RoomCreate 房间创建

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||
|room_id|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "room_id": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomUpdate
### RoomUpdate 房间信息更新

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|fields|是|多个string||
|room_id|是|integer||
|title|否|string||
|cover|否|string||
|tags|否|string||
|background|否|string||
|description|否|string||
|live_start_time|否|integer||
|live_screen_type|否|integer||
|lock_status|否|integer||
|lock_time|否|integer||
|hidden_time|否|integer||
|area_id|否|integer||
|anchor_round_switch|否|integer||
|anchor_record_switch|否|integer||
|live_type|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomBatchUpdate
### RoomBatchUpdate 房间信息批量更新

#### 方法：GET

#### 请求参数

```javascript
{
    "reqs": [
        {
            "fields": [
                ""
            ],
            "room_id": 0,
            "title": "",
            "cover": "",
            "tags": "",
            "background": "",
            "description": "",
            "live_start_time": 0,
            "live_screen_type": 0,
            "lock_status": 0,
            "lock_time": 0,
            "hidden_time": 0,
            "area_id": 0,
            "anchor_round_switch": 0,
            "anchor_record_switch": 0,
            "live_type": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomExtendUpdate
### RoomExtendUpdate 房间扩展信息更新

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|fields|是|多个string||
|room_id|是|integer||
|keyframe|否|string||
|danmu_count|否|integer||
|popularity_count|否|integer||
|audience_count|否|integer||
|gift_count|否|integer||
|gift_gold_amount|否|integer||
|gift_gold_count|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomExtendBatchUpdate
### RoomExtendBatchUpdate 房间扩展信息批量更新

#### 方法：GET

#### 请求参数

```javascript
{
    "reqs": [
        {
            "fields": [
                ""
            ],
            "room_id": 0,
            "keyframe": "",
            "danmu_count": 0,
            "popularity_count": 0,
            "audience_count": 0,
            "gift_count": 0,
            "gift_gold_amount": 0,
            "gift_gold_count": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomExtendIncre
### RoomExtendIncre 房间信息增量更新

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|req_id|是|string||
|fields|是|多个string||
|room_id|是|integer||
|danmu_count|否|integer||
|popularity_count|否|integer||
|audience_count|否|integer||
|gift_count|否|integer||
|gift_gold_amount|否|integer||
|gift_gold_count|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomExtendBatchIncre
### RoomExtendBatchIncre 房间信息批量增量更新

#### 方法：GET

#### 请求参数

```javascript
{
    "reqs": [
        {
            "req_id": "",
            "fields": [
                ""
            ],
            "room_id": 0,
            "danmu_count": 0,
            "popularity_count": 0,
            "audience_count": 0,
            "gift_count": 0,
            "gift_gold_amount": 0,
            "gift_gold_count": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomTagCreate
### RoomTagCreate 房间Tag创建

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_id|是|integer||
|tag_id|是|integer||
|tag_sub_id|否|integer||
|tag_value|否|integer||
|tag_ext|否|string||
|tag_expire_at|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomAttrCreate
### RoomAttrCreate 房间Attr创建

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_id|是|integer||
|attr_id|是|integer||
|attr_sub_id|否|integer||
|attr_value|否|integer||
|attr_ext|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/RoomAttrSetEx
### RoomAttrSetEx 房间Attr更新

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_id|是|integer||
|attr_id|是|integer||
|attr_sub_id|否|integer||
|attr_value|否|integer||
|attr_ext|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/AnchorUpdate
### AnchorUpdate 主播信息更新

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|fields|是|多个string||
|uid|是|integer||
|profile_type|否|integer||
|san_score|否|integer||
|round_status|否|integer||
|record_status|否|integer||
|exp|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/AnchorBatchUpdate
### AnchorBatchUpdate 主播信息批量更新

#### 方法：GET

#### 请求参数

```javascript
{
    "reqs": [
        {
            "fields": [
                ""
            ],
            "uid": 0,
            "profile_type": 0,
            "san_score": 0,
            "round_status": 0,
            "record_status": 0,
            "exp": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/AnchorIncre
### AnchorIncre 主播信息增量更新

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|req_id|是|string||
|fields|是|多个string||
|uid|是|integer||
|san_score|否|integer||
|exp|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/AnchorBatchIncre
### AnchorBatchIncre 主播信息批量增量更新

#### 方法：GET

#### 请求参数

```javascript
{
    "reqs": [
        {
            "req_id": "",
            "fields": [
                ""
            ],
            "uid": 0,
            "san_score": 0,
            "exp": 0
        }
    ]
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```


## /live.daoanchor.v1.DaoAnchor/FetchAreas
### FetchAreas 根据父分区号查询子分区

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|area_id|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "info": {
            "area_id": 0,
            "area_name": ""
        },
        "areas": [
            {
                "area_id": 0,
                "area_name": ""
            }
        ]
    }
}
```


## /live.daoanchor.v1.DaoAnchor/FetchAttrByIDs
### FetchAttrByIDs 批量根据房间号查询指标

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_ids|是|多个integer||
|attr_id|是|integer||
|attr_sub_id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "attrs": {
            "1": {
                "room_id": 0,
                "attr_id": 0,
                "attr_sub_id": 0,
                "attr_value": 0
            }
        }
    }
}
```


## /live.daoanchor.v1.DaoAnchor/DeleteAttr
### DeleteAttr 删除某一个指标

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|attr_id|是|integer||
|attr_sub_id|是|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

