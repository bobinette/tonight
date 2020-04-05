-- Migration: releases
-- Created at: 2020-04-04 19:25:34
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `releases` (
    `uuid` VARCHAR(36) NOT NULL,

    `title` TEXT NOT NULL,
    `description` TEXT NOT NULL,

    `project_uuid` VARCHAR(36) NOT NULL,

    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,

    PRIMARY KEY (`uuid`),
    FOREIGN KEY (`project_uuid`) REFERENCES `projects` (`uuid`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO releases (uuid, title, description, project_uuid, created_at, updated_at)
SELECT uuid, 'Backlog', '', uuid, NOW(), NOW() FROM projects;

ALTER TABLE `tasks`
    ADD COLUMN `release_uuid` VARCHAR(36) NULL DEFAULT NULL AFTER `project_uuid`;

ALTER TABLE `tasks`
    ADD CONSTRAINT `fk_task_release_uuid` FOREIGN KEY (`release_uuid`) REFERENCES `releases`(`uuid`) ON DELETE CASCADE;

UPDATE tasks SET release_uuid = project_uuid;

ALTER TABLE tasks
    DROP FOREIGN KEY `fk_task_project`,
    DROP COLUMN project_uuid;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `project_uuid` VARCHAR(36) NOT NULL AFTER `title`;

UPDATE tasks
JOIN releases ON tasks.release_uuid = releases.uuid
SET tasks.project_uuid = releases.uuid;

ALTER TABLE `tasks`
    ADD CONSTRAINT `fk_task_project` FOREIGN KEY (`project_uuid`) REFERENCES `projects`(`uuid`) ON DELETE CASCADE,
    DROP FOREIGN KEY `fk_task_release_uuid`,
    DROP COLUMN `release_uuid`;

DROP TABLE IF EXISTS `releases`;

COMMIT;
