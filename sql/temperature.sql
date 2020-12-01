CREATE TABLE `t_device` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `sn` varchar(50) NOT NULL DEFAULT '' COMMENT '序列号',
  `mac` varchar(50) NOT NULL DEFAULT '' COMMENT 'mac',
  `store_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户id',
  `province` int(11) NOT NULL DEFAULT '0' COMMENT '省',
  `city` int(11) NOT NULL DEFAULT '0' COMMENT '市',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '设备状态',
  `remake` varchar(255) NOT NULL DEFAULT '' COMMENT '备注信息',
  `created_at` int(255) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `sn` (`sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='设备表';

CREATE TABLE `t_face` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `baidu_user_id` varchar(128) NOT NULL DEFAULT '' COMMENT '百度用户ID',
  `user_info` varchar(256) NOT NULL DEFAULT '' COMMENT '用户信息',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `face_sn` (`baidu_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='百度face信息';

CREATE TABLE `t_image` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `store_id` int(11) NOT NULL DEFAULT '0',
  `sn` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '设备序列号',
  `url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '图片地址',
  `created_at` int(11) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='照片';

CREATE TABLE `t_record` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `face_id` int(11) NOT NULL DEFAULT '0' COMMENT 'face ID',
  `num` float(3,1) NOT NULL DEFAULT '0.0' COMMENT '体温',
  `sn` varchar(50) NOT NULL COMMENT '设备序列号',
  `store_id` int(11) NOT NULL COMMENT '店铺ID',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '1为注册用户，0为未注册用户',
  `face_image` varchar(500) NOT NULL COMMENT '图片地址',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `sn` (`sn`),
  KEY `face_id` (`face_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='测量记录';

CREATE TABLE `t_store` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '公司名称',
  `address` varchar(255) NOT NULL DEFAULT '' COMMENT '公司地址',
  `tel` varchar(30) NOT NULL DEFAULT '' COMMENT '联系电诱',
  `role` tinyint(1) NOT NULL DEFAULT '1' COMMENT '角色（0：超级管理员，1：普通用户）',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '1',
  `created_at` int(11) NOT NULL DEFAULT '0',
  `updated_at` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='帐号表';

CREATE TABLE `t_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '姓名',
  `user_no` varchar(50) NOT NULL DEFAULT '' COMMENT '用户编号',
  `face_id` int(11) NOT NULL DEFAULT '0' COMMENT 'face id',
  `tel` varchar(15) NOT NULL DEFAULT '' COMMENT '电话号码',
  `id_num` varchar(20) NOT NULL DEFAULT '' COMMENT '身份证编号',
  `age` tinyint(3) NOT NULL DEFAULT '0' COMMENT '年龄',
  `company_id` int(11) NOT NULL DEFAULT '0' COMMENT '公司编号',
  `company` varchar(255) NOT NULL DEFAULT '' COMMENT '公司名',
  `department` varchar(255) DEFAULT '' COMMENT '部门',
  `face_image` varchar(255) NOT NULL DEFAULT '' COMMENT '百度注册人脸',
  `app_id` int(10) NOT NULL DEFAULT '0' COMMENT '百度应用ID',
  `group` varchar(50) NOT NULL DEFAULT '' COMMENT '百度人年库用户组',
  `created_at` int(11) NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `user_no` (`user_no`),
  KEY `tel` (`tel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户';