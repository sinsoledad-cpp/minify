-- 文件: app/shortener/schema/sql/000003_analytics_summary_daily.up.sql
-- 职责: 存储预聚合后的报表数据 (读密集型)，供 API 快速查询。

CREATE TABLE `analytics_summary_daily` (
                                           `id` BIGINT UNSIGNED AUTO_INCREMENT ,
                                           `link_id` BIGINT UNSIGNED NOT NULL COMMENT '关联的链接ID links(id)',
                                           `date` DATE NOT NULL COMMENT '聚合的日期 (e.g., 2025-10-21)',
                                           `dimension_type` VARCHAR(30) NOT NULL COMMENT '维度类型 (e.g., total, timeseries_hourly, referer, country, browser, os, device)',
                                           `dimension_value` VARCHAR(255) NOT NULL COMMENT '维度值 (e.g., 10 (for hourly), google.com, USA, Chrome)',
                                           `click_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '该维度在该天的总点击次数',

                                           PRIMARY KEY (`id`),
    -- 唯一键确保后台任务可重复运行 (INSERT ... ON DUPLICATE KEY UPDATE)
                                           UNIQUE KEY `uniq_link_date_dim` (`link_id`, `date`, `dimension_type`, `dimension_value`(100)),

    -- 关键索引: 优化 "Top N" 和 "TimeSeries" 查询
    -- (此索引可覆盖 GetAnalytics API 的所有查询)
                                           INDEX `idx_query_analytics` (`link_id`, `date`, `dimension_type`, `click_count` DESC),

    -- 外键 (确保数据一致性，随 links 表删除)
                                           CONSTRAINT `fk_asd_link_id` FOREIGN KEY (`link_id`) REFERENCES `links` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='报表：按天聚合的多维度统计表 (读密集型)';