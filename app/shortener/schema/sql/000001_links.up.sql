-- 文件: app/shortener/schema/sql/000001_links.up.sql (最终优化版)
-- 职责: 存储短码和长链接的映射，支持软删除和高效的列表过滤。

CREATE TABLE `links` (
                         `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                         `user_id` BIGINT UNSIGNED NOT NULL COMMENT '创建者ID，关联 users(id)',
                         `short_code` VARCHAR(16) NOT NULL COMMENT '短链接码 (e.g., aZ89bC)',
                         `original_url` TEXT NOT NULL COMMENT '原始长链接',
                         `visit_count` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '总访问次数 (冗余字段，由报表系统异步更新)',
                         `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用 (1=启用, 0=禁用)',
                         `expiration_time` DATETIME NULL DEFAULT NULL COMMENT '过期时间 (NULL 为永不过期)',
                         `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间 (微秒精度)',
                         `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间 (微秒精度)',
                         `deleted_at` DATETIME(6) NULL DEFAULT NULL COMMENT '软删除时间 (NULL 表示未删除)',

    -- 核心索引：确保 short_code 唯一
                         UNIQUE INDEX `uniq_short_code` (`short_code`),

    -- 优化索引：用于高性能重定向 (GET /:code)
                         INDEX `idx_redirect` (`short_code`, `deleted_at`, `is_active`),

    -- 优化索引：用于按 "status" (active/inactive) 过滤列表 (GET /links?status=active)
    -- 这个索引覆盖了 WHERE, ORDER BY, 和 LIMIT，性能最高
                         INDEX `idx_user_status_sort` (`user_id`, `deleted_at`, `is_active`, `created_at` DESC),

    -- 优化索引：用于按 "status" (expired) 过滤列表 (GET /links?status=expired)
                         INDEX `idx_user_expiration_sort` (`user_id`, `deleted_at`, `expiration_time`, `created_at` DESC),

    -- 外键约束
                         CONSTRAINT `fk_links_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短链接表 (支持软删除和高效过滤)';