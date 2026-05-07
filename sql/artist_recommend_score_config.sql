CREATE TABLE `artist_recommend_score_config` (
                               `id` int NOT NULL AUTO_INCREMENT,
                               `artist_id` bigint NOT NULL NULL DEFAULT 0 COMMENT '艺术家ID',
                               `artist_name` varchar(128) NOT NULL DEFAULT '' COMMENT '艺术家名称',
                               `nft_num` int NOT NULL NULL DEFAULT 0 COMMENT '藏品数',
                               `activity_num` int NOT NULL NULL DEFAULT 0 COMMENT '活动数',
                               `score` varchar(128) NOT NULL DEFAULT '100' COMMENT '口碑分,初始值100',
                               `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
                               `updated_at` datetime DEFAULT NULL,
                               `deleted_at` datetime DEFAULT NULL,
                               PRIMARY KEY (`id`),
                               UNIQUE KEY `uniq_artist_id` (`artist_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT '艺术家推荐分设置';