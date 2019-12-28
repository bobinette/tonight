-- Migration: add-task-status
-- Created at: 2019-12-28 11:44:11
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `status` VARCHAR(30) NOT NULL AFTER `title`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `status`;

COMMIT;
