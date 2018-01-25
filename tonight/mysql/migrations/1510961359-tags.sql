-- Migration: tags
-- Created at: 2017-11-18 00:29:19
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `tags` (
    `task_id` BIGINT NOT NULL,
    `tag` VARCHAR(512) NOT NULL,

    PRIMARY KEY(`task_id`, `tag`),

    CONSTRAINT `fk_tags_tasks` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `tags`;

COMMIT;
