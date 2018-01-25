-- Migration: dependencies
-- Created at: 2017-11-25 13:25:49
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `task_dependencies` (
    `task_id` BIGINT NOT NULL,
    `dependency_task_id` BIGINT NOT NULL,

    `created_at` DATETIME NOT NULL,

    PRIMARY KEY(`task_id`, `dependency_task_id`),
    CONSTRAINT `fk_task_dependencies_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_task_dependencies_dependency_task` FOREIGN KEY (`dependency_task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

COMMIT;

-- ==== DOWN ====

BEGIN;

DROP TABLE IF EXISTS `task_dependencies`;

COMMIT;
