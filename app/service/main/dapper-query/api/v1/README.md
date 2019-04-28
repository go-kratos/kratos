api/v1/api.proto
================
**Version:** version not set

### /x/internal/dapper/clt-status
---
##### ***GET***
**Summary:** CltStatus 获取 collector 信息

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1CltStatusReply](#v1cltstatusreply) |

### /x/internal/dapper/depends-rank
---
##### ***GET***
**Summary:** DependsRank 查询某一个 service_name:operation_name 下所有依赖组件排名

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| service_name | query |  | No | string |
| start | query |  | No | string (int64) |
| end | query |  | No | string (int64) |
| rank_type | query | 排序类型 max_duration 最大耗时, min_duration 最小耗时, avg_duration 平均耗时, errors 错误数. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1DependsRankReply](#v1dependsrankreply) |

### /x/internal/dapper/depends-topology
---
##### ***GET***
**Summary:** DependsTopology 获取依赖拓扑图

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1DependsTopologyReply](#v1dependstopologyreply) |

### /x/internal/dapper/list-span
---
##### ***GET***
**Summary:** ListSpan 列出一个 service_name 某一 operation_name 所有采样到 Span

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| service_name | query |  | No | string |
| operation_name | query |  | No | string |
| start | query |  | No | string (int64) |
| end | query |  | No | string (int64) |
| order | query | 目前支持的 order  time:desc time:asc 按时间排序 duration:desc duration:asc 按耗时排序. | No | string |
| only_error | query | 只显示 error 的 span. | No | boolean (boolean) |
| offset | query |  | No | integer |
| limit | query |  | No | integer |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1ListSpanReply](#v1listspanreply) |

### /x/internal/dapper/operation-names
---
##### ***GET***
**Summary:** ListOperationName 列出某一 service  下所有 operation_name 仅 span.kind 为 server 的 operation_name

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| service_name | query |  | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1ListOperationNameReply](#v1listoperationnamereply) |

### /x/internal/dapper/operation-names-rank
---
##### ***GET***
**Summary:** OperationNameRank 查询 OperationName 排名列表

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| service_name | query |  | No | string |
| start | query |  | No | string (int64) |
| end | query |  | No | string (int64) |
| rank_type | query | 排序类型 max_duration 最大耗时, min_duration 最小耗时, avg_duration 平均耗时, errors 错误数. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1OperationNameRankReply](#v1operationnamerankreply) |

### /x/internal/dapper/ops-log
---
##### ***GET***
**Summary:** OpsLog 获取 通过 trace-id 获取 opslog 记录
如果请求的 trace-id 没有被记录到, 则需要提供 service_name operation_name 和 timestamp 进行模糊查询

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| trace_id | query |  | No | string |
| span_id | query |  | No | string |
| trace_field | query |  | No | string |
| service_name | query |  | No | string |
| operation_name | query |  | No | string |
| start | query | 开始时间. | No | string (int64) |
| end | query | 结束时间. | No | string (int64) |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1OpsLogReply](#v1opslogreply) |

### /x/internal/dapper/raw-trace
---
##### ***GET***
**Summary:** RawTrace 原始 Trace 数据

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| trace_id | query |  | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1RawTraceReply](#v1rawtracereply) |

### /x/internal/dapper/sample-point
---
##### ***GET***
**Summary:** SamplePoint 获取采样点数据

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| service_name | query |  | No | string |
| operation_name | query |  | No | string |
| only_error | query | only_error 在 errors 那个图可以指定为 true. | No | boolean (boolean) |
| interval | query | interval 使用 span-series 返回的 interval 即可. | No | string (int64) |
| time | query | time 使用 time-series 返回的时间即可，相同格式型如 2006-01-02T15:04:05. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1SamplePointReply](#v1samplepointreply) |

### /x/internal/dapper/service-depend
---
##### ***GET***
**Summary:** ServiceDepend 查询服务的直接依赖
TODO: 通过最近收集的到3 个 span 实时计算的，在当前查询的服务出现不正常的时候，查询结果可能不准确

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| service_name | query | service_name 不解释!. | No | string |
| operation_name | query | operation_name 当 operation_name 为空时查询所有 operation_name 然后 merge 结果. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1ServiceDependReply](#v1servicedependreply) |

### /x/internal/dapper/service-names
---
##### ***GET***
**Summary:** ListServiceName 列出所有 service

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1ListServiceNameReply](#v1listservicenamereply) |

### /x/internal/dapper/span-series
---
##### ***GET***
**Summary:** SpanSeries 获取 span 的时间序列数据

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| service_name | query |  | No | string |
| operation_name | query |  | No | string |
| start | query |  | No | string (int64) |
| end | query |  | No | string (int64) |
| fields | query | 可选的 fields 有 max_duration, min_duration, avg_duration, errors 其中除 errors 返回的是一段时间内的总数 其他返回的都是平均数 fields 是个数组可以通过 fields=max_duration,min_duration,avg_duration 逗号分隔. | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1SpanSeriesReply](#v1spanseriesreply) |

### /x/internal/dapper/trace
---
##### ***GET***
**Summary:** Trace 查询一个 Trace

**Parameters**

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| trace_id | query |  | No | string |
| span_id | query |  | No | string |

**Responses**

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | A successful response. | [v1TraceReply](#v1tracereply) |

### Models
---

### v1Client  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| addr | string |  | No |
| err_count | string (int64) |  | No |
| rate | string (int64) |  | No |
| up_time | string (int64) |  | No |

### v1CltNode  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| node | string |  | No |
| queue_len | string (int64) |  | No |
| clients | [ [v1Client](#v1client) ] |  | No |

### v1CltStatusReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| nodes | [ [v1CltNode](#v1cltnode) ] |  | No |

### v1DependsRankReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| rank_type | string |  | No |
| items | [ [v1RankItem](#v1rankitem) ] |  | No |

### v1DependsTopologyItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| service_name | string |  | No |
| depend_on | string |  | No |

### v1DependsTopologyReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [v1DependsTopologyItem](#v1dependstopologyitem) ] |  | No |

### v1Field  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| key | string |  | No |
| value | string |  | No |

### v1ListOperationNameReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| operation_names | [ string ] |  | No |

### v1ListServiceNameReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| service_names | [ string ] |  | No |

### v1ListSpanReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [v1SpanListItem](#v1spanlistitem) ] |  | No |

### v1Log  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| timestamp | string (int64) |  | No |
| fields | [ [v1Field](#v1field) ] |  | No |

### v1OperationNameRankReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| rank_type | string |  | No |
| items | [ [v1RankItem](#v1rankitem) ] |  | No |

### v1OpsLogRecord  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| time | string |  | No |
| fields | object |  | No |
| level | string |  | No |
| message | string |  | No |

### v1OpsLogReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| records | [ [v1OpsLogRecord](#v1opslogrecord) ] |  | No |

### v1RankItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| service_name | string |  | No |
| operation_name | string |  | No |
| value | double |  | No |

### v1RawTraceReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [v1Span](#v1span) ] |  | No |

### v1SamplePointItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| trace_id | string |  | No |
| span_id | string |  | No |
| duration | string (int64) |  | No |
| is_error | boolean (boolean) |  | No |

### v1SamplePointReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [v1SamplePointItem](#v1samplepointitem) ] |  | No |

### v1SeriesItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| field | string |  | No |
| values | [ string (int64) ] |  | No |

### v1ServiceDependItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| service_name | string |  | No |
| component | string |  | No |
| operation_names | [ string ] |  | No |

### v1ServiceDependReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| items | [ [v1ServiceDependItem](#v1servicedependitem) ] |  | No |

### v1Span  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| service_name | string |  | No |
| operation_name | string |  | No |
| trace_id | string |  | No |
| span_id | string |  | No |
| parent_id | string |  | No |
| start_time | string (int64) |  | No |
| duration | string (int64) |  | No |
| tags | object |  | No |
| logs | [ [v1Log](#v1log) ] |  | No |
| level | integer |  | No |
| childs | [ [v1Span](#v1span) ] |  | No |

### v1SpanListItem  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| trace_id | string |  | No |
| span_id | string |  | No |
| parent_id | string |  | No |
| service_name | string |  | No |
| operation_name | string |  | No |
| start_time | string |  | No |
| duration | string |  | No |
| tags | object |  | No |
| is_error | boolean (boolean) |  | No |
| container_ip | string |  | No |
| region_zone | string |  | No |
| mark | string |  | No |

### v1SpanSeriesReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| interval | string (int64) |  | No |
| times | [ string ] |  | No |
| items | [ [v1SeriesItem](#v1seriesitem) ] |  | No |

### v1TagValue  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| string_value | string |  | No |
| int64_value | string (int64) |  | No |
| bool_value | boolean (boolean) |  | No |
| float_value | float |  | No |

### v1TraceReply  

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| service_count | integer |  | No |
| span_count | integer |  | No |
| max_level | integer |  | No |
| root | [v1Span](#v1span) |  | No |