<!-- package=live.appblink.v1 -->
- [/xlive/app-blink/v1/topic/GetTopicList](#xliveapp-blinkv1topicGetTopicList) 获取话题列表
- [/xlive/app-blink/v1/topic/CheckTopic](#xliveapp-blinkv1topicCheckTopic) 检验话题是否有效

## /xlive/app-blink/v1/topic/GetTopicList
###获取话题列表

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "topic_list": [
            ""
        ]
    }
}
```


## /xlive/app-blink/v1/topic/CheckTopic
###检验话题是否有效

> 需要登录

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|platform|是|string||
|topic|是|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```

