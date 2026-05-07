drop table if exists `table_scan_record`;
CREATE TABLE `table_scan_record` (
         `id` int NOT NULL AUTO_INCREMENT,
         `table` varchar(64) NOT NULL DEFAULT '' COMMENT '表名',
         `last_id` bigint NOT NULL NULL DEFAULT 0 COMMENT '最后一个处理的ID',
         `task_type` int NOT NULL NULL DEFAULT 0 COMMENT '任务类型，1统计艺术家活动数',
         `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
         `updated_at` datetime DEFAULT NULL,
         `deleted_at` datetime DEFAULT NULL,
         PRIMARY KEY (`id`),
         UNIQUE KEY `uniq_table_task`  (`table`,`task_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据表扫描记录';