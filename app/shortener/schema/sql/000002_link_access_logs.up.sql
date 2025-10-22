-- 文件: app/shortener/schema/sql/000002_link_access_logs.up.sql
-- 职责: 存储每一次访问的原始日志 (写密集型)。

CREATE TABLE `link_access_logs` (
                                    `id` BIGINT UNSIGNED AUTO_INCREMENT,
                                    `link_id` BIGINT UNSIGNED NOT NULL COMMENT '关联的链接ID links(id)',
                                    `short_code` VARCHAR(16) NOT NULL COMMENT '访问的短码 (冗余，方便查询)',
                                    `accessed_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '访问时间 (微秒精度)',
                                    `ip_address` VARCHAR(45) NOT NULL COMMENT '访问者 IP (需考虑 IPv6)',
                                    `user_agent` TEXT NULL COMMENT '访问者 User-Agent',
                                    `referer` TEXT NULL COMMENT '访问来源',
                                    `geo_country` VARCHAR(100) NULL COMMENT 'IP解析-国家 (ETL处理后填入)',
                                    `geo_city` VARCHAR(100) NULL COMMENT 'IP解析-城市 (ETL处理后填入)',
                                    `device_type` VARCHAR(50) NULL COMMENT 'UA解析-设备类型 (ETL处理后填入)',
                                    `browser_name` VARCHAR(50) NULL COMMENT 'UA解析-浏览器 (ETL处理后填入)',
                                    `os_name` VARCHAR(50) NULL COMMENT 'UA解析-操作系统 (ETL处理后填入)',

                                    PRIMARY KEY (`id`),
    -- 核心索引：仅供后台 Cron 任务聚合数据时使用
                                    INDEX `idx_link_id_accessed_at` (`link_id`, `accessed_at`),

    -- 核心索引：供后台 Cron 任务按时间范围拉取日志
                                    INDEX `idx_accessed_at` (`accessed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短链接访问日志表 (写密集型)';