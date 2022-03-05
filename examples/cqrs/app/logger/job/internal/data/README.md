# Data

## 注意

**Kafka的客户端**
- 不采用 [Sarama](https://github.com/Shopify/sarama)   
- 采用 [Confluent](https://github.com/confluentinc/confluent-kafka-go/kafka)  

**为什么?** 

Sarama Go客户端存在以下已知问题：

- 当Topic新增分区时，Sarama Go客户端无法感知并消费新增分区，需要客户端重启后，才能消费到新增分区。
- 当Sarama Go客户端同时订阅两个以上的Topic时，有可能会导致部分分区无法正常消费消息。
- 当Sarama Go客户端的消费位点重置策略设置为```Oldest(earliest)```时，如果客户端宕机或服务端版本升级，由于Sarama Go客户端自行实现OutOfRange机制，有可能会导致客户端从最小位点开始重新消费所有消息。

参考阿里云的这篇文章:  
- [为什么不推荐使用Sarama Go客户端收发消息？](https://help.aliyun.com/document_detail/266782.html)
- [关于 Kafka 应用开发知识点的整理](https://pandaychen.github.io/2022/01/01/A-KAFKA-USAGE-SUMUP-2/)
