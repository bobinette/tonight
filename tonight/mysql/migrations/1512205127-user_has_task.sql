-- Migration: user_has_task
-- Created at: 2017-12-02 09:58:47
-- ====  UP  ====

BEGIN;

RENAME TABLE `user` TO `users`;
RENAME TABLE `task_log` TO `task_logs`;

CREATE TABLE IF NOT EXISTS `user_has_tasks` (
    `user_id` BIGINT NOT NULL,
    `task_id` BIGINT NOT NULL,

    `created_at` DATETIME NOT NULL,

    PRIMARY KEY(`user_id`, `task_id`),
    CONSTRAINT `fk_user_has_task_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_has_task_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
)
ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

ALTER TABLE `planning`
    ADD COLUMN `user_id` BIGINT DEFAULT NULL AFTER `id`,
    ADD CONSTRAINT `fk_planning_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `planning`
    DROP FOREIGN KEY `fk_planning_user`,
    DROP COLUMN `user_id`;

DROP TABLE IF EXISTS `user_has_tasks`;

RENAME TABLE `task_logs` TO `task_log`;
RENAME TABLE `users` TO `user`;

COMMIT;
