<!-- package=live.liveadmin.v1 -->
- [/xlive/live-admin/v1/upload/file](#xlivelive-adminv1uploadfile) 

## /xlive/live-admin/v1/upload/file
### 无标题

#### 方法：POST

#### 请求参数

|参数名|必选|类型|描述|
|:---|:---|:---|:---|
|bucket|是|string| 上传到 BFS 的 bucket|
|dir|否|string| 上传到指定的 BFS 目录（可以用来区分业务）|
|filename|否|string| 上传的到bfs的文件名（存储在bfs的文件名，不传bfs会根据文件的sha1值生成并返回）|
|contentType|否|string| 上传的文件的类型（不指定时会自动检测文件类型）|
|wmKey|否|string| 图片水印key，添加图片水印需要上传该参数, 新业务需要提前向bfs申请|
|wmText|否|string| 文字水印，限制不超过20个字符|
|wmPaddingX|否|integer| 水印位置右下角 到原图右下角 水平距离，默认10px|
|wmPaddingY|否|integer| 水印位置右下角 到原图右下角 垂直距离，默认10px|
|wmScale|否|float| 水印宽度占原图高度的比例(0,1) （只支持按照宽度压缩)，默认值: 0.035|
|token|是|string| 上传 Token，通过 obtainToken 接口获取|

#### 响应

```javascript
{
    "code": 0,
    "message": "ok",
    "data": {
        "url": ""
    }
}
```

