CREATE TABLE IF NOT EXISTS `url_analytics` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID, 主键',
    `short_url_id` bigint(20) unsigned NOT NULL COMMENT '关联的short_urls表ID',
    `ip_address` varchar(45) NOT NULL COMMENT '访问者IP地址',
    `user_agent` text COMMENT '访问者User-Agent',
    `referer` text COMMENT '访问来源',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
    PRIMARY KEY (`id`),
    KEY `idx_short_url_id_created_at` (`short_url_id`, `created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '短链接访问统计表';