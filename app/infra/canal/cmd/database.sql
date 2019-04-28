create database bilibili_canal;

 CREATE TABLE `master_info` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `addr` varchar(64) NOT NULL COMMENT 'db addr hostname:port',
  `bin_name` varchar(20) NOT NULL DEFAULT '' COMMENT 'binlog name',
  `bin_pos` int(11) NOT NULL DEFAULT '0' COMMENT 'binlog position',
  `remark` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `cluster` varchar(50) NOT NULL DEFAULT '' COMMENT 'cluster',
  `leader` varchar(20) NOT NULL DEFAULT '' COMMENT 'leader',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除 0-否 1-是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_addr` (`addr`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=46585559 DEFAULT CHARSET=utf8 COMMENT='canal位置信息记录'

CREATE TABLE `canal_apply` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `addr` varchar(64) NOT NULL COMMENT 'db addr hostname:port',
  `remark` varchar(100) NOT NULL DEFAULT '' COMMENT 'remark',
  `cluster` varchar(50) NOT NULL DEFAULT '' COMMENT '集群',
  `leader` varchar(20) NOT NULL DEFAULT '' COMMENT 'leader',
  `comment` text NOT NULL COMMENT 'comment',
  `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT 'state',
  `operator` varchar(32) NOT NULL DEFAULT '' COMMENT 'operator',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `conf_id` int(11) NOT NULL DEFAULT '0' COMMENT '配置id',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_addr` (`addr`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8 COMMENT='canal申请信息'

CREATE TABLE `hbase_info` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `cluster_name` varchar(20) NOT NULL DEFAULT '' COMMENT '集群名称',
  `table_name` varchar(60) NOT NULL DEFAULT '' COMMENT '表名',
  `latest_ts` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'lastest ts',
  `remark` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name_table` (`cluster_name`,`table_name`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='hbase latest_ts表'

CREATE TABLE `tidb_info` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(20) NOT NULL DEFAULT '' COMMENT 'name',
  `cluster_id` varchar(40) NOT NULL DEFAULT '' COMMENT 'cluster id',
  `offset` bigint(20) NOT NULL DEFAULT 0 COMMENT 'offset',
  `tso` bigint(20) NOT NULL DEFAULT '0' COMMENT '全局时间戳',
  `remark` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `ix_mtime` (`mtime`)
) COMMENT='tidb info';