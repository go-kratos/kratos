<!-- package=live.xuserex.v1 -->
- [/live.xuserex.v1.RoomNotice/buy_guard](#live.xuserex.v1.RoomNoticebuy_guard)  是否弹出大航海购买提示
- [/live.xuserex.v1.RoomNotice/is_task_finish](#live.xuserex.v1.RoomNoticeis_task_finish)  habse 任务是否结束
- [/live.xuserex.v1.RoomNotice/set_task_finish](#live.xuserex.v1.RoomNoticeset_task_finish)  手动设置base 任务结束

## /live.xuserex.v1.RoomNotice/buy_guard
### 是否弹出大航海购买提示

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|uid|是|integer| UID|
|target_id|是|integer| 主播UID|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否提示 1:弹出提示,0:不弹出
        "should_notice": 0,
        //  状态有效开始时间
        "begin": 0,
        //  状态有效结束时间
        "end": 0,
        //  当前的时间戳
        "now": 0,
        //  标题
        "title": "",
        //  内容
        "content": "",
        //  按钮
        "button": ""
    }
}
```


## /live.xuserex.v1.RoomNotice/is_task_finish
### habse 任务是否结束

#### 方法：GET

#### 请求参数


#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否完成
        "result": 0
    }
}
```


## /live.xuserex.v1.RoomNotice/set_task_finish
### 手动设置base 任务结束

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|result|否|integer| 是否完成|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        //  是否完成
        "result": 0
    }
}
```

