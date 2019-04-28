# dapper-collector 收集 trace 数据写入 influxdb 与 HBase

### Influxdb 存储格式

| measurement | tags                                               | fields                                        |
|-------------|----------------------------------------------------|-----------------------------------------------|
| span        | service_name,operation_name,peer.service,span.kind | max_duration,min_duration,avg_duration,errors |

### HBase 存储格式

dapper:listidx

| rowkey                                                            | cf:kind:d:{duration nanosecond} | cf:kind:e:{span_id} |
|-------------------------------------------------------------------|----------------------------------------|--------------------|
| hex(hash({service_name})hex(hash({operation_name}))){timestamp/5} | hex({trace_id}):hex({span_id})         | hex({trace_id})    |

```
create 'dapper:listidx', {NAME=>'kind', VERSION=>1, TTL=>604800}
```

dapper:rawtrace

| rowkey          | cf:pb:hex({span_id})_{c,s} |
|-----------------|--------------------------------|
| hex({trace_id}) | protobuf({raw_data})           |

```
create 'dapper:rawtrace', {NAME=>'pb', VERSION=>1, TTL=>604800}
```
