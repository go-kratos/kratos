##

`GET http://api.live.bilibili.com/xlive/gift/v1/gift/room_gift_list`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|room_id|否|integer||
|area_v2_parent_id|否|integer||
|area_v2_id|否|integer||
|platform|否|string||
|build|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
            {
                "id": 0,
                "position": 0,
                "plan_id": 0
            }
        ],
        "silver_list": [
            {
                "id": 0,
                "position": 0,
                "plan_id": 0
            }
        ],
        "show_count_map": 0,
        "old_list": [
            {
                "id": 0,
                "name": "",
                "price": 0,
                "type": 0,
                "coin_type": {
                    "mapKey": ""
                },
                "img": "",
                "gift_url": "",
                "count_set": "",
                "combo_num": 0,
                "super_num": 0,
                "count_map": {
                    "1": ""
                }
            }
        ]
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/gift/v1/gift/gift_config`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string||
|build|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "data": [
            {
                "id": 0,
                "name": "",
                "price": 0,
                "type": 0,
                "coin_type": "",
                "bag_gift": 0,
                "effect": 0,
                "corner_mark": "",
                "broadcast": 0,
                "draw": 0,
                "stay_time": 0,
                "animation_frame_num": 0,
                "desc": "",
                "rule": "",
                "rights": "",
                "privilege_required": 0,
                "count_map": [
                    {
                        "num": 0,
                        "text": ""
                    }
                ],
                "img_basic": "",
                "img_dynamic": "",
                "frame_animation": "",
                "gif": "",
                "webp": "",
                "full_sc_web": "",
                "full_sc_horizontal": "",
                "full_sc_vertical": "",
                "full_sc_horizontal_svga": "",
                "full_sc_vertical_svga": "",
                "bullet_head": "",
                "bullet_tail": ""
            }
        ]
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/gift/v1/gift/discount_gift_list`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer||
|roomid|否|integer||
|area_v2_parent_id|否|integer||
|area_v2_id|否|integer||
|platform|否|string||
|ruid|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "discount_list": [
            {
                "gift_id": 0,
                "price": 0,
                "discount_price": 0,
                "corner_mark": "",
                "corner_position": 0,
                "corner_color": ""
            }
        ]
    }
}
```

##

`GET http://api.live.bilibili.com/xlive/gift/v1/gift/daily_bag`

### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer||

```json
{
    "code": 0,
    "message": "ok",
    "data": {
        "bag_status": 0,
        "bag_expire_status": 0,
        "bag_toast": {
            "toast_status": 0,
            "toast_message": ""
        },
        "bag_list": [
            {
                "type": 0,
                "bag_name": "",
                "gift_list": [
                    {
                        "gift_id": 0,
                        "gift_num": 0,
                        "expire_at": 0
                    }
                ]
            }
        ]
    }
}
```

