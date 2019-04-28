/*
 Navicat Premium Data Transfer

 Source Server         : 172.22.33.22
 Source Server Type    : MySQL
 Source Server Version : 50633
 Source Host           : 172.22.33.22:3306
 Source Schema         : test

 Target Server Type    : MySQL
 Target Server Version : 50633
 File Encoding         : 65001

 Date: 17/12/2018 11:48:32
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for apply
-- ----------------------------
DROP TABLE IF EXISTS `apply`;
CREATE TABLE `apply` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `from` varchar(16) NOT NULL DEFAULT '' COMMENT '申请人',
  `path` varchar(50) NOT NULL DEFAULT '' COMMENT '服务树path',
  `to` varchar(16) NOT NULL DEFAULT '' COMMENT '操作人',
  `status` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '状态 -1 申请中 1 生效',
  `active` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '-1 無效 1生效',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `start_time` varchar(16) NOT NULL DEFAULT '' COMMENT '压测开始时间',
  `end_time` varchar(16) NOT NULL DEFAULT '' COMMENT '压测结束时间',
  PRIMARY KEY (`id`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='权限申请表';

-- ----------------------------
-- Table structure for client_moni
-- ----------------------------
DROP TABLE IF EXISTS `client_moni`;
CREATE TABLE `client_moni` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `script_id` int(11) NOT NULL DEFAULT '0' COMMENT '脚本ID',
  `report_su_id` int(11) NOT NULL DEFAULT '0' COMMENT '报告ID',
  `job_name` varchar(20) NOT NULL DEFAULT '' COMMENT '容器名',
  `job_name_all` varchar(25) NOT NULL DEFAULT '' COMMENT '容器全名',
  `cpu_used` varchar(25) NOT NULL DEFAULT '' COMMENT 'cpu使用率',
  `elapsd_time` int(11) NOT NULL DEFAULT '0' COMMENT '执行时间',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_report_su_id` (`report_su_id`) USING BTREE,
  KEY `ix_job_name` (`job_name`) USING BTREE,
  KEY `ix_ctime` (`ctime`) USING BTREE,
  KEY `ix_mtime` (`mtime`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客户端监控表';

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '评论编号',
  `report_id` int(11) NOT NULL DEFAULT '0' COMMENT '压测报告id',
  `content` varchar(100) NOT NULL DEFAULT '' COMMENT '评论内容',
  `user_name` varchar(500) NOT NULL DEFAULT '' COMMENT '用户名',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '评论状态 1 正常 2 已删除',
  `submit_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论提交时间',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_report_id` (`report_id`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='评论表';

-- ----------------------------
-- Table structure for draft
-- ----------------------------
DROP TABLE IF EXISTS `draft`;
CREATE TABLE `draft` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '草稿箱id',
  `scene_id` int(11) NOT NULL COMMENT '场景id',
  `user_name` varchar(30) NOT NULL DEFAULT '' COMMENT '用户名',
  `is_active` tinyint(4) NOT NULL COMMENT '是否有效 0 无效 1 有效',
  `ctime` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `mtime` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Table structure for grpc
-- ----------------------------
DROP TABLE IF EXISTS `grpc`;
CREATE TABLE `grpc` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `department` varchar(50) NOT NULL DEFAULT '' COMMENT '部门',
  `project` varchar(50) NOT NULL DEFAULT '' COMMENT '项目',
  `app` varchar(50) NOT NULL DEFAULT '' COMMENT '应用',
  `threads_sum` int(11) NOT NULL DEFAULT '1' COMMENT '线程数',
  `ramp_up` int(11) NOT NULL DEFAULT '5' COMMENT '预热时间',
  `loops` int(11) NOT NULL DEFAULT '-1' COMMENT '循环次数：-1:永久',
  `load_time` int(11) NOT NULL DEFAULT '0' COMMENT '运行时长',
  `host_name` varchar(50) NOT NULL DEFAULT '' COMMENT '域名|IP',
  `port` int(6) NOT NULL DEFAULT '9000' COMMENT '端口',
  `service_name` varchar(50) NOT NULL DEFAULT '' COMMENT '服务名称',
  `proto_class_name` varchar(50) NOT NULL DEFAULT '' COMMENT 'proto类名称',
  `pkg_path` varchar(50) NOT NULL DEFAULT '' COMMENT '包名称',
  `asyn_call` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '-1:false 1:true, 0:--',
  `request_type` varchar(50) NOT NULL DEFAULT '' COMMENT '请求函数',
  `request_method` varchar(50) NOT NULL DEFAULT '' COMMENT 'grpc方法',
  `request_content` varchar(500) NOT NULL DEFAULT '' COMMENT 'grpc请求内容',
  `response_type` varchar(50) NOT NULL DEFAULT '' COMMENT '返回函数',
  `script_path` varchar(200) NOT NULL DEFAULT '' COMMENT 'proto文件路径',
  `jar_path` varchar(255) NOT NULL COMMENT 'jar文件路径',
  `jmx_path` varchar(200) NOT NULL DEFAULT '' COMMENT '生成jmx文件路径',
  `jmx_log` varchar(200) NOT NULL DEFAULT '' COMMENT 'jmx执行log',
  `jtl_log` varchar(200) NOT NULL DEFAULT '' COMMENT 'jtl log',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `active` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0 失效 1 生效',
  `update_by` varchar(20) NOT NULL DEFAULT '' COMMENT '更新人',
  `is_async` varchar(4) NOT NULL DEFAULT '' COMMENT '是否异步',
  `param_file_path` varchar(200) NOT NULL DEFAULT '' COMMENT '参数文件路径',
  `param_names` varchar(100) NOT NULL DEFAULT '' COMMENT '参数名称,以逗号分隔',
  `param_delimiter` varchar(5) NOT NULL DEFAULT '' COMMENT '参数分隔符,默认,',
  `param_enable` varchar(16) NOT NULL DEFAULT '' COMMENT '是否可用',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `ix_app` (`app`) USING BTREE,
  KEY `ix_department_project_app` (`app`,`department`,`project`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='grpc脚本表';

-- ----------------------------
-- Table structure for grpc_snap
-- ----------------------------
DROP TABLE IF EXISTS `grpc_snap`;
CREATE TABLE `grpc_snap` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `grpc_id` int(11) NOT NULL COMMENT '复制的脚本id',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `department` varchar(50) NOT NULL DEFAULT '' COMMENT '部门',
  `project` varchar(50) NOT NULL DEFAULT '' COMMENT '项目',
  `app` varchar(50) NOT NULL DEFAULT '' COMMENT '应用',
  `threads_sum` int(11) NOT NULL DEFAULT '1' COMMENT '线程数',
  `ramp_up` int(11) NOT NULL DEFAULT '5' COMMENT '预热时间',
  `loops` int(11) NOT NULL DEFAULT '-1' COMMENT '循环次数：-1:永久',
  `load_time` int(11) NOT NULL DEFAULT '0' COMMENT '运行时长',
  `host_name` varchar(50) NOT NULL DEFAULT '' COMMENT '域名|IP',
  `port` int(6) NOT NULL DEFAULT '9000' COMMENT '端口',
  `service_name` varchar(50) NOT NULL DEFAULT '' COMMENT '服务名称',
  `proto_class_name` varchar(50) NOT NULL DEFAULT '' COMMENT 'proto类名称',
  `pkg_path` varchar(50) NOT NULL DEFAULT '' COMMENT '包名称',
  `asyn_call` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '-1:false 1:true, 0:--',
  `request_type` varchar(50) NOT NULL DEFAULT '' COMMENT '请求函数',
  `request_method` varchar(50) NOT NULL DEFAULT '' COMMENT 'grpc方法',
  `request_content` varchar(500) NOT NULL DEFAULT '' COMMENT 'grpc请求内容',
  `response_type` varchar(50) NOT NULL DEFAULT '' COMMENT '返回函数',
  `script_path` varchar(200) NOT NULL DEFAULT '' COMMENT 'proto文件路径',
  `jar_path` varchar(255) NOT NULL COMMENT 'jar文件路径',
  `jmx_path` varchar(200) NOT NULL DEFAULT '' COMMENT '生成jmx文件路径',
  `jmx_log` varchar(200) NOT NULL DEFAULT '' COMMENT 'jmx执行log',
  `jtl_log` varchar(200) NOT NULL DEFAULT '' COMMENT 'jtl log',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `active` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0 失效 1 生效',
  `update_by` varchar(20) NOT NULL DEFAULT '' COMMENT '更新人',
  `execute_id` varchar(20) NOT NULL COMMENT '执行id',
  `is_async` varchar(4) NOT NULL DEFAULT '' COMMENT '是否异步',
  `param_file_path` varchar(200) NOT NULL DEFAULT '' COMMENT '参数文件路径',
  `param_names` varchar(100) NOT NULL DEFAULT '' COMMENT '参数名称,以逗号分隔',
  `param_delimiter` varchar(5) NOT NULL DEFAULT '' COMMENT '参数分隔符,默认,',
  `param_enable` varchar(16) NOT NULL DEFAULT '' COMMENT '是否支持参数化',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `ix_app` (`app`) USING BTREE,
  KEY `ix_department_project_app` (`app`,`department`,`project`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='grpc脚本表';

-- ----------------------------
-- Table structure for label
-- ----------------------------
DROP TABLE IF EXISTS `label`;
CREATE TABLE `label` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '标签名称',
  `description` varchar(100) NOT NULL DEFAULT '' COMMENT '描述',
  `color` varchar(100) NOT NULL DEFAULT '' COMMENT '标签颜色',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `active` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0 失效 1 生效',
  PRIMARY KEY (`id`),
  KEY `ix_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='标签表';

-- ----------------------------
-- Table structure for label_relation
-- ----------------------------
DROP TABLE IF EXISTS `label_relation`;
CREATE TABLE `label_relation` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `label_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '标签ID label.id',
  `label_name` varchar(50) NOT NULL DEFAULT '' COMMENT '标签名称',
  `target_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '目标ID',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 默认 1 脚本 2报告',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `active` tinyint(4) NOT NULL DEFAULT '1' COMMENT '0 失效 1 生效',
  PRIMARY KEY (`id`),
  KEY `ix_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='标签关系表';

-- ----------------------------
-- Table structure for order
-- ----------------------------
DROP TABLE IF EXISTS `order`;
CREATE TABLE `order` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '工单编号',
  `name` varchar(500) NOT NULL DEFAULT '' COMMENT '工单名称',
  `broker` varchar(100) NOT NULL DEFAULT '' COMMENT '研发对接人',
  `test_background` varchar(500) NOT NULL DEFAULT '' COMMENT '测试背景',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0: 开发自测， 1:工程效能团队测试',
  `test_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1压力测试 2负载测试 3 容量测试 4 健壮性测试 5 恢复性测试 6 浪涌测试 7配置选型测试 8 稳定性测试 9 特殊业务场景测试',
  `test_target` text NOT NULL COMMENT '测试指标',
  `api_list` text NOT NULL COMMENT '接口列 隔开',
  `api_doc` varchar(500) NOT NULL COMMENT '接口文档',
  `limit_user` varchar(500) NOT NULL DEFAULT '' COMMENT '用户限制',
  `limit_ip` varchar(100) NOT NULL DEFAULT '' COMMENT 'ip限制',
  `limit_visit` varchar(100) NOT NULL DEFAULT '' COMMENT '访问次数限制',
  `server_conf` varchar(100) NOT NULL DEFAULT '' COMMENT '服务器配置',
  `dependent_component` varchar(500) NOT NULL DEFAULT '' COMMENT '依赖组件',
  `dependent_business` varchar(500) NOT NULL DEFAULT '' COMMENT '依赖业务方',
  `test_data_from` varchar(500) NOT NULL DEFAULT '' COMMENT '测试数据获取',
  `test_host` varchar(100) NOT NULL DEFAULT '' COMMENT '测试机器地址',
  `moni_redis` varchar(200) NOT NULL COMMENT 'redis moni address',
  `moni_memcache` varchar(200) NOT NULL COMMENT 'memcache moni address',
  `moni_docker` varchar(200) NOT NULL COMMENT 'docker moni address',
  `moni_api` varchar(200) NOT NULL COMMENT 'api moni address',
  `moni_mysql` varchar(200) NOT NULL COMMENT 'mysql moni address',
  `moni_elasticsearch` varchar(200) NOT NULL COMMENT 'elasticsearch moni address',
  `moni_other` varchar(200) NOT NULL COMMENT 'other moni address',
  `test_cycles` int(11) NOT NULL COMMENT '测试周期',
  `script_id` varchar(500) NOT NULL DEFAULT '',
  `machine_id` varchar(200) NOT NULL DEFAULT '' COMMENT '机器编号',
  `department` varchar(20) NOT NULL COMMENT '部门',
  `project` varchar(20) NOT NULL COMMENT '项目',
  `app` varchar(20) NOT NULL COMMENT '应用',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '工单状态 0 申请中 -1 打回 1 排期中 2 进行中 3 测试完成',
  `update_by` int(11) NOT NULL DEFAULT '0' COMMENT '更新者',
  `apply_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请日期',
  `handler` varchar(50) NOT NULL DEFAULT '' COMMENT '处理人',
  `active` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '状态 -1 无效 1 生效',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_name` (`name`(255)),
  KEY `ix_active` (`active`),
  KEY `ix_apply_date` (`apply_date`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='工单表';

-- ----------------------------
-- Table structure for order_admin
-- ----------------------------
DROP TABLE IF EXISTS `order_admin`;
CREATE TABLE `order_admin` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_name` varchar(16) NOT NULL DEFAULT '' COMMENT '用户姓名',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='工单管理员';

-- ----------------------------
-- Table structure for project
-- ----------------------------
DROP TABLE IF EXISTS `project`;
CREATE TABLE `project` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '项目id',
  `name` varchar(120) NOT NULL COMMENT '项目名称',
  `update_by` bigint(20) NOT NULL COMMENT '修改人员id',
  `create_time` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '创建时间',
  `manager` varchar(100) DEFAULT '' COMMENT '项目管理员',
  `active` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '-1:失效；1:生效',
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='项目表';

-- ----------------------------
-- Table structure for ptest_job
-- ----------------------------
DROP TABLE IF EXISTS `ptest_job`;
CREATE TABLE `ptest_job` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `script_id` int(11) NOT NULL DEFAULT '0' COMMENT '脚本id',
  `report_su_id` int(11) NOT NULL DEFAULT '0' COMMENT '报告ID',
  `job_name` varchar(20) NOT NULL DEFAULT '' COMMENT 'job 名',
  `active` int(11) NOT NULL DEFAULT '1' COMMENT '是否有效,1 有效,-1 无效',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `execute_id` varchar(50) NOT NULL COMMENT '执行id',
  `host_ip` varchar(50) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `ix_report_su_id` (`report_su_id`) USING BTREE,
  KEY `ix_script_id` (`script_id`) USING BTREE,
  KEY `ix_mtime` (`mtime`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='报告脚本容器关联表';

-- ----------------------------
-- Table structure for report_graph
-- ----------------------------
DROP TABLE IF EXISTS `report_graph`;
CREATE TABLE `report_graph` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `test_name` varchar(50) NOT NULL DEFAULT '' COMMENT '接口名',
  `test_name_nick` varchar(50) NOT NULL DEFAULT '' COMMENT '接口别名',
  `count` bigint(20) NOT NULL DEFAULT '0' COMMENT '总请求数',
  `qps` int(11) NOT NULL DEFAULT '0' COMMENT 'qps',
  `avg_time` int(11) NOT NULL DEFAULT '0' COMMENT '平均时间',
  `min` int(11) NOT NULL DEFAULT '0' COMMENT '最小时间',
  `max` int(11) NOT NULL DEFAULT '0' COMMENT '最大时间',
  `error` int(11) NOT NULL DEFAULT '0' COMMENT '错误数',
  `fail_percent` varchar(11) NOT NULL DEFAULT '' COMMENT '失败率',
  `ninety_time` int(11) NOT NULL DEFAULT '0' COMMENT '90 分位',
  `ninety_five_time` int(11) NOT NULL DEFAULT '0' COMMENT '95分位',
  `ninety_nine_time` int(11) NOT NULL DEFAULT '0' COMMENT '99分位',
  `net_io` int(11) NOT NULL DEFAULT '0' COMMENT '网络流量',
  `code_ell` int(11) NOT NULL DEFAULT '0' COMMENT 'code200',
  `code_wll` int(11) NOT NULL DEFAULT '0' COMMENT 'code500',
  `code_wly` int(11) NOT NULL DEFAULT '0' COMMENT 'code501',
  `code_wle` int(11) NOT NULL DEFAULT '0' COMMENT 'code502',
  `code_wls` int(11) NOT NULL DEFAULT '0' COMMENT 'code504',
  `code_sll` int(11) NOT NULL DEFAULT '0' COMMENT 'code400',
  `code_sly` int(11) NOT NULL DEFAULT '0' COMMENT 'code401',
  `code_sls` int(11) NOT NULL DEFAULT '0' COMMENT 'code404',
  `code_kong` int(11) NOT NULL DEFAULT '0' COMMENT 'code_kong',
  `code_non_http` int(11) NOT NULL DEFAULT '0' COMMENT 'code_non_http',
  `code_others` int(11) NOT NULL DEFAULT '0' COMMENT 'code_others',
  `pod_name` varchar(25) NOT NULL DEFAULT '' COMMENT '容器全名',
  `threads_sum` int(11) NOT NULL DEFAULT '0' COMMENT '实时线程数',
  `elapsd_time` int(11) NOT NULL DEFAULT '0' COMMENT '持续时间',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `fifty_time` int(11) NOT NULL DEFAULT '0' COMMENT '50分位',
  `code301` int(11) NOT NULL DEFAULT '0' COMMENT 'code301',
  `code302` int(11) NOT NULL DEFAULT '0' COMMENT 'code302',
  PRIMARY KEY (`id`),
  KEY `ix_test_name_nick` (`test_name_nick`) USING BTREE,
  KEY `ix_test_name` (`test_name`) USING BTREE,
  KEY `ix_mtime` (`mtime`) USING BTREE,
  KEY `ix_pod_name` (`pod_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='报告图表';

-- ----------------------------
-- Table structure for report_info
-- ----------------------------
DROP TABLE IF EXISTS `report_info`;
CREATE TABLE `report_info` (
  `id` int(20) NOT NULL AUTO_INCREMENT COMMENT '报告id',
  `job_name` varchar(100) DEFAULT '',
  `project_name` varchar(100) DEFAULT '',
  `test_name` varchar(200) DEFAULT '' COMMENT '接口名',
  `request_count` varchar(200) DEFAULT '' COMMENT '总请求数',
  `avg_time` varchar(200) DEFAULT '' COMMENT '平均响应时间',
  `mid_time` varchar(200) DEFAULT '' COMMENT '中分位',
  `ninety_time` varchar(200) DEFAULT '' COMMENT '90分位',
  `ninety_five_time` varchar(200) DEFAULT '' COMMENT '95分位',
  `ninety_nine_time` varchar(200) DEFAULT '' COMMENT '收件人',
  `min` varchar(200) DEFAULT '',
  `max` varchar(200) DEFAULT '',
  `fail_percent` varchar(200) DEFAULT '' COMMENT '失败率',
  `qps` varchar(200) DEFAULT '' COMMENT 'qps',
  `net_io` varchar(200) DEFAULT '' COMMENT '网络流量',
  `ctime` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `src_name` varchar(100) DEFAULT '' COMMENT '测试报告源文件',
  `update_by` varchar(20) DEFAULT '' COMMENT '创建人',
  `final` int(4) DEFAULT '0' COMMENT '0 中间报告，1 最终报告',
  `active` tinyint(4) DEFAULT '1' COMMENT '状态：0 无效；1 生效',
  `mtime` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `job_name` (`job_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='报告表';

-- ----------------------------
-- Table structure for report_summary
-- ----------------------------
DROP TABLE IF EXISTS `report_summary`;
CREATE TABLE `report_summary` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `script_id` int(11) NOT NULL,
  `script_snap_id` int(11) NOT NULL,
  `execute_id` varchar(50) NOT NULL,
  `department` varchar(50) NOT NULL,
  `project` varchar(50) NOT NULL,
  `app` varchar(50) NOT NULL,
  `test_name` varchar(50) NOT NULL DEFAULT '' COMMENT '接口名',
  `test_name_nick` varchar(2000) NOT NULL DEFAULT '' COMMENT '接口别名',
  `job_name` varchar(20) NOT NULL DEFAULT '' COMMENT '容器名',
  `count` bigint(20) NOT NULL DEFAULT '0' COMMENT '总请求数',
  `qps` int(11) NOT NULL DEFAULT '0' COMMENT 'qps',
  `avg_time` int(11) NOT NULL DEFAULT '0' COMMENT '平均时间',
  `min` int(11) NOT NULL DEFAULT '0' COMMENT '最小时间',
  `max` int(11) NOT NULL DEFAULT '0' COMMENT '最大时间',
  `error` int(11) NOT NULL DEFAULT '0' COMMENT '错误数',
  `fail_percent` varchar(11) NOT NULL DEFAULT '' COMMENT '失败率',
  `ninety_time` int(11) NOT NULL DEFAULT '0' COMMENT '90 分位',
  `ninety_five_time` int(11) NOT NULL DEFAULT '0' COMMENT '95分位',
  `ninety_nine_time` int(11) NOT NULL DEFAULT '0' COMMENT '99分位',
  `net_io` int(11) NOT NULL DEFAULT '0' COMMENT '网络流量',
  `elapsd_time` int(11) NOT NULL DEFAULT '0' COMMENT '持续时间',
  `test_status` int(11) NOT NULL DEFAULT '2' COMMENT '1 :完成, 2 :执行中，3 中断',
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '执行人',
  `res_jtl` varchar(500) NOT NULL,
  `jmeter_log` varchar(500) DEFAULT NULL,
  `docker_sum` int(11) NOT NULL DEFAULT '0' COMMENT '容器数',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `debug` int(4) NOT NULL,
  `active` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否有效  1 有效',
  `scene_id` int(11) NOT NULL DEFAULT '0' COMMENT '场景id',
  `type` tinyint(4) DEFAULT '0',
  `load_time` int(11) NOT NULL DEFAULT '0' COMMENT '执行时间',
  `fifty_time` int(11) NOT NULL DEFAULT '0' COMMENT '50分位',
  PRIMARY KEY (`id`),
  KEY `ix_test_name` (`test_name`) USING BTREE,
  KEY `ix_mtime` (`mtime`) USING BTREE,
  KEY `ix_excute_id` (`execute_id`) USING BTREE,
  KEY `ix_scene_id` (`scene_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='报告汇总表';

-- ----------------------------
-- Table structure for report_timely
-- ----------------------------
DROP TABLE IF EXISTS `report_timely`;
CREATE TABLE `report_timely` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `test_name` varchar(50) CHARACTER SET utf8mb4 DEFAULT '' COMMENT '接口名',
  `count` bigint(20) NOT NULL DEFAULT '0' COMMENT '总请求数',
  `qps` int(11) NOT NULL DEFAULT '0' COMMENT 'qps',
  `avg_time` int(11) NOT NULL DEFAULT '0' COMMENT '平均响应时间',
  `min` int(11) NOT NULL DEFAULT '0' COMMENT '最小时间',
  `max` int(11) NOT NULL DEFAULT '0' COMMENT '最大时间',
  `error` int(11) NOT NULL DEFAULT '0' COMMENT '错误数',
  `fail_percent` varchar(11) NOT NULL DEFAULT '' COMMENT '失败率',
  `ninety_time` int(11) NOT NULL DEFAULT '0' COMMENT '90分位',
  `ninety_five_time` int(11) NOT NULL DEFAULT '0' COMMENT '95分位',
  `ninety_nine_time` int(11) NOT NULL DEFAULT '0' COMMENT '99分位',
  `net_io` int(11) NOT NULL DEFAULT '0' COMMENT '网络流量',
  `code_ell` int(11) NOT NULL,
  `code_wll` int(11) NOT NULL,
  `code_wly` int(11) DEFAULT NULL,
  `code_wle` int(11) DEFAULT NULL,
  `code_wls` int(11) DEFAULT NULL,
  `code_sll` int(11) DEFAULT NULL,
  `code_sly` int(11) DEFAULT NULL,
  `code_sls` int(11) DEFAULT NULL,
  `code_kong` int(11) DEFAULT NULL,
  `code_non_http` int(11) DEFAULT NULL,
  `code_others` int(11) DEFAULT NULL,
  `pod_name` varchar(25) NOT NULL DEFAULT '' COMMENT '容器全名',
  `threads_sum` int(11) NOT NULL DEFAULT '0' COMMENT '实时线程数',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `fifty_time` int(11) NOT NULL DEFAULT '0' COMMENT '50分位',
  `code301` int(11) NOT NULL DEFAULT '0' COMMENT 'code301',
  `code302` int(11) NOT NULL DEFAULT '0' COMMENT 'code302',
  PRIMARY KEY (`id`),
  KEY `ix_test_name` (`test_name`) USING BTREE,
  KEY `ix_pod_name` (`pod_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for scene
-- ----------------------------
DROP TABLE IF EXISTS `scene`;
CREATE TABLE `scene` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '场景id',
  `scene_name` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '场景名称',
  `scene_type` tinyint(4) NOT NULL COMMENT '场景类型 1 自动分组 2 自定义分组 根据最后保存接口/接口组设置的页面类型来更新该字段的值',
  `user_name` varchar(30) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '用户名',
  `is_draft` varchar(4) NOT NULL COMMENT '是否为草稿 0非草稿 1草稿',
  `is_debug` varchar(4) NOT NULL COMMENT 'is_debug   是否调试 0 执行压测 1 调试',
  `jmeter_file_path` varchar(100) DEFAULT NULL,
  `department` varchar(20) DEFAULT NULL,
  `project` varchar(20) DEFAULT NULL,
  `app` varchar(20) DEFAULT NULL,
  `jmeter_log` varchar(100) DEFAULT NULL,
  `res_jtl` varchar(100) DEFAULT NULL,
  `ctime` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `mtime` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `is_active` varchar(4) NOT NULL COMMENT '草稿是否有效',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Table structure for script
-- ----------------------------
DROP TABLE IF EXISTS `script`;
CREATE TABLE `script` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `tree_id` bigint(20) DEFAULT NULL,
  `project_id` bigint(20) DEFAULT NULL COMMENT '脚本id',
  `type` int(2) DEFAULT NULL COMMENT '脚本类型，默认 1 为 jmeter',
  `project_name` varchar(100) DEFAULT NULL COMMENT '项目名称',
  `test_name` varchar(50) DEFAULT NULL COMMENT '接口名称',
  `threads_sum` int(6) DEFAULT NULL COMMENT '总线程数',
  `ready_time` int(6) DEFAULT NULL,
  `load_time` int(6) DEFAULT NULL COMMENT '压测持续时间',
  `proc_type` varchar(16) DEFAULT '' COMMENT '脚本协议类型',
  `url` varchar(500) DEFAULT '',
  `domain` varchar(50) DEFAULT '' COMMENT '被测试的域名',
  `port` varchar(16) NOT NULL DEFAULT '' COMMENT '端口',
  `login` varchar(16) NOT NULL DEFAULT '' COMMENT '是否登录',
  `path` varchar(500) DEFAULT NULL COMMENT '路径',
  `method` varchar(10) DEFAULT NULL COMMENT '方法，post 或者 get ',
  `content_type` varchar(50) CHARACTER SET latin1 DEFAULT '',
  `cookie` varchar(500) CHARACTER SET latin1 DEFAULT '',
  `data` varchar(1000) DEFAULT NULL COMMENT 'json body',
  `assertion` varchar(50) DEFAULT NULL COMMENT '断言',
  `update_by` varchar(50) DEFAULT NULL COMMENT '更新人',
  `save_path` varchar(200) DEFAULT NULL,
  `res_jtl` varchar(100) DEFAULT '',
  `jmeter_log` varchar(100) DEFAULT NULL,
  `ctime` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `active` tinyint(4) DEFAULT '1' COMMENT '状态 ，1 为有效  -1 为无效',
  `upload` varchar(16) DEFAULT '' COMMENT '是否上传',
  `mtime` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `department` varchar(50) NOT NULL DEFAULT '' COMMENT '部門',
  `project` varchar(50) NOT NULL DEFAULT '' COMMENT '项目',
  `app` varchar(50) NOT NULL DEFAULT '' COMMENT '应用',
  `api_header` varchar(500) NOT NULL DEFAULT '' COMMENT '请求头',
  `argument_map` varchar(500) NOT NULL DEFAULT '' COMMENT 'POST提交参数',
  `use_data_file` varchar(16) NOT NULL DEFAULT '' COMMENT '使用文件: 0 不使用 1 使用',
  `file_name` varchar(100) NOT NULL DEFAULT '' COMMENT '文件名称',
  `params_name` varchar(16) NOT NULL DEFAULT '' COMMENT '参数名称',
  `delimiter` varchar(16) NOT NULL DEFAULT '' COMMENT '文本切割符',
  `loops` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '脚本循环次数:-1 永久循环',
  `file_split` varchar(16) NOT NULL DEFAULT '' COMMENT '是否切割文件，0不切割 1切割',
  `split_num` tinyint(4) NOT NULL,
  `use_sign` varchar(2) NOT NULL DEFAULT '' COMMENT '是否需要签名',
  `conn_time_out` int(6) NOT NULL DEFAULT '0' COMMENT '连接超时时间',
  `resp_time_out` int(6) NOT NULL DEFAULT '0' COMMENT '响应超时时间',
  `test_type` tinyint(4) NOT NULL COMMENT '压测类型 0 http 1 grpc 2 场景',
  `scene_id` int(11) NOT NULL COMMENT '场景id 关联scene表中的自增长id',
  `output_params` varchar(255) NOT NULL DEFAULT '' COMMENT '接口输出参数，多个用英文,隔开',
  `group_id` int(11) NOT NULL,
  `run_order` int(11) NOT NULL,
  `script_path` varchar(200) NOT NULL DEFAULT '' COMMENT '脚本路径',
  `json_path` varchar(100) NOT NULL DEFAULT '' COMMENT 'JSON 解析参数路径',
  `is_async` varchar(4) NOT NULL DEFAULT '' COMMENT '是否异步',
  `multipart_path` varchar(100) NOT NULL DEFAULT '' COMMENT 'multipart 路径',
  `multipart_file` varchar(50) NOT NULL DEFAULT '' COMMENT 'multipart 文件名',
  `multipart_param` varchar(50) NOT NULL DEFAULT '' COMMENT 'multipart 参数',
  `mime_type` varchar(50) NOT NULL DEFAULT '' COMMENT 'mime_type 类型',
  `fusing` int(4) NOT NULL DEFAULT '0' COMMENT '自动熔断成功率',
  `keep_alive` varchar(4) NOT NULL DEFAULT '1' COMMENT '是否使用长连接',
  PRIMARY KEY (`id`),
  KEY `tree_id` (`tree_id`) USING BTREE,
  KEY `project_id` (`project_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for script_snap
-- ----------------------------
DROP TABLE IF EXISTS `script_snap`;
CREATE TABLE `script_snap` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `script_id` int(11) NOT NULL DEFAULT '0' COMMENT '脚本id',
  `tree_id` bigint(20) DEFAULT NULL,
  `project_id` bigint(20) DEFAULT NULL COMMENT '脚本id',
  `execute_id` varchar(50) DEFAULT NULL,
  `type` tinyint(2) DEFAULT NULL COMMENT '脚本类型，默认 1 为 jmeter',
  `project_name` varchar(100) DEFAULT NULL COMMENT '项目名称',
  `test_name` varchar(50) DEFAULT NULL COMMENT '接口名称',
  `threads_sum` int(6) DEFAULT NULL COMMENT '总线程数',
  `ready_time` int(6) DEFAULT NULL,
  `load_time` int(6) DEFAULT NULL COMMENT '压测持续时间',
  `proc_type` varchar(16) DEFAULT '' COMMENT '脚本协议类型',
  `url` varchar(500) DEFAULT '',
  `domain` varchar(50) DEFAULT '' COMMENT '被测试的域名',
  `port` varchar(16) NOT NULL DEFAULT '' COMMENT '端口',
  `login` varchar(16) NOT NULL DEFAULT '' COMMENT '是否登录',
  `path` varchar(500) DEFAULT NULL COMMENT '路径',
  `method` varchar(10) DEFAULT NULL COMMENT '方法，post 或者 get ',
  `content_type` varchar(50) CHARACTER SET latin1 DEFAULT '',
  `cookie` varchar(500) CHARACTER SET latin1 DEFAULT '',
  `data` varchar(1000) DEFAULT NULL COMMENT 'json body',
  `assertion` varchar(50) DEFAULT NULL COMMENT '断言',
  `update_by` varchar(50) DEFAULT NULL COMMENT '更新人',
  `save_path` varchar(200) DEFAULT NULL,
  `res_jtl` varchar(100) DEFAULT '',
  `jmeter_log` varchar(100) DEFAULT NULL,
  `ctime` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `active` tinyint(4) DEFAULT '1' COMMENT '状态 ，1 为有效  -1 为无效',
  `upload` varchar(16) DEFAULT '' COMMENT '是否上传',
  `mtime` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `department` varchar(50) NOT NULL DEFAULT '' COMMENT '部門',
  `project` varchar(50) NOT NULL DEFAULT '' COMMENT '项目',
  `app` varchar(50) NOT NULL DEFAULT '' COMMENT '应用',
  `api_header` varchar(500) NOT NULL DEFAULT '' COMMENT '请求头',
  `argument_map` varchar(500) NOT NULL DEFAULT '' COMMENT 'POST提交参数',
  `use_data_file` varchar(16) NOT NULL DEFAULT '' COMMENT '使用文件: 0 不使用 1 使用',
  `file_name` varchar(200) NOT NULL DEFAULT '' COMMENT '文件名',
  `params_name` varchar(16) NOT NULL DEFAULT '' COMMENT '参数名称',
  `delimiter` varchar(16) NOT NULL DEFAULT '' COMMENT '文本切割符',
  `loops` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '脚本循环次数:-1 永久循环',
  `file_split` varchar(16) NOT NULL,
  `split_num` tinyint(4) NOT NULL,
  `use_sign` varchar(16) NOT NULL,
  `conn_time_out` int(6) NOT NULL DEFAULT '0' COMMENT '连接超时时间',
  `scene_id` int(11) NOT NULL,
  `resp_time_out` int(6) NOT NULL DEFAULT '0' COMMENT '响应超时时间',
  `json_path` varchar(100) NOT NULL DEFAULT '' COMMENT 'JSON 解析参数路径',
  `group_id` int(11) NOT NULL DEFAULT '0' COMMENT '分组id',
  `is_async` varchar(4) NOT NULL DEFAULT '' COMMENT '是否异步',
  `multipart_path` varchar(100) NOT NULL DEFAULT '' COMMENT 'multipart 路径',
  `multipart_file` varchar(50) NOT NULL DEFAULT '' COMMENT 'multipart 文件名',
  `multipart_param` varchar(50) NOT NULL DEFAULT '' COMMENT 'multipart 参数',
  `mime_type` varchar(50) NOT NULL DEFAULT '' COMMENT 'mime_type 类型',
  `fusing` int(4) NOT NULL DEFAULT '0' COMMENT '自动熔断成功率',
  `keep_alive` varchar(4) NOT NULL DEFAULT '1' COMMENT '是否使用长连接',
  PRIMARY KEY (`id`),
  KEY `tree_id` (`tree_id`) USING BTREE,
  KEY `project_id` (`project_id`) USING BTREE,
  KEY `excute_id` (`execute_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='脚本快照表';

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '用户名字',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '用户邮箱',
  `active` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '是否有效:-1 无效，1 有效',
  `accept` tinyint(4) NOT NULL DEFAULT '-1' COMMENT '-1 不允许访问 1 允许访问',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

-- ----------------------------
-- Table structure for work_order
-- ----------------------------
DROP TABLE IF EXISTS `work_order`;
CREATE TABLE `work_order` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '工单id',
  `name` varchar(1000) NOT NULL COMMENT '工单名称',
  `content` mediumtext COMMENT '工单正文',
  `type` tinyint(4) DEFAULT '0' COMMENT '0: 开发自测， 1:EP测试',
  `script_id` bigint(20) DEFAULT '0' COMMENT '脚本id，默认0',
  `machine_id` bigint(20) DEFAULT '0' COMMENT '机器id，默认0',
  `project_id` bigint(20) NOT NULL COMMENT '项目id',
  `status` tinyint(4) DEFAULT '0' COMMENT '工单状态:0：申请中，-1：打回，1：排期中，2：进行中，3、测试完成',
  `update_by` bigint(20) NOT NULL COMMENT '更新者',
  `apply_date` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '申请日期',
  `active` tinyint(4) DEFAULT '-1' COMMENT '状态：-1 无效；1 生效',
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`(255)),
  KEY `idx_machine_id` (`machine_id`),
  KEY `idx_project_id` (`project_id`),
  KEY `idx_active` (`active`),
  KEY `idx_apply_date` (`apply_date`)
) ENGINE=InnoDB CHARSET=utf8 COMMENT='工单表';

SET FOREIGN_KEY_CHECKS = 1;
