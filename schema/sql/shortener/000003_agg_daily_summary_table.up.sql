-- 每日统计聚合表
CREATE TABLE IF NOT EXISTS `agg_daily_summary` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `short_url_id` bigint(20) unsigned NOT NULL,
    `summary_date` date NOT NULL COMMENT '统计日期',
    `total_clicks` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '总点击量',
    `unique_visitors` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '独立访客数 (UV)',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_url_id_date` (
        `short_url_id`,
        `summary_date`
    )
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '每日统计聚合表';