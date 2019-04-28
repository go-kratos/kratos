# (fat 1, uat 1, prod 0)
CREATE TABLE `sign_up` (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '签约ID',
  sex tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别，性别 0:保密 1:男 2:女',
  mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  begin_date date NOT NULL DEFAULT '0000-00-00' COMMENT '签约开始时间',
  end_date date NOT NULL DEFAULT '0000-00-00' COMMENT '签约结束时间',
  state tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常，100删除',
  country varchar(16) NOT NULL DEFAULT '' COMMENT '国家',
  province varchar(16) NOT NULL DEFAULT '' COMMENT '省',
  city varchar(16) NOT NULL DEFAULT '' COMMENT '市',
  note varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (id),
  KEY ix_mid (mid),
  KEY ix_begin_date (begin_date),
  KEY ix_end_date (end_date),
  KEY ix_mtime (mtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='签约表';

-- 付款表
CREATE TABLE sign_pay (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT '签约ID',
  due_date date NOT NULL DEFAULT '0000-00-00' COMMENT '签约结束时间',
  pay_value BIGINT(20) NOT NULL DEFAULT '0' COMMENT '金额',
  state tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0未支付，1已支付，100删除',
  note varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (id),
  KEY ix_signid (sign_id),
  KEY ix_mid (mid),
  KEY ix_date (due_date),
  KEY ix_mtime (mtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='付款表';

CREATE TABLE sign_task (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '任务ID',
  mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT '签约ID',
  task_type tinyint(4) NOT NULL DEFAULT '0' COMMENT '任务类型',
  task_counter int(11) NOT NULL DEFAULT '0' COMMENT '任务计数器',
  task_condition int(11) NOT NULL DEFAULT '0' COMMENT '任务条件',
  task_data varchar(1024) NOT NULL DEFAULT '' COMMENT '任务存储相关数据', -- 任务数据，比如用来存已经投过的稿件id
  state tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常，1过期，100删除',
  generate_date date NOT NULL DEFAULT '0000-00-00' COMMENT '任务生成时间',
  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',

  PRIMARY KEY (id),
  KEY ix_mid (mid),
  KEY ix_signid (sign_id),
  KEY ix_date (generate_date),
  KEY ix_mtime (mtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务表';

-- 合同表
CREATE TABLE sign_contract (
  id int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '合同ID',
  sign_id int(11) unsigned NOT NULL DEFAULT '0' COMMENT '签约ID',
  mid int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'up主id',
  filename varchar(255) NOT NULL DEFAULT '' COMMENT '合同名',
  filelink varchar(255) NOT NULL DEFAULT '' COMMENT '文件链接地址',
  state tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态，0正常，100删除',

  ctime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  mtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (id),
  KEY ix_mid (mid),
  KEY ix_signid (sign_id),
  KEY ix_mtime (mtime)
)  ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='合同表';


# (fat 1, uat 1, prod 1)
alter table sign_up add column admin_id int(11) not null default 0 comment '管理员id';
alter table sign_up add column admin_name varchar(32) not null default '' comment '管理员name';

# (fat 1, uat 1, prod 1) --2018.06.26
alter table sign_up add column email_state tinyint(4) not null default 0 comment '邮件发送情况，0未发送，1已发送过提醒邮件';
alter table sign_pay add column email_state tinyint(4) not null default 0 comment '邮件发送情况，0未发送，1已发送过提醒邮件';
