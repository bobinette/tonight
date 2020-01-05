-- Migration: add-project
-- Created at: 2019-12-24 09:46:26
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `projects` (
    `uuid` VARCHAR(36) NOT NULL,

    `name` TEXT NOT NULL,

    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,

    PRIMARY KEY (`uuid`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE `tasks`
    ADD COLUMN `project_uuid` VARCHAR(36) NOT NULL AFTER `title`,
    ADD CONSTRAINT `fk_task_project` FOREIGN KEY (`project_uuid`) REFERENCES `projects`(`uuid`) ON DELETE CASCADE;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP FOREIGN KEY `fk_task_project`,
    DROP COLUMN `project_uuid`;

DROP TABLE IF EXISTS `projects`;

COMMIT;
