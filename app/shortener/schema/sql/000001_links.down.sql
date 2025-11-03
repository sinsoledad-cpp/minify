-- 文件: app/shortener/schema/sql/000001_links.up.sql (最终优化版)
-- 职责: 存储短码和长链接的映射，支持软删除和高效的列表过滤。

CREATE TABLE `links` (
                         `id` BIGINT UNSIGNED AUTO_INCREMENT,
                         `user_id` BIGINT UNSIGNED NOT NULL COMMENT '创建者ID，关联 users(id)',
                         `short_code` VARCHAR(16) NOT NULL COMMENT '短链接码 (e.g., aZ89bC)',
                         `original_url` TEXT NOT NULL COMMENT '原始长链接',
                         `is_active` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否启用 (1=启用, 0=禁用)',
                         `expiration_time` DATETIME NULL DEFAULT NULL COMMENT '过期时间 (NULL 为永不过期)',
                         `created_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间 (微秒精度)',
                         `updated_at` DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间 (微秒精度)',
                         `deleted_at` DATETIME(6) NULL DEFAULT NULL COMMENT '软删除时间 (NULL 表示未删除)',

                         PRIMARY KEY (`id`),
    -- 核心索引：确保 short_code 唯一
                         UNIQUE INDEX `uniq_short_code` (`short_code`),

    -- 优化索引：用于高性能重定向 (GET /:code)
                         INDEX `idx_redirect` (`short_code`, `deleted_at`, `is_active`),

    -- 优化索引：用于 Admin 后台 "status" (expired) 过滤列表 (基于 OFFSET)
                         INDEX `idx_user_expiration_sort` (`user_id`, `deleted_at`, `expiration_time`, `created_at` DESC),

    -- ⭐ 优化索引：用于 C 端游标分页 和 Admin 端 OFFSET 分页
    -- 覆盖 WHERE(user_id, deleted_at, is_active)
    -- 覆盖 ORDER BY(created_at DESC, id DESC) [C端游标]
    -- 覆盖 ORDER BY(created_at DESC) [Admin端OFFSET]
                         INDEX `idx_user_status_cursor` (`user_id`, `deleted_at`, `is_active`, `created_at` DESC, `id` DESC)

    -- 外键约束
#                          CONSTRAINT `fk_links_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短链接表 (支持软删除和高效过滤)';