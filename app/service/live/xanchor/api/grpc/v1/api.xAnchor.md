## FetchRoomByIDs 查询房间信息

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/FetchRoomByIDs`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_ids|否|多个integer||
|uids|否|多个integer||
|fields|否|多个string||
|default_fields|否|integer||

```json
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
                "room_tag_list": [
                    {
                        "tag_id": 0,
                        "tag_type": 0,
                        "tag_value": 0,
                        "tag_attribute": ""
                    }
                ],
                "anchor_tag_list": [
                    {
                        "tag_id": 0,
                        "tag_type": 0,
                        "tag_value": 0,
                        "tag_attribute": ""
                    }
                ],
                "anchor_profile_type": 0,
                "anchor_exp": [
                    {
                        "level": 0,
                        "next_level": 0,
                        "level_color": 0,
                        "exp": 0,
                        "current_level_exp": 0,
                        "next_level_exp": 0
                    }
                ],
                "anchor_round_switch": 0,
                "anchor_round_status": 0,
                "anchor_record_switch": 0,
                "anchor_record_status": 0,
                "anchor_san": 0,
                "live_type": 0
            }
        }
    }
}
```

## RoomOnlineList 在线房间列表

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomOnlineList`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|filter|否|string||
|sort|否|string||
|page|否|integer||
|page_size|否|integer||
|fields|否|多个string||

```json
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
                "room_tag_list": [
                    {
                        "tag_id": 0,
                        "tag_type": 0,
                        "tag_value": 0,
                        "tag_attribute": ""
                    }
                ],
                "anchor_tag_list": [
                    {
                        "tag_id": 0,
                        "tag_type": 0,
                        "tag_value": 0,
                        "tag_attribute": ""
                    }
                ],
                "anchor_profile_type": 0,
                "anchor_exp": [
                    {
                        "level": 0,
                        "next_level": 0,
                        "level_color": 0,
                        "exp": 0,
                        "current_level_exp": 0,
                        "next_level_exp": 0
                    }
                ],
                "anchor_round_switch": 0,
                "anchor_round_status": 0,
                "anchor_record_switch": 0,
                "anchor_record_status": 0,
                "anchor_san": 0,
                "live_type": 0
            }
        }
    }
}
```

## RoomCreate 房间创建

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomCreate`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer||
|room_id|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "room_id": 0
    }
}
```

## RoomUpdate 房间信息更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomUpdate`

### 请求参数

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

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## RoomBatchUpdate 房间信息批量更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomBatchUpdate`

### 请求参数

```json
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

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## RoomExtendUpdate 房间扩展信息更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomExtendUpdate`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|fields|是|多个string||
|room_id|是|integer||
|key_frame|否|string||
|danmu_count|否|integer||
|popularity_count|否|integer||
|audience_count|否|integer||
|gift_count|否|integer||
|gift_gold_amount|否|integer||
|gift_gold_count|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## RoomExtendBatchUpdate 房间扩展信息批量更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomExtendBatchUpdate`

### 请求参数

```json
{
    "reqs": [
        {
            "fields": [
                ""
            ],
            "room_id": 0,
            "key_frame": "",
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

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## RoomExtendIncre 房间信息增量更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomExtendIncre`

### 请求参数

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

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## RoomExtendBatchIncre 房间信息批量增量更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomExtendBatchIncre`

### 请求参数

```json
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

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## RoomTagSet 房间Tag更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/RoomTagSet`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|fields|是|多个string||
|room_id|是|integer||
|tag_type|是|integer||
|tag_value|否|integer||
|tag_attribute|否|string||
|tag_expire_at|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## AnchorUpdate 主播信息更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/AnchorUpdate`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|fields|是|多个string||
|uid|是|integer||
|profile_type|否|integer||
|san_score|否|integer||
|round_status|否|integer||
|record_status|否|integer||
|exp|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## AnchorBatchUpdate 主播信息批量更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/AnchorBatchUpdate`

### 请求参数

```json
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

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## AnchorIncre 主播信息增量更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/AnchorIncre`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|req_id|是|string||
|fields|是|多个string||
|uid|是|integer||
|san_score|否|integer||
|exp|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## AnchorBatchIncre 主播信息批量增量更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/AnchorBatchIncre`

### 请求参数

```json
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

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

## AnchorTagSet 主播Tag更新

`GET http://api.live.bilibili.com/xlive/xanchor/v1/xAnchor/AnchorTagSet`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|fields|是|多个string||
|anchor_id|是|integer||
|tag_type|是|integer||
|tag_value|否|integer||
|tag_attribute|否|string||
|tag_expire_at|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "affected_rows": 0
    }
}
```

