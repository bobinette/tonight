-- Migration: done_becomes_status
-- Created at: 2017-12-15 20:04:52
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `done`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `done` BOOLEAN NOT NULL AFTER `rank`;

COMMIT;
