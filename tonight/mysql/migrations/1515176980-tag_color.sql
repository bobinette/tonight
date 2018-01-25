-- Migration: tag_color
-- Created at: 2018-01-05 19:29:40
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `user_customs_tags` (
    `user_id` BIGINT NOT NULL,
    `tag` VARCHAR(512) NOT NULL,
    `colour` CHAR(7), -- hexadecimal colour

    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,

    PRIMARY KEY(`user_id`, `tag`),
    CONSTRAINT `fk_user_customs_tags_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `user_customs_tags`;

COMMIT;
