<!-- package=live.rtc.v1 -->
- [/live.rtc.v1.Rtc/JoinChannel](#live.rtc.v1.RtcJoinChannel) 
- [/live.rtc.v1.Rtc/LeaveChannel](#live.rtc.v1.RtcLeaveChannel) 
- [/live.rtc.v1.Rtc/PublishStream](#live.rtc.v1.RtcPublishStream) 
- [/live.rtc.v1.Rtc/TerminateStream](#live.rtc.v1.RtcTerminateStream) 
- [/live.rtc.v1.Rtc/Channel](#live.rtc.v1.RtcChannel) 
- [/live.rtc.v1.Rtc/Stream](#live.rtc.v1.RtcStream) 
- [/live.rtc.v1.Rtc/SetRtcConfig](#live.rtc.v1.RtcSetRtcConfig) 
- [/live.rtc.v1.Rtc/VerifyToken](#live.rtc.v1.RtcVerifyToken) 

## /live.rtc.v1.Rtc/JoinChannel
### 无标题

#### 方法：POST

#### 请求参数

```javascript
{
    "channel_id": 0,
    "user_id": 0,
    "proto_version": 0,
    "source": [
        {
            "type": 0,
            "codec": "",
            "media_specific": "",
            "ssrc": 0,
            "user_id": 0
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
        "call_id": 0,
        "token": "",
        "source": [
            {
                "type": 0,
                "codec": "",
                "media_specific": "",
                "ssrc": 0,
                "user_id": 0
            }
        ]
    }
}
```


## /live.rtc.v1.Rtc/LeaveChannel
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|channel_id|否|integer||
|user_id|否|integer||
|call_id|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.rtc.v1.Rtc/PublishStream
### 无标题

#### 方法：POST

#### 请求参数

```javascript
{
    "channel_id": 0,
    "user_id": 0,
    "call_id": 0,
    "encoder_config": {
        "width": 0,
        "height": 0,
        "bitrate": 0,
        "frame_rate": 0,
        "video_codec": "",
        "video_profile": "",
        "channel": 0,
        "sample_rate": 0,
        "audio_codec": ""
    },
    "mix_config": ""
}
```

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.rtc.v1.Rtc/TerminateStream
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|channel_id|否|integer||
|user_id|否|integer||
|call_id|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.rtc.v1.Rtc/Channel
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|channel_id|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "media_source": [
            {
                "type": 0,
                "codec": "",
                "media_specific": "",
                "ssrc": 0,
                "user_id": 0
            }
        ],
        "server": "",
        "tcp_port": 0,
        "udp_port": 0
    }
}
```


## /live.rtc.v1.Rtc/Stream
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|channel_id|否|integer||
|user_id|否|integer||
|call_id|否|integer||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "encoder_config": {
            "width": 0,
            "height": 0,
            "bitrate": 0,
            "frame_rate": 0,
            "video_codec": "",
            "video_profile": "",
            "channel": 0,
            "sample_rate": 0,
            "audio_codec": ""
        },
        "mix_config": ""
    }
}
```


## /live.rtc.v1.Rtc/SetRtcConfig
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|channel_id|否|integer||
|user_id|否|integer||
|call_id|否|integer||
|config|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
    }
}
```


## /live.rtc.v1.Rtc/VerifyToken
### 无标题

#### 方法：GET

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|channel_id|否|integer||
|call_id|否|integer||
|token|否|string||

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "pass": true
    }
}
```

