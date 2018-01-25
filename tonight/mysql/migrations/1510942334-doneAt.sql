-- Migration: doneAt
-- Created at: 2017-11-17 19:12:14
-- ====  UP  ====

BEGIN;

ALTER TABLE `tasks`
    ADD COLUMN `doneAt` DATETIME(6) DEFAULT NULL AFTER `done`;

COMMIT;

-- ==== DOWN ====

BEGIN;

ALTER TABLE `tasks`
    DROP COLUMN `doneAt`;

COMMIT;
