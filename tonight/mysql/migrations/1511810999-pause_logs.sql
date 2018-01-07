-- Migration: pause_logs
-- Created at: 2017-11-27 20:29:59
-- ====  UP  ====

BEGIN;

ALTER TABLE `task_log`
    DROP PRIMARY KEY,
    ADD COLUMN `id` BIGINT NOT NULL AUTO_INCREMENT FIRST,
    DROP FOREIGN KEY `fk_task_log_task`,
    ADD PRIMARY KEY (`id`),
    ADD COLUMN `type` VARCHAR(16) DEFAULT NULL AFTER `task_id`;

ALTER TABLE `task_log`
    ADD CONSTRAINT `fk_task_log_task` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE;

UPDATE `task_log` SET `type` = "COMPLETION";

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `task_log`
    DROP COLUMN `type`,
    DROP PRIMARY KEY,
    ADD PRIMARY KEY (`task_id`, `completion`),
    DROP COLUMN `id`;

COMMIT;
