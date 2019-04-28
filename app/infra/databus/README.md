#### archive-service

##### 项目简介
> 1.databus是一个通过使用redis协议来简化kafka的消费方/生产方的一个中间件

##### 编译环境
> 请只用golang v1.8.x以上版本编译执行。  

##### 依赖包
> 1.公共包go-common  
> 2.kafka包：sarama和sarama-cluster

##### 特别说明
> 1.databus服务只是使用redis协议，并不是和redis用法就完全一样，所以必须完全参照文档进行使用。
> 2.文档地址：http://info.bilibili.co/pages/viewpage.action?pageId=2491209

##### 特性

> 1. 采用kafka 0.10版本，使用kafka新协议做offset提交
> 2. auth 时使用dsn协议，key:secret@group/topic=?&role=?&offset=?
> 3. 采用新的redis命令进行交互，批量返回消息
> 4. 不再自动提交offset，由客户端手动提交partition的offset
> 5. databus新增集群的概念，一个appkey,secret只能用于授权一个databus集群
