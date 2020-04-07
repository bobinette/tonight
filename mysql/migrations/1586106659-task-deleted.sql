-- Migration: task-deleted
-- Created at: 2020-04-05 19:10:59
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `deleted` BOOLEAN NOT NULL DEFAULT 0 AFTER `updated_at`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `deleted`;

COMMIT;
