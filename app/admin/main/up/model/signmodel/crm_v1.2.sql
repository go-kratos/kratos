
CREATE TABLE `sign_task_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  `mid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  `sign_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '签约ID',
  `task_template_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'sign_task模板表中的任务ID',
  `task_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型 0累积 1日 2周 3月 4季度',
  `task_counter` int(11) NOT NULL DEFAULT '0' COMMENT '任务计数器',
  `task_condition` int(11) NOT NULL DEFAULT '0' COMMENT '任务条件',
  `attribute` bigint(20) NOT NULL DEFAULT '0' COMMENT '属性位',
  `task_data` varchar(1024) NOT NULL DEFAULT '' COMMENT '任务存储相关数据',
  `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1、未完成，2、完成， 100、删除',
  `generate_date` date NOT NULL DEFAULT '0000-00-00' COMMENT '任务开始时间',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_template_id_generate_date` (`task_template_id`,`generate_date`),
  KEY `ix_mid` (`mid`),
  KEY `ix_sign_id` (`sign_id`),
  KEY `ix_generate_date` (`generate_date`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='任务历史数据表';

CREATE TABLE `sign_task_absence` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `sign_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '签约ID',
  `mid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  `task_history_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'sign_task_history表中ID',
  `absence_count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '请假数量',
  `reason` varchar(255) NOT NULL DEFAULT '' COMMENT '请假理由',
  `state` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常，100删除',
  `admin_id` int(11) NOT NULL DEFAULT '0' COMMENT '管理员id',
  `admin_name` varchar(32) NOT NULL DEFAULT '' COMMENT '管理员name',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_mid` (`mid`),
  KEY `ix_sign_id` (`sign_id`),
  KEY `ix_task_history_id` (`task_history_id`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 COMMENT='任务请假表';

CREATE TABLE `sign_violation_history` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `sign_id` int(11) NOT NULL DEFAULT '0' COMMENT '签约id',
  `mid` bigint(20) NOT NULL DEFAULT '0' COMMENT '违约人',
  `admin_id` int(11) NOT NULL DEFAULT '0' COMMENT '操作人id',
  `admin_name` varchar(32) NOT NULL DEFAULT '' COMMENT '操作人名字',
  `violation_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '违约原因',
  `state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '违约状态 1:违约 100:删除',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_sign_id` (`sign_id`),
  KEY `ix_mtime` (`mtime`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT '违约历史表';

#新增字段
alter table sign_task add column attribute bigint(20) not null default '0' comment '属性位';
alter table sign_task add column finish_note varchar(255) NOT NULL DEFAULT '' COMMENT '任务完成方式';

alter table sign_pay add column `in_tax` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否含税:1 不含税  2含税';
alter table sign_up add column `organization` tinyint(4) NOT NULL DEFAULT '1' COMMENT '组织属性: 1个人 2公司';
alter table sign_up add column `sign_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '签约类型: 0 其他 、1独家、2首发、3独家系列、4独家（双微除外）、5独家（微博除外）';
alter table sign_up add column `age` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '年龄';
alter table sign_up add column `residence` varchar(255) NOT NULL DEFAULT '' COMMENT '居住地';
alter table sign_up add column `id_card` varchar(20) NOT NULL DEFAULT '' COMMENT '身份证';
alter table sign_up add column `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '联系方式';
alter table sign_up add column `qq` bigint(20) unsigned NOT NULL COMMENT 'qq号';
alter table sign_up add column `wechat` varchar(16) NOT NULL DEFAULT '' COMMENT '微信号';
alter table sign_up add column `is_economic` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是非签署经济约 1否 2是';
alter table sign_up add column `economic_company` varchar(16) NOT NULL DEFAULT '' COMMENT '签约的经济公司';
alter table sign_up add column `task_state` tinyint(4) NOT NULL DEFAULT '1' COMMENT '任务完成度: 1 未完成 2 已完成';
alter table sign_up add column `wechat` varchar(16) NOT NULL DEFAULT '' COMMENT '微信号';
alter table sign_up add column `leave_times` int(11) NOT NULL COMMENT '请假次数';
alter table sign_up add column `violation_times` int(11) NOT NULL COMMENT '违约次数';
alter table sign_up add column `active_tid` smallint(6) unsigned NOT NULL DEFAULT '0' COMMENT 'up所属主分区';

#新增索引
ALTER TABLE sign_up ADD INDEX ix_active_tid(`active_tid`);
ALTER TABLE sign_up ADD INDEX ix_sex(`sex`);
ALTER TABLE sign_up ADD INDEX ix_country(`country`);
ALTER TABLE sign_up ADD INDEX ix_task_state(`task_state`);
ALTER TABLE sign_up ADD INDEX ix_sign_type(`sign_type`);

alter table sign_up add `economic_begin` date NOT NULL DEFAULT '0000-00-00' COMMENT '经济约的签约开始时间';
alter table sign_up add `economic_end` date NOT NULL DEFAULT '0000-00-00' COMMENT '经济约的签约结束时间';