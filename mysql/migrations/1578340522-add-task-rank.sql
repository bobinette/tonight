-- Migration: add-task-rank
-- Created at: 2020-01-06 20:55:22
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `rank` SMALLINT NULL DEFAULT NULL AFTER `status`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `rank`;

COMMIT;
