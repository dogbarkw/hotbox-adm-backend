CREATE TABLE `hotbox_yop_test_user` (
                                    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
                                    `user_id` bigint NOT NULL DEFAULT '0' COMMENT '平台用户ID',
                                    `mobile` varchar(32) NOT NULL DEFAULT '' COMMENT '手机号',
                                    `real_name` varchar(32) NOT NULL DEFAULT '' COMMENT '账号实名',
                                    `user_type` tinyint NOT NULL DEFAULT '0' COMMENT '账号类型 1实名账号 2测试账号',
                                    `rate` int NOT NULL DEFAULT '0' COMMENT '分成比例',
                                    `total_income` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '累计进账',
                                    `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
                                    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
                                    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                    `count_time` datetime DEFAULT NULL COMMENT '上次进账统计时间',
                                    PRIMARY KEY (`id`),
                                    UNIQUE KEY `uniq_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='特殊账号表';

CREATE TABLE `hotbox_yop_test_user_rate_record` (
                                                `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
                                                `yop_test_user_id` bigint NOT NULL DEFAULT 0 COMMENT '特殊账号ID',
                                                `rate` int NOT NULL DEFAULT '0' COMMENT '分成比例',
                                                `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
                                                `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                                PRIMARY KEY (`id`),
                                                KEY `idx_id_created_at` (`yop_test_user_id`,`created_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='特殊账号分成记录表';

CREATE TABLE `hotbox_yop_test_user_income_record` (
                                                  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键 ID',
                                                  `yop_test_user_id` bigint NOT NULL DEFAULT 0 COMMENT '特殊账号ID',
                                                  `fee` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '金额',
                                                  `rate` int NOT NULL DEFAULT '0' COMMENT '当前分成比例',
                                                  `income` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '进账',
                                                  `income_time` datetime DEFAULT NULL COMMENT '进账时间',
                                                  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
                                                  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                                  PRIMARY KEY (`id`),
                                                  KEY `yop_test_user_id` (`yop_test_user_id`),
                                                  KEY `idx_id_created_at` (`yop_test_user_id`,`created_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='特殊账号进账记录表';