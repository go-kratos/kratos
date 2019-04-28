<!-- package=live.appucenter.v1 -->
- [/xlive/app-ucenter/v1/topic/GetTopicList](#xliveapp-ucenterv1topicGetTopicList) 获取话题列表
- [/xlive/app-ucenter/v1/topic/CheckTopic](#xliveapp-ucenterv1topicCheckTopic) 检验话题是否有效

## /xlive/app-ucenter/v1/topic/GetTopicList
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


## /xlive/app-ucenter/v1/topic/CheckTopic
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

