CREATE TABLE
  `activity_score` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `activity_id` int(11) NOT NULL,
    `activity_type` tinyint(3) NOT NULL,
    `activity_start_ts` bigint NOT NULL DEFAULT '0',
    `activity_title` varchar(255) DEFAULT NULL,
    `total_cost` decimal(10, 2) NOT NULL DEFAULT '0.00' COMMENT '总成本',
    `activity_duration` int(11) NOT NULL DEFAULT '0' COMMENT '活动持续时长',
    `warn_time` int(11) NOT NULL DEFAULT '0' COMMENT '报警次数',
    `notch_number` int(11) NOT NULL DEFAULT '0' COMMENT '缺口数',
    `expected_product_amount` int(11) NOT NULL DEFAULT 0 COMMENT '预期流通份数',
    `recommend_score_before_end` decimal(10, 2) NOT NULL DEFAULT '0.00' COMMENT '活动结束前的推荐分',
    `praise_score` decimal(10, 2) NOT NULL DEFAULT '0.00' COMMENT '活动结束口碑分',
    `status` tinyint(2) NOT NULL COMMENT '状态:0未结束；1：已结束',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `uniq_idx_activity_id_type` (`activity_id`, `activity_type`),
    KEY `idx_deleted_at` (`deleted_at`)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci


ALTER TABLE `activity_score` CHANGE COLUMN `expected_product_amount` `expected_product_amount` int(11) NOT NULL DEFAULT 0 COMMENT '预期流通份数';

ALTER TABLE `activity_score` MODIFY COLUMN `activity_start_ts` bigint NOT NULL DEFAULT '0' COMMENT '活动开始时间';