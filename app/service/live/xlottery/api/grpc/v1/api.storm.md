<!-- package=live.xlottery.v1 -->
- [/live.xlottery.v1.Storm/Start](#live.xlottery.v1.StormStart)  开启节奏风暴
- [/live.xlottery.v1.Storm/CanStart](#live.xlottery.v1.StormCanStart) 节奏风暴是否能开启
- [/live.xlottery.v1.Storm/Join](#live.xlottery.v1.StormJoin) 加入节奏风暴
- [/live.xlottery.v1.Storm/Check](#live.xlottery.v1.StormCheck) 检查是否加入节奏风暴 

## /live.xlottery.v1.Storm/Start
### 开启节奏风暴

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户id|
|ruid|否|integer| 主播id|
|roomid|否|integer|房间号|
|useShield|否|bool|是否开启敏感词过滤|
|num|否|integer|道具数量|
|beatid|否|integer|节奏内容id|
|skipExternalCheck|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  错误码
        "code": 0,
        //  错误信息
        "msg": "",
        "start": {
            // 倒计时,秒
            "time": 0,
            // 抽奖标识
            "id": 0
        }
    }
}
```


## /live.xlottery.v1.Storm/CanStart
###节奏风暴是否能开启

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|否|integer| 用户id|
|ruid|否|integer| 主播id|
|roomid|否|integer|房间号|
|useShield|否|bool|是否开启敏感词过滤|
|num|否|integer|道具数量|
|beatid|否|integer|节奏内容id|
|skipExternalCheck|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  错误码
        "code": 0,
        //  错误信息
        "msg": "",
        // 是否能开启节奏风暴
        "ret_status": true
    }
}
```


## /live.xlottery.v1.Storm/Join
###加入节奏风暴

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|id|否|integer|抽奖id|
|roomid|否|integer|房间id|
|color|否|string|弹幕颜色 |
|mid|否|integer|userid  |
|platform|否|string|平台 web，ios，android|
|captcha_token|否|string|验证码标识|
|captcha_phrase|否|string|验证码明文 |

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  错误码
        "code": 0,
        //  错误信息
        "msg": "",
        //  加入成功返回数据
        "join": {
            // 礼物id
            "gift_id": 0,
            // 标题
            "title": "",
            // 礼物web内容
            "content": "",
            // 礼物移动端内容
            "mobile_content": "",
            // 礼物图片
            "gift_img": "",
            // 礼物数量
            "gift_num": 0,
            // 礼物名字
            "gift_name": ""
        }
    }
}
```


## /live.xlottery.v1.Storm/Check
###检查是否加入节奏风暴 

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|roomid|是|integer|房间号|
|uid|否|integer|用户id|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  错误码
        "code": 0,
        //  错误信息
        "msg": "",
        "check": {
            // 用户id
            "id": 0,
            // 房间号
            "roomid": 0,
            // 数量
            "num": 0,
            // 发送数量
            "send_num": "",
            // 结束时间戳
            "time": 0,
            // 内容
            "content": "",
            // 是否已经加入
            "hasJoin": 0,
            // 图片链接
            "storm_gif": ""
        }
    }
}
```

