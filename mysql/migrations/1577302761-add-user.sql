-- Migration: add-user
-- Created at: 2019-12-25 20:39:21
-- ====  UP  ====

BEGIN;

CREATE TABLE IF NOT EXISTS `users` (
    `id` VARCHAR(256) NOT NULL,
    `name` VARCHAR(1000) NOT NULL,

    PRIMARY KEY (`id`)
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE IF NOT EXISTS `user_permission_on_project` (
    `user_id` VARCHAR(256) NOT NULL,
    `project_uuid` VARCHAR(36) NOT NULL,
    `permission` VARCHAR(100) NOT NULL,

    PRIMARY KEY (`user_id`, `project_uuid`),
    CONSTRAINT `fk_permission_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_permission_project` FOREIGN KEY (`project_uuid`) REFERENCES `projects`(`uuid`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE `events`
    ADD COLUMN `user_id` VARCHAR(256) NULL AFTER `type`,
    ADD CONSTRAINT `fk_event_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `events`
    DROP FOREIGN KEY `fk_event_user`,
    DROP COLUMN `user_id`;

DROP TABLE IF EXISTS `user_permission_on_project`;
DROP TABLE IF EXISTS `users`;

COMMIT;
