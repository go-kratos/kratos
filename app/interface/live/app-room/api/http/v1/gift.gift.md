<!-- package=live.approom.v1 -->
- [/xlive/app-room/v1/gift/daily_bag](#xliveapp-roomv1giftdaily_bag)  每日礼包接口
- [/xlive/app-room/v1/gift/NeedTipRecharge](#xliveapp-roomv1giftNeedTipRecharge) 
- [/xlive/app-room/v1/gift/TipRechargeAction](#xliveapp-roomv1giftTipRechargeAction) 
- [/xlive/app-room/v1/gift/gift_config](#xliveapp-roomv1giftgift_config) 礼物全量配置

## /xlive/app-room/v1/gift/daily_bag
### 每日礼包接口

#### 方法：GET

#### 请求参数


#### 响应

```javascript
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
                "source": {
                    "medal_id": 0,
                    "medal_name": "",
                    "level": 0,
                    "user_level": 0
                },
                "gift_list": [
                    {
                        "gift_id": "",
                        "gift_num": 0,
                        "expire_at": 0
                    }
                ]
            }
        ]
    }
}
```


## /xlive/app-room/v1/gift/NeedTipRecharge
### 无标题

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|from|是|integer| 来源 1金瓜子 2 银瓜子|
|needGold|否|integer| 需要的金瓜子  如果From=2　那么直接传0|
|platform|是|string| 平台 android ios|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否展示
        "show": 0,
        //  bp
        "bp": 0.1,
        //  bp券
        "bpCoupon": 0.1,
        //  需要充值的金瓜子
        "rechargeGold": 0
    }
}
```


## /xlive/app-room/v1/gift/TipRechargeAction
### 无标题

> 需要登录

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|from|是|integer| 来源 1金瓜子 2 银瓜子|
|action|是|integer|行为 1 停止推送|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /xlive/app-room/v1/gift/gift_config
###礼物全量配置

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|否|string||
|build|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "list": [
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
                "bullet_tail": "",
                "limit_interval": 0
            }
        ]
    }
}
```

