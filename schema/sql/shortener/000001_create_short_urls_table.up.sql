CREATE TABLE IF NOT EXISTS `short_urls` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID, 主键',
    `user_id` bigint(20) unsigned NOT NULL COMMENT '创建用户的ID',
    `short_key` varchar(50) NOT NULL COMMENT '短链接的唯一key',
    `original_url` text NOT NULL COMMENT '原始长链接',
    `original_url_md5` char(32) NOT NULL COMMENT '原始长链接的MD5哈希值，用于去重',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `expires_at` timestamp NULL DEFAULT NULL COMMENT '过期时间, NULL表示永不过期',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '软删除时间, NULL表示未删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_short_key` (`short_key`),
    -- 允许多个用户缩短同一个链接，但每个用户对同一个长链接只能有一个短链接
    UNIQUE KEY `uk_user_id_url_md5` (`user_id`, `original_url_md5`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '短链接映射表';