-- Migration: completion
-- Created at: 2017-11-25 09:13:15
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `task_log` (
    `task_id` BIGINT NOT NULL,
    `completion` SMALLINT NOT NULL,

    `description` TEXT DEFAULT NULL,

    `created_at` DATETIME NOT NULL,

    PRIMARY KEY(`task_id`, `completion`),
    CONSTRAINT `fk_task_log_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `task_log`;

COMMIT;
